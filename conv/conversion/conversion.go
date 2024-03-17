//go:build ignore_until_step_12

package conversion

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
	"github.com/Silicon-Ally/etap2sf/utils"
)

func (in *Input) Convert() (*Output, error) {
	doConversion := func(name string, fn func() []error) error {
		errors := fn()
		if len(errors) == 0 {
			return nil
		}
		filePath, _ := utils.WriteErrorsToTempFile(errors, "errors-approaches")
		return fmt.Errorf("converting %s failed, errors at %s", name, filePath)
	}

	result := &io{in: in, out: &Output{
		accountsByRefs:   map[string]*sfenterprise.Account{},
		contactsByRefs:   map[string]*sfenterprise.Contact{},
		refSubstitutions: map[string]string{},
	}}
	if err := doConversion("approaches", result.convertApproaches); err != nil {
		return nil, err
	}
	if err := doConversion("funds", result.convertFunds); err != nil {
		return nil, err
	}
	if err := doConversion("accounts", result.convertAccounts); err != nil {
		return nil, err
	}
	if err := result.linkSecondaryRoleRefsToPrimaryRoleRefInMap(); err != nil {
		return nil, err
	}
	if err := doConversion("relationships", result.convertRelationships); err != nil {
		return nil, err
	}
	if err := doConversion("households", result.assignContactsToHouseholdAccounts); err != nil {
		return nil, err
	}
	if err := doConversion("journalentries", result.convertJournalEntries); err != nil {
		return nil, err
	}
	if err := doConversion("campaigns", result.convertCampaigns); err != nil {
		return nil, err
	}
	if err := doConversion("fund allocations", result.assignGauAllocations); err != nil {
		return nil, err
	}
	if err := doConversion("attachments", result.convertAttachments); err != nil {
		return nil, err
	}
	return result.out, nil
}

func (i *io) convertApproaches() []error {
	errors := []error{}
	// NOTE: we ignore approaches because of the particulars of our implementation.
	// If you want to map approaches, you're welcome to, but there isn't an obvious
	// default way of doing this in salesforce.
	return errors
}

func (i *io) convertCampaigns() []error {
	errors := []error{}

	campaigns := map[string]bool{}
	for _, c := range i.in.Campaigns {
		campaigns[c] = true
	}
	for cn := range campaigns {
		if cn != "" {
			out, err := i.transformETAPCampaignToSalesforceCampaign(cn)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting campaign %q: %w", cn, err))
			}
			i.out.Campaigns = append(i.out.Campaigns, out)
		}
	}
	return errors
}

func (i *io) convertFunds() []error {
	errors := []error{}
	for _, f := range i.in.Funds {
		out, err := i.transformETAPFundToSalesforceGeneralAccountingUnit(f)
		if err != nil {
			errors = append(errors, fmt.Errorf("converting fund %q: %w", *f.Name, err))
		}
		i.out.GeneralAccountingUnits = append(i.out.GeneralAccountingUnits, out)
	}
	return errors
}

func (i *io) convertAccounts() []error {
	errors := []error{}
	for _, a := range i.in.Accounts {
		isOrg, err := i.isAccountAnOrganization(a)
		if err != nil {
			errors = append(errors, fmt.Errorf("checking if account is an organization: %w", err))
			continue
		}
		if isOrg {
			if err := i.convertOrganization(a); err != nil {
				errors = append(errors, fmt.Errorf("converting to org: %w", err))
			}
		} else {
			if err := i.convertIndividual(a); err != nil {
				errors = append(errors, fmt.Errorf("converting to individual: %w", err))
			}
		}
	}
	return errors
}

func (i *io) isAccountAnOrganization(a *generated.Account) (bool, error) {
	df, err := i.in.CustomFields.LookupByName("Account Type", a.AccountDefinedValues)
	if err != nil {
		return false, fmt.Errorf("looking up account type: %w", err)
	}
	if df == nil || df.Value == nil {
		// We assume anyone with the account value not set is an individual.
		return false, nil
	}
	if *df.Value != "Individual" && *df.Value != "Organization" {
		return true, fmt.Errorf("account type is not individual or organization: %q", *df.Value)
	}
	return *df.Value == "Organization", nil
}

func (i *io) convertOrganization(a *generated.Account) error {
	result, err := i.transformETAPAccountToSalesforceAccount(a)
	if err != nil {
		return fmt.Errorf("converting account to account: %w", err)
	}
	if result.Etap_Account_AccountRoleType__c != nil && *result.Etap_Account_AccountRoleType__c != 0 {
		return fmt.Errorf("account has a role type of %d, but is an organization: %q", *a.AccountRoleType, *a.Ref)
	}
	i.out.Accounts = append(i.out.Accounts, result)
	i.out.accountsByRefs[*a.Ref] = result
	return nil
}

func (i *io) convertIndividual(a *generated.Account) error {
	result, err := i.transformETAPAccountToSalesforceContact(a)
	if err != nil {
		return fmt.Errorf("converting account to contact: %w", err)
	}
	// This is either a tribute account or a user account (as opposed to a donor account)
	if result.Etap_Account_AccountRoleType__c != nil && *result.Etap_Account_AccountRoleType__c != 0 {
		betterRef := i.findDonorAccountForUser(a)
		// If we can find a better representation of the user/tribute (i.e. a donor), use that. Otherwise,
		// create the contact in Salesforce as usual.
		if betterRef != "" {
			// In this case, we just make sure we're putting the user into ETap elsewhere.
			i.out.refSubstitutions[*a.Ref] = betterRef
			return nil
		}
	}
	i.out.Contacts = append(i.out.Contacts, result)
	i.out.contactsByRefs[*a.Ref] = result
	return nil
}

func (i *io) findDonorAccountForUser(a *generated.Account) string {
	if a.DonorRoleRef != nil && *a.DonorRoleRef != "" {
		return *a.DonorRoleRef
	}
	// The reason we do it this way is to preference Donor > Tribute > User, and guarantee one of them exists if we return "".
	lowestRoleType := *a.AccountRoleType
	lowestRoleRef := ""
	for _, r := range i.in.Accounts {
		if strings.EqualFold(*r.Email, *a.Email) && *r.AccountRoleType < lowestRoleType {
			lowestRoleRef = *r.Ref
			lowestRoleType = *r.AccountRoleType
		}
	}
	return lowestRoleRef
}

func (result *io) linkSecondaryRoleRefsToPrimaryRoleRefInMap() error {
	for replace, with := range result.out.refSubstitutions {
		if withContact, ok := result.out.contactsByRefs[with]; ok {
			result.out.contactsByRefs[replace] = withContact
		} else if withAccount, ok := result.out.accountsByRefs[with]; ok {
			result.out.accountsByRefs[replace] = withAccount
		} else {
			return fmt.Errorf("substituting %q for %q, but %q not found", replace, with, with)
		}
	}
	return nil
}

func (i *io) assignContactsToHouseholdAccounts() []error {
	type HHID string
	type AccountRef string
	householdToRefs := map[HHID][]AccountRef{}
	refToHousehold := map[AccountRef]HHID{}
	householdToHead := map[HHID]AccountRef{}
	generatesContact := map[AccountRef]bool{}

	join := func(a1, a2 AccountRef) error {
		hh1, ok := refToHousehold[a1]
		if !ok {
			return fmt.Errorf("couldn't find household for %q", a1)
		}
		hh2, ok := refToHousehold[a2]
		if !ok {
			return fmt.Errorf("couldn't find household for %q", a1)
		}
		if hh1 == hh2 {
			return nil
		}
		for _, r := range householdToRefs[hh2] {
			refToHousehold[r] = hh1
		}
		householdToRefs[hh1] = append(householdToRefs[hh1], householdToRefs[hh2]...)
		delete(householdToRefs, hh2)
		return nil
	}
	makeHOH := func(ref AccountRef) error {
		hh, ok := refToHousehold[ref]
		if !ok {
			return fmt.Errorf("couldn't find household for %q", ref)
		}
		if ogh, ok := householdToHead[hh]; !ok {
			householdToHead[hh] = ref
		} else if ogh != ref {
			return fmt.Errorf("household %q already has head %q", hh, ogh)
		}
		return nil
	}

	errors := []error{}
	for i, c := range i.out.Contacts {
		hhid := HHID(fmt.Sprintf("%d", i))
		ref := AccountRef(*c.Etap_Account_Ref__c)
		if ref == "" {
			errors = append(errors, fmt.Errorf("contact %q has no ref", *c.Name))
			continue
		}
		householdToRefs[hhid] = []AccountRef{ref}
		refToHousehold[ref] = hhid
		generatesContact[ref] = true
	}

	for _, r := range i.in.Relationships {
		if r.HohAccount == nil {
			errors = append(errors, fmt.Errorf("at %q: nil HOH account", *r.Ref))
		}
		switch *r.HohAccount {
		case 0:
			continue // Non-household account
		case 1, 2:
			a1 := AccountRef(*r.Account1Ref)
			a2 := AccountRef(*r.Account2Ref)
			ok1 := generatesContact[a1]
			ok2 := generatesContact[a2]
			if !ok1 || !ok2 {
				fmt.Printf("WARNING - Household Account Inclues Organization - %s and %s\n", *r.Account1Name, *r.Account2Name)
				continue
			}
			if err := join(a1, a2); err != nil {
				errors = append(errors, fmt.Errorf("joining %q and %q: %w", *r.Account1Name, *r.Account2Name, err))
				continue
			}
			hoh := a1
			if *r.HohAccount == 2 {
				hoh = a2
			}
			if err := makeHOH(hoh); err != nil {
				errors = append(errors, fmt.Errorf("making %q HOH: %w", hoh, err))
				continue
			}
		}
	}

	// Prevents map non-determinism
	hhids := []HHID{}
	for hhid := range householdToRefs {
		hhids = append(hhids, hhid)
	}
	sort.Slice(hhids, func(i, j int) bool {
		return hhids[i] < hhids[j]
	})

	for _, hhid := range hhids {
		refs := householdToRefs[hhid]
		hoh, ok := householdToHead[hhid]
		if !ok {
			hoh = refs[0]
		}
		hohc, err := i.lookupSFContactByRef(string(hoh))
		if err != nil {
			errors = append(errors, fmt.Errorf("looking up head of household %q: %w", hoh, err))
			continue
		}
		sort.Slice(refs, func(i, j int) bool {
			return refs[i] < refs[j]
		})
		syntheticHouseholdRef := "SynthHH" + string(refs[0])
		household := &sfenterprise.Account{
			Name:                         ptr(fmt.Sprintf("%s %s Household", *hohc.FirstName, *hohc.LastName)),
			Etap_MultiObject_EtapRef__c:  ptr(syntheticHouseholdRef),
			Etap_MigrationExplanation__c: ptr(fmt.Sprintf("Household auto-generated by etapestry migration code based off of ETapestry Household Relationships - includes %+v", refs)),
		}
		household.Etap_MigrationTime__c = NowXSD()
		household.CreatedDate = clonePtr(hohc.CreatedDate)
		household.LastModifiedDate = clonePtr(hohc.LastModifiedDate)
		household.Type = ptr(sfenterprise.Account_Type_Household)
		if i.in.HouseholdAccountRecordType == "" {
			return []error{fmt.Errorf("household account record type not set")}
		}
		household.RecordType = &sfenterprise.RecordType{
			Id:          ptr(i.in.HouseholdAccountRecordType),
			SobjectType: ptr(sfenterprise.RecordType_SobjectType_Account),
			Name:        ptr("Household Account"),
		}

		isJointPersona := hohc.Etap_Account_PersonaTypes__c != nil && strings.Contains(*hohc.Etap_Account_PersonaTypes__c, "joint")
		isSolo := len(refs) == 0
		if isSolo || isJointPersona {
			if hohc.MailingAddress != nil {
				household.ShippingAddress = clonePtr(hohc.MailingAddress)
				household.BillingAddress = clonePtr(hohc.MailingAddress)
			}
		}

		i.out.Accounts = append(i.out.Accounts, household)
		placeholderAccountID, err := idPlaceholderForRef(&syntheticHouseholdRef)
		if err != nil {
			errors = append(errors, fmt.Errorf("generating placeholder account id: %w", err))
			continue
		}
		for _, ref := range refs {
			c, err := i.lookupSFContactByRef(string(ref))
			if err != nil {
				errors = append(errors, fmt.Errorf("looking up contact %q: %w", ref, err))
				continue
			}
			// Why do this? In populating this field later, we use a utility which populates the
			// value directly into this pointer. If it's shared across pointers, you'll get a
			// "expected to need replacement, but was" error, because it was replaced from the
			// conversion of another househol member.
			c.AccountId = clonePtr(placeholderAccountID)
		}
	}
	return errors
}

func (i *io) convertRelationships() []error {
	errors := []error{}
	for _, r := range i.in.Relationships {
		if individuals, accounts, err := i.isRelationshipBetweenIndividualsOrAccounts(r); err != nil {
			errors = append(errors, fmt.Errorf("checking if relationship is between individuals: %w", err))
			continue
		} else if individuals {
			// NOTE: we make out2 + add it to the set inside of here. It's not ideal, but it allows us fuller access to context.
			out, err := i.transformETAPRelationshipToSalesforceRelationship(r)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting relationship between indidividuals %q: %w", *r.Ref, err))
				continue
			}
			i.out.Relationships = append(i.out.Relationships, out)
		} else if accounts {
			fmt.Printf("WARNING - relationship between accounts %s and %s - skipping\n", *r.Account1Name, *r.Account2Name)
		} else {
			out, err := i.transformETAPRelationshipToSalesforceAffiliation(r)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting relationship between account and contact %q: %w", *r.Ref, err))
				continue
			}
			i.out.Affiliations = append(i.out.Affiliations, out)
		}
	}
	return errors
}

func (i *io) isRelationshipBetweenIndividualsOrAccounts(r *generated.Relationship) (bool, bool, error) {
	a1, _ := i.lookupSFAccountByRef(*r.Account1Ref)
	a2, _ := i.lookupSFAccountByRef(*r.Account2Ref)
	if a1 != nil || a2 != nil {
		if a1 != nil && a2 != nil {
			return false, true, nil
		}
		return false, false, nil
	}
	c1, _ := i.lookupSFContactByRef(*r.Account1Ref)
	c2, _ := i.lookupSFContactByRef(*r.Account2Ref)
	if c1 == nil || c2 == nil {
		return false, false, fmt.Errorf("very odd - relationship between neither accounts nor contacts - couldn't find anything")
	}
	return true, false, nil
}

func (i *io) convertJournalEntries() []error {
	errors := []error{}
	for _, je := range i.in.JournalEntries {
		err := i.convertJournalEntry(je)
		if err != nil {
			errors = append(errors, fmt.Errorf("converting journal entry %q: %w", je.Ref(), err))
		}
	}
	return errors
}

func (i *io) convertJournalEntry(je *overrides.JournalEntry) error {
	if je.Gift != nil {
		opp, err := i.transformETAPGiftToSalesforceOpportunity(je.Gift)
		if err != nil {
			return fmt.Errorf("converting gift: %w", err)
		}
		i.out.Opportunities = append(i.out.Opportunities, opp)
		return nil
	}
	if je.Contact != nil {
		task, err := i.transformETAPContactToSalesforceTask(je.Contact)
		if err != nil {
			return fmt.Errorf("converting contact to task: %w", err)
		}
		i.out.Tasks = append(i.out.Tasks, task)
		context, err := i.transformETAPContactToSalesforceEtapAdditionalContext(je.Contact)
		if err != nil {
			return fmt.Errorf("converting contact to additional context: %w", err)
		}
		i.out.AdditionalContexts = append(i.out.AdditionalContexts, context)
		return nil
	}
	if je.Note != nil {
		task, err := i.transformETAPNoteToSalesforceTask(je.Note)
		if err != nil {
			return fmt.Errorf("converting note: %w", err)
		}
		i.out.Tasks = append(i.out.Tasks, task)
		context, err := i.transformETAPNoteToSalesforceEtapAdditionalContext(je.Note)
		if err != nil {
			return fmt.Errorf("converting note to additional context: %w", err)
		}
		i.out.AdditionalContexts = append(i.out.AdditionalContexts, context)
		return nil
	}
	if je.Payment != nil {
		p, err := i.transformETAPPaymentToSalesforcePayment(je.Payment)
		if err != nil {
			return fmt.Errorf("converting payment: %w", err)
		}
		i.out.Payments = append(i.out.Payments, p)
		return nil
	}
	if je.Pledge != nil {
		opp, err := i.transformETAPPledgeToSalesforceOpportunity(je.Pledge)
		if err != nil {
			return fmt.Errorf("converting pledge: %w", err)
		}
		i.out.Opportunities = append(i.out.Opportunities, opp)
		return nil
	}
	if je.RecurringGift != nil {
		rd, err := i.transformETAPRecurringGiftToSalesforceOpportunity(je.RecurringGift)
		if err != nil {
			return fmt.Errorf("converting recurring gift: %w", err)
		}
		i.out.Opportunities = append(i.out.Opportunities, rd)
		return nil
	}
	if je.RecurringGiftSchedule != nil {
		rd, err := i.transformETAPRecurringGiftScheduleToSalesforceRecurringDonation(je.RecurringGiftSchedule)
		if err != nil {
			return fmt.Errorf("converting recurring gift schedule: %w", err)
		}
		i.out.RecurringDonations = append(i.out.RecurringDonations, rd)
		return nil
	}
	if je.SoftCredit != nil {
		if _, ok := i.out.accountsByRefs[*je.SoftCredit.AccountRef]; ok {
			return i.convertAccountSoftCredit(je.SoftCredit)
		} else {
			return i.convertContactSoftCredit(je.SoftCredit)
		}
	}
	if je.SegmentedDonation != nil {
		opp, err := i.transformETAPSegmentedDonationToSalesforceOpportunity(je.SegmentedDonation)
		if err != nil {
			return fmt.Errorf("converting segmented donation: %w", err)
		}
		i.out.Opportunities = append(i.out.Opportunities, opp)
		return nil
	}
	if je.Disbursement != nil {
		// Disbursements are a huge challenge in Salesforce - there isn't anything along the same lines.
		// After talking with the client here, we agreed to skip these, because they didn't have significant
		// usage in their workflows or data.
		return fmt.Errorf("skipping disbursement, which aren't supported: %q", je.Ref)
	}
	return fmt.Errorf("journal entry has unsupported kind: %+v", je)
}

func (i *io) convertContactSoftCredit(in *generated.SoftCredit) error {
	out, err := i.transformETAPSoftCreditToSalesforcePartialSoftCredit(in)
	if err != nil {
		if errors.Is(err, &IsMissingHardCredit{}) {
			fmt.Printf("WARNING - skipping a soft credit (%s) without a corresponding hard credit (%s).\n", *in.Ref, *in.HardCreditRef)
			return nil
		}
		return fmt.Errorf("converting soft credit to partial soft credit: %w", err)
	}
	i.out.PartialSoftCredits = append(i.out.PartialSoftCredits, out)
	return nil
}

func (i *io) convertAccountSoftCredit(in *generated.SoftCredit) error {
	out, err := i.transformETAPSoftCreditToSalesforceAccountSoftCredit(in)
	if err != nil {
		if errors.Is(err, &IsMissingHardCredit{}) {
			fmt.Printf("WARNING - skipping a soft credit (%s) without a corresponding hard credit (%s).\n", *in.Ref, *in.HardCreditRef)
			return nil
		}
		return fmt.Errorf("converting soft credit to account soft credit: %w", err)
	}
	i.out.AccountSoftCredits = append(i.out.AccountSoftCredits, out)
	return nil
}

func (i *io) assignGauAllocations() []error {
	errors := []error{}
	for _, o := range i.out.Opportunities {
		fundName := ""
		if o.Etap_Gift_Fund__c != nil && *o.Etap_Gift_Fund__c != "" {
			fundName = *o.Etap_Gift_Fund__c
		}
		if o.Etap_Pledge_Fund__c != nil && *o.Etap_Pledge_Fund__c != "" {
			fundName = *o.Etap_Pledge_Fund__c
		}
		if o.Etap_RecurringGift_Fund__c != nil && *o.Etap_RecurringGift_Fund__c != "" {
			fundName = *o.Etap_RecurringGift_Fund__c
		}
		if fundName == "" {
			continue
		}

		var fund string
		for _, f := range i.in.Funds {
			if *f.Name == fundName {
				fund = *f.Ref
				break
			}
		}
		if fund == "" {
			return []error{fmt.Errorf("couldn't find fund %q", fundName)}
		}

		out := &sfenterprise.Npsp__Allocation__c{}

		if id, err := idPlaceholderForRef(&fund); err != nil {
			errors = append(errors, fmt.Errorf("creating placeholder for account gift maker: %w", err))
			continue
		} else {
			out.Npsp__General_Accounting_Unit__c = id
		}
		if id, err := idPlaceholderForRef(clonePtr(o.Etap_MultiObject_EtapRef__c)); err != nil {
			errors = append(errors, fmt.Errorf("creating placeholder for : %w", err))
			continue
		} else {
			out.Npsp__Opportunity__c = id
		}
		out.Npsp__Amount__c = clonePtr(o.Amount)
		out.Npsp__Percent__c = ptr(100.0)
		out.Etap_MultiObject_EtapRef__c = ptr(*o.Etap_MultiObject_EtapRef__c + "-alloc")
		if *o.Amount < 0 {
			continue
		}
		i.out.GAUAllocations = append(i.out.GAUAllocations, out)
	}
	return errors
}

func (i *io) convertAttachments() []error {
	errors := []error{}
	for _, je := range i.in.JournalEntries {
		for _, a := range je.Attachments() {
			out, err := i.transformETAPAttachmentToSalesforceContentVersion(a)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting attachment %q: %w", *a.Ref, err))
			}
			i.out.ContentVersions = append(i.out.ContentVersions, out)
			link, err := i.pureManualCreateContentDocumentLink(je, a)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting attachment link: %q: %w", *a.Ref, err))
			}
			i.out.ContentDocumentLinks = append(i.out.ContentDocumentLinks, link)
			addc, err := i.transformETAPAttachmentToSalesforceEtapAdditionalContext(a)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting attachment additional context: %q: %w", *a.Ref, err))
			}
			i.out.AdditionalContexts = append(i.out.AdditionalContexts, addc)
		}
	}
	return errors
}

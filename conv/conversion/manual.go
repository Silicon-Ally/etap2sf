//go:build ignore_until_step_12

package conversion

import (
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/attachments/exportfiles"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
	"github.com/Silicon-Ally/etap2sf/utils"
)

func (io *io) manualTransformETAPCampaignToSalesforceCampaign(in string, out *sfenterprise.Campaign) error {
	if name, err := required(errIfLongerThan(&in, 80)); err != nil {
		return fmt.Errorf("campaign name: %w", err)
	} else {
		out.Name = name
	}

	out.IsActive = ptr(true)
	out.Status = ptr(sfenterprise.Campaign_Status_InProgress)
	out.Description = trimIfLongerThan(ptr(fmt.Sprintf("Auto Generated from eTapestry Campaign %q", in)), 255)
	out.Etap_Campaign_Name__c = ptr(in)

	out.Etap_MultiObject_EtapRef__c = ptr(campaignPlaceholderRef(in))
	explanation := fmt.Sprintf("This campaign was generated from the campaign migration spreadsheet, with campaign column %q.", in)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}

	out.CreatedById = &io.in.AttributedUserId
	out.LastModifiedById = &io.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	return nil
}

func (i *io) manualTransformETAPFundToSalesforceGeneralAccountingUnit(in *generated.Fund, out *sfenterprise.Npsp__General_Accounting_Unit__c) error {
	if name, err := required(errIfLongerThan(in.Name, 80)); err != nil {
		return fmt.Errorf("fund name: %w", err)
	} else {
		out.Name = name
	}

	if in.Disabled != nil {
		out.Npsp__Active__c = ptr(!*in.Disabled)
	}
	out.Npsp__Description__c = in.Note

	explanation := fmt.Sprintf("This general accounting unit was generated from an eTapestry fund named %q.", *in.Name)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	return nil
}

func (io *io) manualTransformETAPRelationshipToSalesforceAffiliation(in *generated.Relationship, out *sfenterprise.Npe5__Affiliation__c) error {
	c1, _ := io.lookupSFContactByRef(*in.Account1Ref)
	c2, _ := io.lookupSFContactByRef(*in.Account2Ref)
	a1, _ := io.lookupSFAccountByRef(*in.Account1Ref)
	a2, _ := io.lookupSFAccountByRef(*in.Account2Ref)
	if c1 == nil && c2 == nil {
		return fmt.Errorf("neither contact exists")
	}
	if c1 != nil && c2 != nil {
		return fmt.Errorf("both contacts exist")
	}
	if a1 == nil && a2 == nil {
		return fmt.Errorf("neither account exists")
	}
	if a1 != nil && a2 != nil {
		return fmt.Errorf("both accounts exists")
	}
	if c1 == nil && a1 == nil {
		return fmt.Errorf("contact 1 and account 1 both nil")
	}
	if c2 == nil && a2 == nil {
		return fmt.Errorf("contact 2 and account 2 both nil")
	}
	contactRef := in.Account1Ref
	accountRef := in.Account2Ref
	if c1 == nil {
		contactRef = in.Account2Ref
		accountRef = in.Account1Ref
	}

	out.Npe5__Status__c = ptr(sfenterprise.Npe5Affiliation_npe5Status_Current)
	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDate(in.StartDate); err != nil {
		return fmt.Errorf("start date: %w", err)
	} else {
		out.Npe5__StartDate__c = date
	}
	if date, err := AttemptToParseNilableDate(in.EndDate); err != nil {
		return fmt.Errorf("end date: %w", err)
	} else if date != nil {
		out.Npe5__EndDate__c = date
		if date.ToGoTime().Before(time.Now()) {
			out.Npe5__Status__c = ptr(sfenterprise.Npe5Affiliation_npe5Status_Former)
		}
	}

	if contactPlaceholder, err := idPlaceholderForRef(contactRef); err != nil {
		return fmt.Errorf("contact placeholder: %w", err)
	} else {
		out.Npe5__Contact__c = contactPlaceholder
	}
	if accountPlaceholder, err := idPlaceholderForRef(accountRef); err != nil {
		return fmt.Errorf("account placeholder: %w", err)
	} else {
		out.Npe5__Organization__c = accountPlaceholder
	}

	roleStr := fmt.Sprintf("%s / %s", *in.Type.Role1, *in.Type.Role2)
	if role, err := errIfLongerThan(&roleStr, 255); err != nil {
		return fmt.Errorf("role: %w", err)
	} else {
		out.Npe5__Role__c = role
	}
	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Npe5__Description__c = desc
	}
	explanation := fmt.Sprintf("This affiliation was generated from an eTapestry relationship of (%s : %s) between %q and %q.", *in.Type.Role1, *in.Type.Role2, *in.Account1Name, *in.Account2Name)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &io.in.AttributedUserId
	out.LastModifiedById = &io.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last mod date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	return nil
}

func (io *io) manualTransformETAPRelationshipToSalesforceRelationship(in *generated.Relationship, out *sfenterprise.Npe4__Relationship__c) error {
	explanation := fmt.Sprintf("This relationship was generated from an eTapestry relationship of (%s : %s) between %q and %q.", *in.Type.Role1, *in.Type.Role2, *in.Account1Name, *in.Account2Name)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &io.in.AttributedUserId
	out.LastModifiedById = &io.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last mod date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Npe4__Description__c = desc
	}
	out.Npe4__Status__c = ptr(sfenterprise.Npe4Relationship_npe4Status_Current)

	// NOTE THINGS ABOVE THIS LINE WILL APPEAR IN BOTH RELATIONSHIPS...
	out2, err := utils.CloneXML(out)
	if err != nil {
		return fmt.Errorf("cloning relationship: %w", err)
	}
	// AND THINGS BELOW THIS LINE WILL ONLY APPEAR IN ONE!

	id1, err := idPlaceholderForRef(in.Account1Ref)
	if err != nil {
		return fmt.Errorf("creating contact 1 placeholder: %w", err)
	}
	id2, err := idPlaceholderForRef(in.Account2Ref)
	if err != nil {
		return fmt.Errorf("creating contact 2 placeholder: %w", err)
	}

	ogRef := *out.Etap_Relationship_Ref__c

	out.Npe4__Contact__c = clonePtr(id1)
	out.Npe4__RelatedContact__c = clonePtr(id2)
	// We append the x-of-2 suffixes so that we don't collide when doing upserts on this key
	out.Etap_Relationship_Ref__c = ptr(ogRef + ".1-of-2")

	out2.Npe4__Contact__c = clonePtr(id2)
	out2.Npe4__RelatedContact__c = clonePtr(id1)
	out2.Etap_Relationship_Ref__c = ptr(ogRef + ".2-of-2")

	io.out.Relationships = append(io.out.Relationships, out2)
	return nil
}

func (io *io) manualTransformETAPAccountToSalesforceContact(in *generated.Account, out *sfenterprise.Contact) error {
	out.Birthdate = out.Etap_PersonalInfo_Birthdate__c
	if out.Birthdate == nil {
		if t, err := AttemptToParseNilableDate(out.Etap_PersonalInfo_BirthdayMMYYYY__c); err != nil {
			return fmt.Errorf("parsing secondary birthday: %w", err)
		} else {
			out.Birthdate = t
		}
	}

	if out.Etap_PersonalInfo_Deceased__c != nil {
		if *out.Etap_PersonalInfo_Deceased__c == "Yes" {
			out.Npsp__Deceased__c = ptr(true)
		} else {
			return fmt.Errorf("unknown deceased value: %q", *out.Etap_PersonalInfo_Deceased__c)
		}
	}
	// Note: we just have no insight into how emails are stored in ETAP, but creating new fields
	// for them would break a bunch of Salesforce NPSP behavior. Here, we just assign into slots
	// and acknowledge that the context of work/personal is lost.
	emails := cleanEmails(in.Email)
	if len(emails) >= 1 {
		out.Email = &emails[0]
	}
	if len(emails) >= 2 {
		out.Npe01__AlternateEmail__c = &emails[1]
	}
	if len(emails) >= 3 {
		out.Npe01__HomeEmail__c = &emails[2]
	}
	if len(emails) >= 4 {
		out.Npe01__WorkEmail__c = &emails[3]
	}
	if len(emails) >= 5 {
		return fmt.Errorf("too many emails: %q", emails)
	}

	out.HasOptedOutOfEmail = in.OptedOut
	if out.Etap_MovesMgmt_CCUnsubscribed__c != nil {
		if *out.Etap_MovesMgmt_CCUnsubscribed__c == "Unsubscribed" {
			out.HasOptedOutOfEmail = ptr(true)
		} else {
			return fmt.Errorf("unknown email opt out value: %q", *out.Etap_MovesMgmt_CCUnsubscribed__c)
		}
	}

	out.MailingAddress = &sfenterprise.Address{
		Street:     in.Address,
		State:      in.State,
		PostalCode: in.PostalCode,
		City:       in.City,
		Country:    in.Country,
	}

	if in.Title != nil && *in.Title != "" {
		sal, err := parseTitleToSalutation(*in.Title)
		if err != nil {
			return fmt.Errorf("parsing title %q: %w", *in.Title, err)
		}
		out.Salutation = ptr(sal)
		if gender, genderIdentity, err := salutationToGender(sal); err == nil {
			out.Gender__c = ptr(gender)
			out.GenderIdentity = ptr(genderIdentity)
		}
	}

	// out.Name is a reserved field, composed out of these five elements:
	out.FirstName = in.FirstName
	out.LastName = in.LastName
	if out.LastName == nil || *out.LastName == "" {
		out.LastName = ptr("[Unknown Last Name]")
	}
	out.MiddleName = in.MiddleName
	out.Suffix = in.Suffix
	out.Title = in.Title

	if out.Etap_PersonalInfo_PreferredPronouns__c != nil {
		ps, err := parsePreferredPronouns(*out.Etap_PersonalInfo_PreferredPronouns__c)
		if err != nil {
			return fmt.Errorf("parsing preferred pronouns %q: %w", *out.Etap_PersonalInfo_PreferredPronouns__c, err)
		}
		out.Pronouns = ptr(ps)
		if gender, genderIdentity, err := pronounsToGender(ps); err == nil {
			out.Gender__c = ptr(gender)
			out.GenderIdentity = ptr(genderIdentity)
		}
	}

	for _, p := range in.Phones.Items {
		switch *p.Type {
		case "Cell", "":
			if out.MobilePhone != nil {
				return fmt.Errorf("too many mobiles: %q and %q", *out.MobilePhone, *p.Number)
			}
			out.MobilePhone = p.Number
		case "Fax":
			if out.Fax != nil {
				return fmt.Errorf("too many faxes: %q and %q", *out.Fax, *p.Number)
			}
			out.Fax = p.Number
		case "Voice":
			if out.OtherPhone != nil {
				return fmt.Errorf("too many other phones: %q and %q", *out.OtherPhone, *p.Number)
			}
			out.OtherPhone = p.Number
		case "Work":
			if out.Npe01__WorkPhone__c != nil {
				return fmt.Errorf("too many work phones: %q and %q", *out.Npe01__WorkPhone__c, *p.Number)
			}
			out.Npe01__WorkPhone__c = p.Number
		case "Home":
			if out.HomePhone != nil {
				return fmt.Errorf("too many home phones: %q and %q", *out.HomePhone, *p.Number)
			}
			out.HomePhone = p.Number
		default:
			return fmt.Errorf("unknown phone type: %q", *p.Type)
		}
	}
	out.CreatedById = &io.in.AttributedUserId
	out.LastModifiedById = &io.in.AttributedUserId
	out.OwnerId = &io.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if out.Etap_AccountInformation_Constituency__c != nil {
		s := string(*out.Etap_AccountInformation_Constituency__c)
		roles := []string{}
		allRoles := []string{
			"Intern",
			"Staff",
			"Volunteer",
			"Board Member",
			"Former Staff",
			"Former Board Member",
		}
		for _, role := range allRoles {
			if strings.Contains(s, role) {
				roles = append(roles, role)
			}
		}
		// Slashes cause an issue in import so we use a differnet API name for them.
		if strings.Contains(s, "Former Volunteer/Intern") {
			roles = append(roles, "Former Volunteer or Intern")
		}
		sort.Strings(roles)
		out.Roles__c = JoinEnumsWithSemicolons(roles)
	}

	if date, err := AttemptToParseNilableDateTime(in.AccountCreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.AccountLastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Description = desc
	}
	explanation := fmt.Sprintf("This contact was generated from an ETapestry account named %q.", *in.Name)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	return nil
}

func (io *io) manualTransformETAPAccountToSalesforceAccount(in *generated.Account, out *sfenterprise.Account) error {
	explanation := fmt.Sprintf("This account was generated from an organization eTapestry account named %q.", *in.Name)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}

	out.CreatedById = &io.in.AttributedUserId
	out.LastModifiedById = &io.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.AccountCreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.AccountLastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}

	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Description = desc
	}

	out.Name = in.Name

	out.BillingAddress = &sfenterprise.Address{
		Street:     in.Address,
		State:      in.State,
		PostalCode: in.PostalCode,
		City:       in.City,
		Country:    in.Country,
	}
	out.ShippingAddress = clonePtr(out.BillingAddress)

	out.Type = ptr(sfenterprise.Account_Type_Other)
	if out.Etap_AccountInformation_Constituency__c != nil {
		s := strings.ToLower(string(*out.Etap_AccountInformation_Constituency__c))
		if strings.Contains(s, "church") {
			out.Type = ptr(sfenterprise.Account_Type_ReligiousOrganization)
		} else if strings.Contains(s, "foundation") {
			out.Type = ptr(sfenterprise.Account_Type_Foundation)
		} else if strings.Contains(s, "business") {
			out.Type = ptr(sfenterprise.Account_Type_Corporate)
		}
	}

	if io.in.OrganizationAccountRecordType == "" {
		return fmt.Errorf("organization account record type not set")
	}
	out.RecordType = &sfenterprise.RecordType{
		Id:          ptr(io.in.OrganizationAccountRecordType),
		SobjectType: ptr(sfenterprise.RecordType_SobjectType_Account),
		Name:        ptr("Organization"),
	}

	for _, p := range in.Phones.Items {
		switch *p.Type {
		case "Cell", "Voice", "Work", "Home", "":
			if out.Phone != nil && *out.Phone != "" && *out.Phone != *p.Number {
				og := ""
				if out.AdditionalPhoneNumbers__c != nil {
					og = *out.AdditionalPhoneNumbers__c
				}
				out.AdditionalPhoneNumbers__c = ptr(strings.TrimSpace(og + " " + *p.Number))
			} else {
				out.Phone = p.Number
			}
		case "Fax":
			if out.Fax != nil && *out.Fax != "" && *out.Fax != *p.Number {
				og := ""
				if out.AdditionalPhoneNumbers__c != nil {
					og = *out.AdditionalPhoneNumbers__c
				}
				out.AdditionalPhoneNumbers__c = ptr(strings.TrimSpace(og + " " + *p.Number))
			} else {
				out.Fax = p.Number
			}
		default:
			return fmt.Errorf("unknown phone type: %q", *p.Type)
		}
	}

	out.ShippingAddress = &sfenterprise.Address{
		Street:     in.Address,
		State:      in.State,
		PostalCode: in.PostalCode,
		City:       in.City,
		Country:    in.Country,
	}

	out.Website = in.WebAddress
	out.Etap_MultiObject_EtapRef__c = in.Ref
	return nil
}

func (i *io) manualTransformETAPGiftToSalesforceOpportunity(in *generated.Gift, out *sfenterprise.Opportunity) error {
	out.Etap_MigrationExplanation__c = ptr(fmt.Sprintf("This opportunity was generated from an eTapestry gift of $%f.", *in.Amount))
	if len(*out.Etap_MigrationExplanation__c) > 255 {
		return fmt.Errorf("migration explanation too long: %d > 255", len(*out.Etap_MigrationExplanation__c))
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.OwnerId = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}

	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Description = desc
	}

	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.ContactId = id
		aid, err := noReplacementForId(makerContact.AccountId)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.AccountId = aid
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.AccountId = id
	} else {
		return fmt.Errorf("could not find maker for gift: %q with ref %q", *in.Ref, *in.AccountRef)
	}
	out.Amount = in.Amount

	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.CloseDate = date
	}

	out.Etap_MultiObject_EtapRef__c = in.Ref
	out.StageName = ptr(sfenterprise.Opportunity_StageName_Received)
	out.Name = ptr(strings.TrimSpace(*in.Campaign + " Donation | " + out.CloseDate.ToGoTime().Format("01/02/2006")))
	out.Name = trimIfLongerThan(out.Name, 120)
	return nil
}

func (i *io) manualTransformETAPPaymentToSalesforcePayment(in *generated.Payment, out *sfenterprise.Npe01__OppPayment__c) error {
	explanation := fmt.Sprintf("This payment-on-opportunity was generated from an eTapestry payment-on-gift by %s on %s of $%f.", *in.AccountName, *in.Date, *in.Amount)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	if id, err := idPlaceholderForRef(in.PledgeRef); err != nil {
		return fmt.Errorf("creating placeholder for pledge: %w", err)
	} else {
		out.Npe01__Opportunity__c = id
	}
	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.Npe01__Payment_Date__c = date
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	out.Npe01__Payment_Amount__c = in.Amount
	out.Etap_MultiObject_EtapRef__c = in.Ref
	// Name is intentionally omitted. It is not allowed to be set.
	return nil
}

func (i *io) manualTransformETAPPledgeToSalesforceOpportunity(in *generated.Pledge, out *sfenterprise.Opportunity) error {
	explanation := fmt.Sprintf("This opportunity was generated from an eTapestry pledge by %s on %s of $%f.", *in.AccountName, *in.Date, *in.Amount)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.OwnerId = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Description = desc
	}
	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.ContactId = id
		aid, err := noReplacementForId(makerContact.AccountId)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.AccountId = aid
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.AccountId = id
	} else {
		return fmt.Errorf("could not find maker for gift: %q with ref %q", *in.Ref, *in.AccountRef)
	}

	out.Amount = in.Amount

	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.CloseDate = date
	}
	out.Etap_MultiObject_EtapRef__c = in.Ref

	out.Name = ptr(strings.TrimSpace(*in.Campaign + " Pledge | " + out.CloseDate.ToGoTime().Format("01/02/2006")))
	out.Name = trimIfLongerThan(out.Name, 120)

	out.StageName = ptr(sfenterprise.Opportunity_StageName_Received)

	return nil
}

func (i *io) manualTransformETAPRecurringGiftScheduleToSalesforceRecurringDonation(in *generated.RecurringGiftSchedule, out *sfenterprise.Npe03__Recurring_Donation__c) error {
	explanation := fmt.Sprintf("This recurring donation was generated from an eTapestry recurring gift schedule from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}

	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.Npe03__Date_Established__c = date
	}

	out.Name = in.Note
	if out.Name == nil || *out.Name == "" {
		out.Name = ptr(fmt.Sprintf("Recurring Donation by %s on %s", *in.AccountName, *in.Date))
	}
	out.Name = trimIfLongerThan(out.Name, 80)

	out.Npe03__Amount__c = in.Amount

	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.Npe03__Contact__c = id
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.Npe03__Organization__c = id
	} else {
		return fmt.Errorf("could not find maker for gift: %q with ref %q", *in.Ref, *in.AccountRef)
	}

	if in.Schedule.FirstInstallmentDate != nil && *in.Schedule.FirstInstallmentDate != "" {
		if date, err := AttemptToParseNilableDate(in.Schedule.FirstInstallmentDate); err != nil {
			return fmt.Errorf("first installment date: %w", err)
		} else {
			out.Npsp__StartDate__c = date
		}
	}
	if in.Schedule.StopDate != nil && *in.Schedule.StopDate != "" {
		if date, err := AttemptToParseNilableDate(in.Schedule.StopDate); err != nil {
			return fmt.Errorf("first installment date: %w", err)
		} else {
			out.Npsp__StartDate__c = date
		}
		out.Npsp__RecurringType__c = ptr(sfenterprise.Npe03RecurringDonation_npspRecurringType_Fixed)
	} else {
		out.Npsp__RecurringType__c = ptr(sfenterprise.Npe03RecurringDonation_npspRecurringType_Open)
	}
	ifreq, periodt, err := convertInstallmentFrequency(*in.Schedule.Frequency)
	if err != nil {
		return fmt.Errorf("converting installment frequency: %w", err)
	}
	out.Npsp__InstallmentFrequency__c = ifreq
	out.Npe03__Installment_Period__c = periodt
	out.Npe03__Amount__c = in.Schedule.InstallmentAmount

	// Default to the first of the month, if we can't find a better value
	out.Npsp__Day_of_Month__c = ptr(sfenterprise.Npe03RecurringDonation_npspDayofMonth_1)
	if date, err := AttemptToParseNilableDate(in.Schedule.FirstInstallmentDate); err == nil {
		_, _, dayOfMonth := date.ToGoTime().Date()
		dom, err := sfenterprise.Parse_Npe03RecurringDonation_npspDayofMonth_(fmt.Sprintf("%d", dayOfMonth))
		if err == nil {
			out.Npsp__Day_of_Month__c = &dom
		}
	}
	return nil
}

func (i *io) manualTransformETAPRecurringGiftToSalesforceOpportunity(in *generated.RecurringGift, out *sfenterprise.Opportunity) error {
	explanation := fmt.Sprintf("This opportunity was generated from an eTapestry recurring gift of $%f.", *in.Amount)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.OwnerId = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if desc, err := errIfLongerThan(in.Note, 32000); err != nil {
		return fmt.Errorf("note: %w", err)
	} else {
		out.Description = desc
	}
	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.ContactId = id
		aid, err := noReplacementForId(makerContact.AccountId)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.AccountId = aid
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.AccountId = id
	} else {
		return fmt.Errorf("could not find maker for gift: %q with ref %q", *in.Ref, *in.AccountRef)
	}

	if id, err := idPlaceholderForRef(in.RecurringGiftScheduleRef); err != nil {
		return fmt.Errorf("creating placeholder for recurring gift schedule: %w", err)
	} else {
		out.Npe03__Recurring_Donation__c = id
	}

	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.CloseDate = date
	}
	out.Name = ptr(strings.TrimSpace(*in.Campaign + " Recurring Donation | " + out.CloseDate.ToGoTime().Format("01/02/2006")))
	out.Name = trimIfLongerThan(out.Name, 120)
	out.Amount = in.Amount
	out.StageName = ptr(sfenterprise.Opportunity_StageName_Received)
	out.Etap_MultiObject_EtapRef__c = in.Ref
	return nil
}

func (i *io) manualTransformETAPSoftCreditToSalesforceAccountSoftCredit(in *generated.SoftCredit, out *sfenterprise.Npsp__Account_Soft_Credit__c) error {
	explanation := fmt.Sprintf("This partial soft credit was generated from an eTapestry soft credit with reference %s on %s.", *in.Ref, *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if id, err := idPlaceholderForRef(in.AccountRef); err != nil {
		return fmt.Errorf("creating placeholder for contact contact maker: %w", err)
	} else {
		out.Npsp__Account__c = id
	}
	out.Npsp__Amount__c = clonePtr(in.Amount)

	hcr := *in.HardCreditRef
	hardJE := i.in.JournalEntryRefs[hcr]
	if hardJE == nil {
		return &IsMissingHardCredit{msg: fmt.Sprintf("journal entry appears to be missing hard credit %q", hcr)}
	}
	if hardJE.RecurringGiftSchedule != nil {
		return &IsMissingHardCredit{msg: fmt.Sprintf("journal entry appears to be missing hard credit - it's a recurring gift schedule %q", hcr)}
	}
	if id, err := idPlaceholderForRef(in.HardCreditRef); err != nil {
		return fmt.Errorf("creating placeholder for contact contact maker: %w", err)
	} else {
		out.Npsp__Opportunity__c = id
	}
	// Name is intentionally omitted - it cannot be set
	return nil
}

func (i *io) manualTransformETAPSoftCreditToSalesforcePartialSoftCredit(in *generated.SoftCredit, out *sfenterprise.Npsp__Partial_Soft_Credit__c) error {
	explanation := fmt.Sprintf("This partial soft credit was generated from an eTapestry soft credit with reference %s on %s.", *in.Ref, *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if id, err := idPlaceholderForRef(in.AccountRef); err != nil {
		return fmt.Errorf("creating placeholder for contact contact maker: %w", err)
	} else {
		out.Npsp__Contact__c = id
	}
	out.Npsp__Amount__c = clonePtr(in.Amount)

	hcr := *in.HardCreditRef
	hardJE := i.in.JournalEntryRefs[hcr]
	if hardJE == nil {
		return &IsMissingHardCredit{msg: fmt.Sprintf("journal entry appears to be missing hard credit %q", hcr)}
	}
	if id, err := idPlaceholderForRef(&hcr); err != nil {
		return fmt.Errorf("creating placeholder for contact contact maker: %w", err)
	} else {
		if hardJE.RecurringGiftSchedule != nil {
			return &IsMissingHardCredit{msg: fmt.Sprintf("journal entry appears to be missing hard credit - it's a recurring gift schedule %q", hcr)}
		} else {
			out.Npsp__Opportunity__c = id
		}
	}
	// Name is intentionally omitted - it cannot be set.
	return nil
}

func (i *io) manualTransformETAPContactToSalesforceTask(in *generated.Contact, out *sfenterprise.Task) error {
	explanation := fmt.Sprintf("This task was generated from an eTapestry Contact from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.ActivityDate = date
	}
	/* https://ideas.salesforce.com/s/idea/a0B8W00000H5GHDUA3/ability-to-set-completeddatetime-audit-fields
	if date, err := AttemptToParseNilableDateTime(in.Date); err != nil {
		return fmt.Errorf("date time: %w", err)
	} else {
		out.CompletedDateTime = date
	}
	*/
	out.Etap_MultiObject_EtapRef__c = in.Ref
	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact contact maker: %w", err)
		}
		out.WhoId = id
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.WhatId = id
	} else {
		return fmt.Errorf("could not find maker for gift: %q with ref %q", *in.Ref, *in.AccountRef)
	}
	out.Description = trimIfLongerThan(in.Subject, 255)
	if in.Method != nil && *in.Method != "" {
		s, err := sfenterprise.Parse_Task_Subject_(*in.Method)
		if err != nil {
			return fmt.Errorf("parsing task subject %q: %w", *in.Method, err)
		}
		out.Subject = taskSubject(s, in.Subject, in.Attachments)
		t, err := sfenterprise.Parse_Task_Type_(*in.Method)
		if err != nil {
			return fmt.Errorf("parsing task type %q: %w", *in.Method, err)
		}
		out.Type = &t
	}
	out.Status = ptr(sfenterprise.Task_Status_Completed)
	out.Etap_Contact_Note__c = trimIfLongerThan(in.Note, 255)
	out.Etap_Contact_Attachments__c = trimIfLongerThan(out.Etap_Contact_Attachments__c, 255)
	id, err := idPlaceholderForRef(ptr(additionalContextPlaceholderRef(*in.Ref)))
	if err != nil {
		return fmt.Errorf("creating placeholder for additional context: %w", err)
	}
	out.Etap_AdditionalContextForRecord__c = id
	return nil
}

func taskSubject(ts sfenterprise.Task_Subject_, subject *string, attachments *generated.ArrayOfAttachment) *sfenterprise.Task_Subject_ {
	sub := ""
	if subject != nil && *subject != "" {
		sub = ": " + *subject
	}
	valueNoStar := string(ts) + sub
	if attachments != nil && len(attachments.Items) > 0 {
		valueNoStar = "*" + valueNoStar
	}
	return ptr(sfenterprise.Task_Subject_(*trimIfLongerThan(ptr(valueNoStar), 255)))
}

func (i *io) manualTransformETAPContactToSalesforceEtapAdditionalContext(in *generated.Contact, out *sfenterprise.Etap_AdditionalContext__c) error {
	explanation := fmt.Sprintf("This additional context was generated from an eTapestry Contact from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	out.Name = ptr(additionalContextPlaceholderRef(*in.Ref))
	out.Etap_MultiObject_EtapRef__c = ptr(additionalContextPlaceholderRef(*in.Ref))
	return nil
}

func (i *io) manualTransformETAPNoteToSalesforceTask(in *generated.Note, out *sfenterprise.Task) error {
	explanation := fmt.Sprintf("This task was generated from an eTapestry Note from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()

	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.ActivityDate = date
	}
	/* https://ideas.salesforce.com/s/idea/a0B8W00000H5GHDUA3/ability-to-set-completeddatetime-audit-fields
	if date, err := AttemptToParseNilableDateTime(in.Date); err != nil {
		return fmt.Errorf("date time: %w", err)
	} else {
		out.CompletedDateTime = date
	}
	// guaddothedoggo!#132
	*/
	out.Etap_MultiObject_EtapRef__c = in.Ref
	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact contact maker: %w", err)
		}
		out.WhoId = id
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.WhatId = id
	} else {
		return fmt.Errorf("could not find maker for note: %q with ref %q", *in.Ref, *in.AccountRef)
	}
	out.Status = ptr(sfenterprise.Task_Status_Completed)
	out.Subject = taskSubject(sfenterprise.Task_Subject_Note, in.Note, in.Attachments)
	out.Description = trimIfLongerThan(in.Note, 255)
	out.Etap_Note_Note__c = trimIfLongerThan(in.Note, 255)
	out.Etap_Note_Attachments__c = trimIfLongerThan(out.Etap_Note_Attachments__c, 255)
	id, err := idPlaceholderForRef(ptr(additionalContextPlaceholderRef(*in.Ref)))
	if err != nil {
		return fmt.Errorf("creating placeholder for additional context: %w", err)
	}
	out.Etap_AdditionalContextForRecord__c = id
	return nil
}

func (i *io) manualTransformETAPNoteToSalesforceEtapAdditionalContext(in *generated.Note, out *sfenterprise.Etap_AdditionalContext__c) error {
	explanation := fmt.Sprintf("This additional context was generated from an eTapestry Note from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()
	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	out.Name = ptr(additionalContextPlaceholderRef(*in.Ref))
	out.Etap_MultiObject_EtapRef__c = ptr(additionalContextPlaceholderRef(*in.Ref))
	return nil
}

func (i *io) manualTransformETAPAttachmentToSalesforceEtapAdditionalContext(in *generated.Attachment, out *sfenterprise.Etap_AdditionalContext__c) error {
	explanation := fmt.Sprintf("This additional context was generated from an eTapestry Attachment from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()
	if date, err := AttemptToParseNilableDateTime(in.Date); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
		out.LastModifiedDate = date
	}
	out.Name = ptr(additionalContextPlaceholderRef(*in.Ref))
	out.Etap_MultiObject_EtapRef__c = ptr(additionalContextPlaceholderRef(*in.Ref))
	return nil
}

func (i *io) manualTransformETAPAttachmentToSalesforceContentVersion(in *generated.Attachment, out *sfenterprise.ContentVersion) error {
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	if date, err := AttemptToParseNilableDateTime(in.Date); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
		out.LastModifiedDate = date
	}
	out.Title = in.Filename
	filePath, err := GetPathToAttachment(in)
	if err != nil {
		return fmt.Errorf("getting path to attachment: %w", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading attachment file: %w", err)
	}
	versionData := encodeToBase64(data)
	if len(versionData) == 0 {
		return fmt.Errorf("attachment file is empty")
	}
	out.VersionData = &versionData
	out.PathOnClient = in.Filename
	out.Etap_MultiObject_EtapRef__c = ptr(*in.Ref)
	// out.FirstPublishLocationId = &i.in.CallerUserId
	out.OwnerId = &i.in.AttributedUserId
	return nil
}

func encodeToBase64(data []byte) []byte {
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	encodedData := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encodedData, data)
	return encodedData
}

func (i *io) pureManualCreateContentDocumentLink(je *overrides.JournalEntry, a *generated.Attachment) (*sfenterprise.ContentDocumentLink, error) {
	out := &sfenterprise.ContentDocumentLink{}
	jeRef := je.Ref()
	if id, err := idPlaceholderForRef(&jeRef); err != nil {
		return nil, fmt.Errorf("creating placeholder for journal entry ref for attachment: %w", err)
	} else {
		out.LinkedEntityId = id
	}
	if id, err := idPlaceholderForRef(a.Ref); err != nil {
		return nil, fmt.Errorf("creating placeholder for attachment ref: %w", err)
	} else {
		out.ContentDocumentId = id
	}
	// https://developer.salesforce.com/docs/atlas.en-us.object_reference.meta/object_reference/sforce_api_objects_contentdocumentlink.htm
	// Means "all users who have the permission to see the file (i.e. shared via record)"
	out.Visibility = ptr(sfenterprise.ContentDocumentLink_Visibility_AllUsers)
	// Means "inferred" - content is shared with users who can see the linked record
	out.ShareType = ptr(sfenterprise.ContentDocumentLink_ShareType_I)
	return out, nil
}

func (i *io) manualTransformETAPSegmentedDonationToSalesforceOpportunity(in *generated.SegmentedDonation, out *sfenterprise.Opportunity) error {
	explanation := fmt.Sprintf("This additional context was generated from an eTapestry Segmented Donation from %s.", *in.Date)
	if exp, err := errIfLongerThan(&explanation, 255); err != nil {
		return fmt.Errorf("explanation: %w", err)
	} else {
		out.Etap_MigrationExplanation__c = exp
	}
	out.CreatedById = &i.in.AttributedUserId
	out.LastModifiedById = &i.in.AttributedUserId
	out.OwnerId = &i.in.AttributedUserId
	out.Etap_MigrationTime__c = NowXSD()
	if date, err := AttemptToParseNilableDateTime(in.CreatedDate); err != nil {
		return fmt.Errorf("created date: %w", err)
	} else {
		out.CreatedDate = date
	}
	if date, err := AttemptToParseNilableDateTime(in.LastModifiedDate); err != nil {
		return fmt.Errorf("last modified date: %w", err)
	} else {
		out.LastModifiedDate = date
	}
	if date, err := AttemptToParseNilableDate(in.Date); err != nil {
		return fmt.Errorf("date: %w", err)
	} else {
		out.CloseDate = date
	}
	out.Etap_MultiObject_EtapRef__c = ptr(*in.Ref)
	if makerContact, ok := i.out.contactsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerContact.Etap_Account_Ref__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.ContactId = id
		aid, err := noReplacementForId(makerContact.AccountId)
		if err != nil {
			return fmt.Errorf("creating placeholder for contact gift maker: %w", err)
		}
		out.AccountId = aid
	} else if makerAccount, ok := i.out.accountsByRefs[*in.AccountRef]; ok {
		id, err := idPlaceholderForRef(makerAccount.Etap_MultiObject_EtapRef__c)
		if err != nil {
			return fmt.Errorf("creating placeholder for account gift maker: %w", err)
		}
		out.AccountId = id
	} else {
		return fmt.Errorf("could not find maker for gift: %q with ref %q", *in.Ref, *in.AccountRef)
	}

	out.Amount = clonePtr(in.TotalAmount)
	out.StageName = ptr(sfenterprise.Opportunity_StageName_Received)
	out.Name = ptr(strings.TrimSpace("Segmented Donation | " + out.CloseDate.ToGoTime().Format("01/02/2006")))
	out.Name = trimIfLongerThan(out.Name, 120)
	return nil
}

func (i *io) manualTransformETAPDisbursementToSalesforceOpportunity(in *generated.Disbursement, result *sfenterprise.Opportunity) error {
	return fmt.Errorf("not supported")
}

func (i *io) manualTransformETAPAttachmentToSalesforceContentDocumentLink(in *generated.Attachment, out *sfenterprise.ContentDocumentLink) error {
	return fmt.Errorf("not supported - this conversion requires a journal entry to be present to do the conversion - this should not be called.")
}

func campaignPlaceholderRef(campaign string) string {
	return "campaign-" + campaign
}
func approachPlaceholderRef(approach string) string {
	return "approach-" + approach
}
func additionalContextPlaceholderRef(ref string) string {
	return "Overflow_Information_For_" + ref
}

func GetPathToAttachment(a *generated.Attachment) (string, error) {
	if a.Ref == nil {
		return "", fmt.Errorf("attachment ref is nil")
	}
	if a.Filename == nil {
		return "", fmt.Errorf("attachment filename is nil")
	}
	return fmt.Sprintf("%s/%s/%s", exportfiles.AttachmentsFolder, *a.Ref, *sanitizeFilename(a.Filename)), nil
}

func sanitizeFilename(s *string) *string {
	if s == nil {
		return nil
	}
	var builder strings.Builder
	for _, runeValue := range *s {
		if runeValue != '$' {
			builder.WriteRune(runeValue)
		}
	}
	ss := builder.String()
	return &ss
}

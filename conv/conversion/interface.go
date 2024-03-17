package conversion

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/etap/data"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
	"github.com/Silicon-Ally/etap2sf/etap/inference/customfields"
	"github.com/Silicon-Ally/etap2sf/salesforce/clients/enterprise/utils"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
)

//nolint:unused // Used in files after step 12.
type io struct {
	in  *Input
	out *Output
}

type Input struct {
	CustomFields                  *customfields.CustomFields
	Accounts                      []*generated.Account
	Approaches                    []string
	Campaigns                     []string
	Funds                         []*generated.Fund
	JournalEntries                []*overrides.JournalEntry
	Relationships                 []*generated.Relationship
	JournalEntryRefs              map[string]*overrides.JournalEntry
	AttributedUserId              sfenterprise.ID
	CallerUserId                  sfenterprise.ID
	OrganizationAccountRecordType sfenterprise.ID
	HouseholdAccountRecordType    sfenterprise.ID
}

type Output struct {
	Accounts               []*sfenterprise.Account
	AccountSoftCredits     []*sfenterprise.Npsp__Account_Soft_Credit__c
	AdditionalContexts     []*sfenterprise.Etap_AdditionalContext__c
	Affiliations           []*sfenterprise.Npe5__Affiliation__c
	Campaigns              []*sfenterprise.Campaign
	Contacts               []*sfenterprise.Contact
	ContentDocumentLinks   []*sfenterprise.ContentDocumentLink
	ContentVersions        []*sfenterprise.ContentVersion
	ContentNotes           []*sfenterprise.ContentNote
	GeneralAccountingUnits []*sfenterprise.Npsp__General_Accounting_Unit__c
	GAUAllocations         []*sfenterprise.Npsp__Allocation__c
	Households             []*sfenterprise.Npo02__Household__c
	Opportunities          []*sfenterprise.Opportunity
	Payments               []*sfenterprise.Npe01__OppPayment__c
	PartialSoftCredits     []*sfenterprise.Npsp__Partial_Soft_Credit__c
	RecurringDonations     []*sfenterprise.Npe03__Recurring_Donation__c
	Relationships          []*sfenterprise.Npe4__Relationship__c
	Tasks                  []*sfenterprise.Task

	refSubstitutions map[string]string                //nolint:unused // Used in files after step 12.
	accountsByRefs   map[string]*sfenterprise.Account //nolint:unused // Used in files after step 12.
	contactsByRefs   map[string]*sfenterprise.Contact //nolint:unused // Used in files after step 12.
}

func GetInput() (*Input, error) {
	client, err := utils.NewSandboxClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get sandbox client: %v", err)
	}
	attributedUserID, err := client.LookupUserByEmail(conversionsettings.AttributedUserEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup attributed user: %v", err)
	}
	callerUserID, err := client.LookupUserByEmail(conversionsettings.CallerUserEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup attributed user: %v", err)
	}
	orgRTID, hhRTID, err := client.GetAccountRecordTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get record types: %v", err)
	}
	if orgRTID == "" {
		return nil, fmt.Errorf("failed to get organization record type")
	}
	if hhRTID == "" {
		return nil, fmt.Errorf("failed to get household record type")
	}
	accounts, err := data.GetAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}
	approaches, err := data.GetApproaches()
	if err != nil {
		return nil, fmt.Errorf("failed to get approaches: %v", err)
	}
	campaigns, err := data.GetCampaigns()
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %v", err)
	}
	customFields, err := customfields.GetCustomFields()
	if err != nil {
		return nil, fmt.Errorf("failed to get custom fields: %v", err)
	}
	funds, err := data.GetFunds()
	if err != nil {
		return nil, fmt.Errorf("failed to get funds: %v", err)
	}
	journalEntries, err := data.GetJournalEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entries: %v", err)
	}
	relationships, err := data.GetRelationships()
	if err != nil {
		return nil, fmt.Errorf("failed to get relationships: %v", err)
	}
	jes := make(map[string]*overrides.JournalEntry)
	for _, je := range journalEntries {
		ref := je.Ref()
		if jes[ref] != nil {
			return nil, fmt.Errorf("duplicate jounal entry ref: %q", ref)
		}
		jes[ref] = je
	}

	return &Input{
		Accounts:                      accounts,
		Approaches:                    approaches,
		Campaigns:                     campaigns,
		CustomFields:                  customFields,
		Funds:                         funds,
		JournalEntries:                journalEntries,
		Relationships:                 relationships,
		JournalEntryRefs:              jes,
		AttributedUserId:              attributedUserID,
		CallerUserId:                  callerUserID,
		OrganizationAccountRecordType: orgRTID,
		HouseholdAccountRecordType:    hhRTID,
	}, nil
}

package download_all_data_from_etap

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/data"
	"github.com/Silicon-Ally/etap2sf/etap/inference/customfields"
)

func Run() error {
	accounts, err := data.GetAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %v", err)
	}
	approaches, err := data.GetApproaches()
	if err != nil {
		return fmt.Errorf("failed to get approaches: %v", err)
	}
	campaigns, err := data.GetCampaigns()
	if err != nil {
		return fmt.Errorf("failed to get campaigns: %v", err)
	}
	definedFields, err := data.GetDefinedFields()
	if err != nil {
		return fmt.Errorf("failed to get defined fields: %v", err)
	}
	funds, err := data.GetFunds()
	if err != nil {
		return fmt.Errorf("failed to get funds: %v", err)
	}
	journalEntries, err := data.GetJournalEntries()
	if err != nil {
		return fmt.Errorf("failed to get journal entries: %v", err)
	}
	relationships, err := data.GetRelationships()
	if err != nil {
		return fmt.Errorf("failed to get relationships: %v", err)
	}
	customFields, err := customfields.GetCustomFields()
	if err != nil {
		return fmt.Errorf("failed to get custom fields: %v", err)
	}
	fmt.Printf(`
Download complete.

Found %d Accounts
Found %d Approaches
Found %d Campaigns
Found %d Defined Fields
Found %d Funds
Found %d Journal Entries
Found %d Relationships
Found %d Custom Fields

Your metadata has successfully been downloaded from eTapestry. You may proceed to the next step.
`, len(accounts), len(approaches), len(campaigns),
		len(definedFields), len(funds), len(journalEntries),
		len(relationships), len(customFields.Fields))
	return nil
}

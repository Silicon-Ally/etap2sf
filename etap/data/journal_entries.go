package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var journalEntries []*overrides.JournalEntry

func GetJournalEntries() ([]*overrides.JournalEntry, error) {
	if journalEntries != nil {
		return journalEntries, nil
	}
	data, err := utils.MemoizeOperation("etap-journal-entries.json", doGetJournalEntryData)
	if err != nil {
		return nil, fmt.Errorf("failed to get je data: %v", err)
	}
	result := []*overrides.JournalEntry{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal je data: %v", err)
	}
	refs := map[string]bool{}
	for _, je := range result {
		ref := je.Ref()
		if ref == "" {
			return nil, fmt.Errorf("je has no ref: %+v", je)
		}
		if _, ok := refs[ref]; ok {
			return nil, fmt.Errorf("duplicate je ref: %s", ref)
		}
		refs[ref] = true
	}
	journalEntries = result
	return result, nil
}

func doGetJournalEntryData() ([]byte, error) {
	accounts, err := GetAccounts()
	if err != nil {
		return nil, fmt.Errorf("getting accounts: %w", err)
	}

	return client.WithClient(func(c *client.Client) ([]byte, error) {
		getJournalEntries := func(account *generated.Account) ([]*overrides.JournalEntry, error) {
			jeData, err := utils.MemoizeOperation(fmt.Sprintf("journal-entries/%s.json", *account.Ref), func() ([]byte, error) {
				journalEntries, err := c.GetAllJournalEntries(account)
				if err != nil {
					return nil, fmt.Errorf("getting journal entries: %w", err)
				}
				data, err := json.MarshalIndent(journalEntries, "", "  ")
				if err != nil {
					return nil, fmt.Errorf("marshaling journal entries: %w", err)
				}
				// Prevents us from destroying the server
				time.Sleep(500 * time.Millisecond)
				return data, nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get journal entries: %v", err)
			}
			var journalEntries []*overrides.JournalEntry
			err = json.Unmarshal(jeData, &journalEntries)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal journal entries: %v", err)
			}
			return journalEntries, nil
		}

		jes := []*overrides.JournalEntry{}
		for i, account := range accounts {
			journalEntries, err := getJournalEntries(account)
			if err != nil {
				return nil, fmt.Errorf("error getting journal entries for account %s: %v", *account.Ref, err)
			}
			for _, je := range journalEntries {
				if !je.Empty && je.Ref() == "" {
					return nil, fmt.Errorf("early - journal entry for account %s has no ref: %+v", *account.Ref, je)
				}
				if !je.Empty {
					jes = append(jes, je)
				}
			}
			fmt.Printf("Account %d/%d: %s\n", i+1, len(accounts), *account.Name)
		}

		result, err := json.MarshalIndent(jes, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal journal entries: %v", err)
		}
		return result, nil
	})
}

func GetContacts() ([]*generated.Contact, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.Contact { return je.Contact })
}

func GetNotes() ([]*generated.Note, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.Note { return je.Note })
}

func GetGifts() ([]*generated.Gift, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.Gift { return je.Gift })
}

func GetPayments() ([]*generated.Payment, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.Payment { return je.Payment })
}

func GetSoftCredits() ([]*generated.SoftCredit, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.SoftCredit { return je.SoftCredit })
}

func GetRecurringGiftSchedules() ([]*generated.RecurringGiftSchedule, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.RecurringGiftSchedule { return je.RecurringGiftSchedule })
}

func GetRecurringGifts() ([]*generated.RecurringGift, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.RecurringGift { return je.RecurringGift })
}

func GetPledges() ([]*generated.Pledge, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.Pledge { return je.Pledge })
}

func GetDisbursements() ([]*generated.Disbursement, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.Disbursement { return je.Disbursement })
}

func GetSegmentedDonations() ([]*generated.SegmentedDonation, error) {
	return getJEs(func(je *overrides.JournalEntry) *generated.SegmentedDonation { return je.SegmentedDonation })
}

func getJEs[T any](fn func(je *overrides.JournalEntry) *T) ([]*T, error) {
	jes, err := GetJournalEntries()
	if err != nil {
		return nil, fmt.Errorf("getting journal entries: %w", err)
	}
	result := []*T{}
	for _, je := range jes {
		t := fn(je)
		if t != nil {
			result = append(result, t)
		}
	}
	return result, nil
}

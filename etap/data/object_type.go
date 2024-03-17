package data

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap"
)

func Get(o etap.ObjectType) ([]any, error) {
	switch o {
	case etap.ObjectType_Account:
		accounts, err := GetAccounts()
		if err != nil {
			return nil, fmt.Errorf("getting accounts: %w", err)
		}
		return asAny(accounts), nil
	case etap.ObjectType_Attachment:
		attachments, err := GetAttachments()
		if err != nil {
			return nil, fmt.Errorf("getting attachments: %w", err)
		}
		return asAny(attachments), nil
	case etap.ObjectType_Relationship:
		relationships, err := GetRelationships()
		if err != nil {
			return nil, fmt.Errorf("getting relationships: %w", err)
		}
		return asAny(relationships), nil
	case etap.ObjectType_Approach:
		approaches, err := GetApproaches()
		if err != nil {
			return nil, fmt.Errorf("getting approaches: %w", err)
		}
		return asAny(approaches), nil
	case etap.ObjectType_Campaign:
		campaigns, err := GetCampaigns()
		if err != nil {
			return nil, fmt.Errorf("getting campaigns: %w", err)
		}
		return asAny(campaigns), nil
	case etap.ObjectType_Fund:
		funds, err := GetFunds()
		if err != nil {
			return nil, fmt.Errorf("getting funds: %w", err)
		}
		return asAny(funds), nil
	case etap.ObjectType_Contact:
		data, err := GetContacts()
		if err != nil {
			return nil, fmt.Errorf("getting contacts: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_Note:
		data, err := GetNotes()
		if err != nil {
			return nil, fmt.Errorf("getting notes: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_Gift:
		data, err := GetGifts()
		if err != nil {
			return nil, fmt.Errorf("getting gifts: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_Payment:
		data, err := GetPayments()
		if err != nil {
			return nil, fmt.Errorf("getting payments: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_SoftCredit:
		data, err := GetSoftCredits()
		if err != nil {
			return nil, fmt.Errorf("getting softcredits: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_RecurringGiftSchedule:
		data, err := GetRecurringGiftSchedules()
		if err != nil {
			return nil, fmt.Errorf("getting recurringgiftschedules: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_RecurringGift:
		data, err := GetRecurringGifts()
		if err != nil {
			return nil, fmt.Errorf("getting recurringgifts: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_Pledge:
		data, err := GetPledges()
		if err != nil {
			return nil, fmt.Errorf("getting pledges: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_Disbursement:
		data, err := GetDisbursements()
		if err != nil {
			return nil, fmt.Errorf("getting disbursements: %w", err)
		}
		return asAny(data), nil
	case etap.ObjectType_SegmentedDonation:
		data, err := GetSegmentedDonations()
		if err != nil {
			return nil, fmt.Errorf("getting segmenteddonations: %w", err)
		}
		return asAny(data), nil
	}
	return nil, fmt.Errorf("unknown object type get-all-data: %q", o)
}

func asAny[T any](ts []T) []any {
	result := make([]any, len(ts))
	for i, t := range ts {
		result[i] = t
	}
	return result
}

package etap

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
)

type ObjectType string

const (
	ObjectType_Account               ObjectType = "Account"
	ObjectType_Approach              ObjectType = "Approach"
	ObjectType_Attachment            ObjectType = "Attachment"
	ObjectType_Campaign              ObjectType = "Campaign"
	ObjectType_Contact               ObjectType = "Contact"
	ObjectType_Disbursement          ObjectType = "Disbursement"
	ObjectType_Fund                  ObjectType = "Fund"
	ObjectType_Gift                  ObjectType = "Gift"
	ObjectType_Note                  ObjectType = "Note"
	ObjectType_Payment               ObjectType = "Payment"
	ObjectType_Pledge                ObjectType = "Pledge"
	ObjectType_Purchase              ObjectType = "Purchase"
	ObjectType_RecurringGift         ObjectType = "RecurringGift"
	ObjectType_RecurringGiftSchedule ObjectType = "RecurringGiftSchedule"
	ObjectType_Relationship          ObjectType = "Relationship"
	ObjectType_SegmentedDonation     ObjectType = "SegmentedDonation"
	ObjectType_SoftCredit            ObjectType = "SoftCredit"
)

func (o ObjectType) String() string {
	return string(o)
}

func (o ObjectType) IsString() bool {
	if o == ObjectType_Campaign || o == ObjectType_Approach {
		return true
	}
	return false
}

func (o ObjectType) Struct() (any, error) {
	switch o {
	case ObjectType_Account:
		return generated.Account{}, nil
	case ObjectType_Attachment:
		return generated.Attachment{}, nil
	case ObjectType_Relationship:
		return generated.Relationship{}, nil
	case ObjectType_Campaign:
		return "", fmt.Errorf("campaign is not a struct")
	case ObjectType_Approach:
		return "", fmt.Errorf("approach is not a struct")
	case ObjectType_Fund:
		return generated.Fund{}, nil
	case ObjectType_Contact:
		return generated.Contact{}, nil
	case ObjectType_Note:
		return generated.Note{}, nil
	case ObjectType_Gift:
		return generated.Gift{}, nil
	case ObjectType_Payment:
		return generated.Payment{}, nil
	case ObjectType_SoftCredit:
		return generated.SoftCredit{}, nil
	case ObjectType_RecurringGiftSchedule:
		return generated.RecurringGiftSchedule{}, nil
	case ObjectType_RecurringGift:
		return generated.RecurringGift{}, nil
	case ObjectType_Pledge:
		return generated.Pledge{}, nil
	case ObjectType_Disbursement:
		return generated.Disbursement{}, nil
	case ObjectType_SegmentedDonation:
		return generated.SegmentedDonation{}, nil
	default:
		return nil, fmt.Errorf("unknown object type etap-struct: %s", o)
	}
}

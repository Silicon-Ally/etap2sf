package overrides

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
)

// The underlying collection type in the generated code is typeless, so it can't instantiate
// it's unknown-typed children. This is a workaround to allow the generated code to instantiate
// a generic collection to instantiate its children with their correct type.
type Collection[T any] struct {
	Items []T `xml:"item,omitempty" json:"item,omitempty" yaml:"item,omitempty"`
}

type PagedQueryResultsResponse[T any] struct {
	Count     *int           `xml:"count,omitempty" json:"count,omitempty" yaml:"count,omitempty"`
	Data      *Collection[T] `xml:"data,omitempty" json:"data,omitempty" yaml:"data,omitempty"`
	Pages     *int           `xml:"pages,omitempty" json:"pages,omitempty" yaml:"pages,omitempty"`
	Start     *int           `xml:"start,omitempty" json:"start,omitempty" yaml:"start,omitempty"`
	Total     *int           `xml:"total,omitempty" json:"total,omitempty" yaml:"total,omitempty"`
	UsedCache *bool          `xml:"usedCache,omitempty" json:"usedCache,omitempty" yaml:"usedCache,omitempty"`
}

type OperationMessagingService_getExistingQueryResultsResponse[T any] struct {
	Result *PagedQueryResultsResponse[T] `xml:"result,omitempty" json:"result,omitempty" yaml:"result,omitempty"`
}

type AccountBody struct {
	M OperationMessagingService_getExistingQueryResultsResponse[generated.Account] `xml:"getExistingQueryResultsResponse"`
}

type JournalEntry struct {
	Empty bool
	// Used In Our Data
	Contact               *generated.Contact               `xml:"contact,omitempty" json:"contact,omitempty" yaml:"contact,omitempty"`
	Disbursement          *generated.Disbursement          `xml:"disbursement,omitempty" json:"disbursement,omitempty" yaml:"disbursement,omitempty"`
	Gift                  *generated.Gift                  `xml:"gift,omitempty" json:"gift,omitempty" yaml:"gift,omitempty"`
	Note                  *generated.Note                  `xml:"note,omitempty" json:"note,omitempty" yaml:"note,omitempty"`
	Payment               *generated.Payment               `xml:"payment,omitempty" json:"payment,omitempty" yaml:"payment,omitempty"`
	Pledge                *generated.Pledge                `xml:"pledge,omitempty" json:"pledge,omitempty" yaml:"pledge,omitempty"`
	RecurringGift         *generated.RecurringGift         `xml:"recurringGift,omitempty" json:"recurringGift,omitempty" yaml:"recurringGift,omitempty"`
	RecurringGiftSchedule *generated.RecurringGiftSchedule `xml:"recurringGiftSchedule,omitempty" json:"recurringGiftSchedule,omitempty" yaml:"recurringGiftSchedule,omitempty"`
	SegmentedDonation     *generated.SegmentedDonation     `xml:"segmentedDonation,omitempty" json:"segmentedDonation,omitempty" yaml:"segmentedDonation,omitempty"`
	SoftCredit            *generated.SoftCredit            `xml:"softCredit,omitempty" json:"softCredit,omitempty" yaml:"softCredit,omitempty"`
	// Not used in our our data
	Declaration     *generated.Declaration     `xml:"declaration,omitempty" json:"declaration,omitempty" yaml:"declaration,omitempty"`
	SegmentedPledge *generated.SegmentedPledge `xml:"segmentedPledge,omitempty" json:"segmentedPledge,omitempty" yaml:"segmentedPledge,omitempty"`
	Invitation      *generated.Invitation      `xml:"invitation,omitempty" json:"invitation,omitempty" yaml:"invitation,omitempty"`
	Purchase        *generated.Purchase        `xml:"purchase,omitempty" json:"purchase,omitempty" yaml:"purchase,omitempty"`
	SegmentedOrder  *generated.SegmentedOrder  `xml:"segmentedOrder,omitempty" json:"segmentedOrder,omitempty" yaml:"segmentedOrder,omitempty"`
	// Calendar isn't specified in the WSDL, though it is specified in the documentation
	// Calendar              *generated.Calendar              `xml:"calendar,omitempty" json:"calendar,omitempty" yaml:"calendar,omitempty"`
	// Participation isn't specified in the WSDL, though it is specified in the documentation.
	// Participation         *generated.Participation         `xml:"participation,omitempty" json:"participation,omitempty" yaml:"participation,omitempty"`
}

func (j *JournalEntry) Ref() string {
	if j.Note != nil && j.Note.Ref != nil {
		return *j.Note.Ref
	}
	if j.Contact != nil && j.Contact.Ref != nil {
		return *j.Contact.Ref
	}
	if j.Declaration != nil && j.Declaration.Ref != nil {
		return *j.Declaration.Ref
	}
	if j.Gift != nil && j.Gift.Ref != nil {
		return *j.Gift.Ref
	}
	if j.Pledge != nil && j.Pledge.Ref != nil {
		return *j.Pledge.Ref
	}
	if j.Payment != nil && j.Payment.Ref != nil {
		return *j.Payment.Ref
	}
	if j.RecurringGiftSchedule != nil && j.RecurringGiftSchedule.Ref != nil {
		return *j.RecurringGiftSchedule.Ref
	}
	if j.RecurringGift != nil && j.RecurringGift.Ref != nil {
		return *j.RecurringGift.Ref
	}
	if j.Disbursement != nil && j.Disbursement.Ref != nil {
		return *j.Disbursement.Ref
	}
	if j.Purchase != nil && j.Purchase.Ref != nil {
		return *j.Purchase.Ref
	}
	if j.SoftCredit != nil && j.SoftCredit.Ref != nil {
		return *j.SoftCredit.Ref
	}
	if j.SegmentedDonation != nil && j.SegmentedDonation.Ref != nil {
		return *j.SegmentedDonation.Ref
	}
	return ""
}

func (j *JournalEntry) Attachments() []*generated.Attachment {
	result := []*generated.Attachment{}
	if j.Note != nil && j.Note.Attachments != nil {
		result = append(result, j.Note.Attachments.Items...)
	}
	if j.Contact != nil && j.Contact.Attachments != nil {
		result = append(result, j.Contact.Attachments.Items...)
	}
	if j.Declaration != nil && j.Declaration.Attachments != nil {
		result = append(result, j.Declaration.Attachments.Items...)
	}
	if j.Gift != nil && j.Gift.Attachments != nil {
		result = append(result, j.Gift.Attachments.Items...)
	}
	if j.Pledge != nil && j.Pledge.Attachments != nil {
		result = append(result, j.Pledge.Attachments.Items...)
	}
	if j.Payment != nil && j.Payment.Attachments != nil {
		result = append(result, j.Payment.Attachments.Items...)
	}
	if j.RecurringGiftSchedule != nil && j.RecurringGiftSchedule.Attachments != nil {
		result = append(result, j.RecurringGiftSchedule.Attachments.Items...)
	}
	if j.RecurringGift != nil && j.RecurringGift.Attachments != nil {
		result = append(result, j.RecurringGift.Attachments.Items...)
	}
	if j.Disbursement != nil && j.Disbursement.Attachments != nil {
		result = append(result, j.Disbursement.Attachments.Items...)
	}
	if j.Purchase != nil && j.Purchase.Attachments != nil {
		result = append(result, j.Purchase.Attachments.Items...)
	}
	return result
}

func (j *JournalEntry) AccountRef() string {
	if j.Note != nil && j.Note.AccountRef != nil {
		return *j.Note.AccountRef
	}
	if j.Contact != nil && j.Contact.AccountRef != nil {
		return *j.Contact.AccountRef
	}
	if j.Declaration != nil && j.Declaration.AccountRef != nil {
		return *j.Declaration.AccountRef
	}
	if j.Gift != nil && j.Gift.AccountRef != nil {
		return *j.Gift.AccountRef
	}
	if j.Pledge != nil && j.Pledge.AccountRef != nil {
		return *j.Pledge.AccountRef
	}
	if j.Payment != nil && j.Payment.AccountRef != nil {
		return *j.Payment.AccountRef
	}
	if j.RecurringGiftSchedule != nil && j.RecurringGiftSchedule.AccountRef != nil {
		return *j.RecurringGiftSchedule.AccountRef
	}
	if j.RecurringGift != nil && j.RecurringGift.AccountRef != nil {
		return *j.RecurringGift.AccountRef
	}
	if j.Disbursement != nil && j.Disbursement.AccountRef != nil {
		return *j.Disbursement.AccountRef
	}
	if j.Purchase != nil && j.Purchase.AccountRef != nil {
		return *j.Purchase.AccountRef
	}
	if j.SegmentedDonation != nil && j.SegmentedDonation.AccountRef != nil {
		return *j.SegmentedDonation.AccountRef
	}
	return ""
}

func (j *JournalEntry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data, err := innerXML(d, start)
	if err != nil {
		return fmt.Errorf("when trying to read inner XML: %w", err)
	}

	var structType struct {
		Type int `xml:"type"`
	}
	if err := xml.Unmarshal(data, &structType); err != nil {
		return err
	}

	switch structType.Type {
	case 1:
		j.Note = &generated.Note{}
		return xml.Unmarshal(data, j.Note)
	case 2:
		j.Contact = &generated.Contact{}
		return xml.Unmarshal(data, j.Contact)
	case 4:
		j.Declaration = &generated.Declaration{}
		return xml.Unmarshal(data, j.Declaration)
	case 5:
		j.Gift = &generated.Gift{}
		return xml.Unmarshal(data, j.Gift)
	case 6:
		j.Pledge = &generated.Pledge{}
		return xml.Unmarshal(data, j.Pledge)
	case 7:
		j.Payment = &generated.Payment{}
		return xml.Unmarshal(data, j.Payment)
	case 8:
		j.RecurringGiftSchedule = &generated.RecurringGiftSchedule{}
		return xml.Unmarshal(data, j.RecurringGiftSchedule)
	case 9:
		j.RecurringGift = &generated.RecurringGift{}
		return xml.Unmarshal(data, j.RecurringGift)
	case 10:
		j.SegmentedDonation = &generated.SegmentedDonation{}
		return xml.Unmarshal(data, j.SegmentedDonation)
	case 11:
		j.SoftCredit = &generated.SoftCredit{}
		return xml.Unmarshal(data, j.SoftCredit)
	case 12:
		j.Disbursement = &generated.Disbursement{}
		return xml.Unmarshal(data, j.Disbursement)
	case 13:
		j.SegmentedPledge = &generated.SegmentedPledge{}
		return xml.Unmarshal(data, j.SegmentedPledge)
	case 14:
		j.Invitation = &generated.Invitation{}
		return xml.Unmarshal(data, j.Invitation)
	case 15:
		j.Purchase = &generated.Purchase{}
		return xml.Unmarshal(data, j.Purchase)
	case 18:
		j.SegmentedOrder = &generated.SegmentedOrder{}
		return xml.Unmarshal(data, j.SegmentedOrder)
	case 0:
		j.Empty = true
		return nil // This is a blank journal entry, so we can just ignore it.
	case 3:
		return fmt.Errorf("unsupported journal entry type (calendar, type=3)")
	case 19:
		return fmt.Errorf("unsupported journal entry type (participation, type=19)")
	default:
		return fmt.Errorf("unsupported journal entry type (unknown, type=%d)", structType.Type)
	}
}

func innerXML(d *xml.Decoder, start xml.StartElement) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	// Write the starting element to the buffer.
	if err := enc.EncodeToken(start); err != nil {
		return nil, err
	}

	depth := 1 // Start with depth of 1 to consider the starting element.
	for depth > 0 {
		token, err := d.Token()
		if err != nil {
			return nil, err
		}

		switch token.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}

		if err := enc.EncodeToken(token); err != nil {
			return nil, err
		}
	}

	if err := enc.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type PagedJournalEntriesResponse[T any] struct {
	Count     *int           `xml:"count,omitempty" json:"count,omitempty" yaml:"count,omitempty"`
	Data      *Collection[T] `xml:"data,omitempty" json:"data,omitempty" yaml:"data,omitempty"`
	Pages     *int           `xml:"pages,omitempty" json:"pages,omitempty" yaml:"pages,omitempty"`
	Start     *int           `xml:"start,omitempty" json:"start,omitempty" yaml:"start,omitempty"`
	Total     *int           `xml:"total,omitempty" json:"total,omitempty" yaml:"total,omitempty"`
	UsedCache *bool          `xml:"usedCache,omitempty" json:"usedCache,omitempty" yaml:"usedCache,omitempty"`
}

type OperationMessagingService_getJournalEntriesResponse[T any] struct {
	Result *PagedJournalEntriesResponse[T] `xml:"result,omitempty" json:"result,omitempty" yaml:"result,omitempty"`
}

type JournalEntryBody struct {
	M OperationMessagingService_getJournalEntriesResponse[JournalEntry] `xml:"getJournalEntriesResponse"`
}

type PagedRelationshipsResponse struct {
	Count     *int                                `xml:"count,omitempty" json:"count,omitempty" yaml:"count,omitempty"`
	Data      *Collection[generated.Relationship] `xml:"data,omitempty" json:"data,omitempty" yaml:"data,omitempty"`
	Pages     *int                                `xml:"pages,omitempty" json:"pages,omitempty" yaml:"pages,omitempty"`
	Start     *int                                `xml:"start,omitempty" json:"start,omitempty" yaml:"start,omitempty"`
	Total     *int                                `xml:"total,omitempty" json:"total,omitempty" yaml:"total,omitempty"`
	UsedCache *bool                               `xml:"usedCache,omitempty" json:"usedCache,omitempty" yaml:"usedCache,omitempty"`
}

type OperationMessagingService_getRelationships struct {
	Result *PagedRelationshipsResponse `xml:"result,omitempty" json:"result,omitempty" yaml:"result,omitempty"`
}

type RelationshipsBody struct {
	M OperationMessagingService_getRelationships `xml:"getRelationshipsResponse"`
}

type PagedDefinedFieldsResponse struct {
	Count     *int                                `xml:"count,omitempty" json:"count,omitempty" yaml:"count,omitempty"`
	Data      *Collection[generated.DefinedField] `xml:"data,omitempty" json:"data,omitempty" yaml:"data,omitempty"`
	Pages     *int                                `xml:"pages,omitempty" json:"pages,omitempty" yaml:"pages,omitempty"`
	Start     *int                                `xml:"start,omitempty" json:"start,omitempty" yaml:"start,omitempty"`
	Total     *int                                `xml:"total,omitempty" json:"total,omitempty" yaml:"total,omitempty"`
	UsedCache *bool                               `xml:"usedCache,omitempty" json:"usedCache,omitempty" yaml:"usedCache,omitempty"`
}

type OperationMessagingService_getDefinedFields struct {
	Result *PagedDefinedFieldsResponse `xml:"result,omitempty" json:"result,omitempty" yaml:"result,omitempty"`
}

type DefinedFieldsBody struct {
	M OperationMessagingService_getDefinedFields `xml:"getDefinedFieldsResponse"`
}

type FundObjectsBody struct {
	M generated.OperationMessagingService_getFundObjectsResponse `xml:"getFundObjectsResponse"`
}

type CampaignsBody struct {
	M generated.OperationMessagingService_getCampaignsResponse `xml:"getCampaignsResponse"`
}

type ApproachesBody struct {
	M generated.OperationMessagingService_getApproachesResponse `xml:"getApproachesResponse"`
}

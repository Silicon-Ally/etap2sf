// Note - this is not in `conversion` to avoid a circular dependency.
package conversionsettings

import (
	"strings"

	"github.com/Silicon-Ally/etap2sf/etap"
	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
)

var AttributedUserEmail = "user-to-attribute-to@your.org"
var CallerUserEmail = "your-email@your.org"

var NovelObjectTypes = []salesforce.ObjectType{
	salesforce.ObjectType_AdditionalContext,
}

var ObjectTypeMap = map[etap.ObjectType][]salesforce.ObjectType{
	etap.ObjectType_Attachment: {
		salesforce.ObjectType_ContentDocumentLink,
		salesforce.ObjectType_ContentVersion,
		salesforce.ObjectType_AdditionalContext,
	},
	etap.ObjectType_Campaign: {
		salesforce.ObjectType_Campaign,
	},
	etap.ObjectType_Fund: {
		salesforce.ObjectType_GeneralAccountingUnit,
	},
	etap.ObjectType_Account: {
		salesforce.ObjectType_Account,
		salesforce.ObjectType_Contact,
	},
	etap.ObjectType_Relationship: {
		salesforce.ObjectType_Relationship,
		salesforce.ObjectType_Affiliation,
	},
	etap.ObjectType_Contact: {
		salesforce.ObjectType_Task,
		salesforce.ObjectType_AdditionalContext,
	},
	etap.ObjectType_Disbursement: {
		salesforce.ObjectType_Opportunity,
	},
	etap.ObjectType_Gift: {
		salesforce.ObjectType_Opportunity,
	},
	etap.ObjectType_Note: {
		salesforce.ObjectType_Task,
		salesforce.ObjectType_AdditionalContext,
	},
	etap.ObjectType_Payment: {
		salesforce.ObjectType_Payment,
	},
	etap.ObjectType_Pledge: {
		salesforce.ObjectType_Opportunity,
	},
	etap.ObjectType_RecurringGift: {
		salesforce.ObjectType_Opportunity,
	},
	etap.ObjectType_RecurringGiftSchedule: {
		salesforce.ObjectType_RecurringDonation,
	},
	etap.ObjectType_SegmentedDonation: {
		salesforce.ObjectType_Opportunity,
	},
	etap.ObjectType_SoftCredit: {
		salesforce.ObjectType_PartialSoftCredit,
		salesforce.ObjectType_AccountSoftCredit,
	},
}

var CategoryNameSubstitutions = map[string]string{
	"PersonalInformation":     "PersonalInfo",
	"OrganizationInformation": "OrgInfo",
}

var CategoryLabelSubstitutions = map[string]string{
	"Personal Information":     "Personal Info",
	"Organization Information": "Org. Info",
	"RecurringGiftSchedule":    "RecurrGiftSched",
}

var FieldNameSubstitutions = map[string]string{
	"etap_RecurringGift_NonDeductibleAmount__c":          "etap_RecurringGift_NonDeductibleAmt__c",
	"etap_RecurringGift_OriginalAccountName__c":          "etap_RecurringGift_OriginalAcctName__c",
	"etap_RecurringGift_OriginalTransactionRef__c":       "etap_RecurringGift_OriginalTransRef__c",
	"etap_RecurringGift_RecurringGiftScheduleRef__c":     "etap_RecurringGift_ScheduleRef__c",
	"etap_SegmentedDonation_LastModifiedDate__c":         "etap_SegmentedDonation_LastModDate__c",
	"etap_RecurringGiftSchedule_AccountName__c":          "etap_RecurringGiftSchedule_ActName__c",
	"etap_RecurringGiftSchedule_Attachments__c":          "etap_RecurringGiftSchedule_Atchments__c",
	"etap_RecurringGiftSchedule_CreatedDate__c":          "etap_RecurringGiftSchedule_CreateDate__c",
	"etap_RecurringGiftSchedule_LastModifiedDate__c":     "etap_RecurringGiftSchedule_LastModDte__c",
	"etap_RecurringGiftSchedule_LinkedGiftsAmount__c":    "etap_RGS_LinkedGiftsAmt__c",
	"etap_RecurringGiftSchedule_NextGiftAmount__c":       "etap_RGS_NextGiftAmt__c",
	"etap_RecurringGiftSchedule_NextGiftDate__c":         "etap_RGS_NextGiftDate__c",
	"etap_RecurringGiftSchedule_ScheduledValuable__c":    "etap_RGS_SchedValuable__c",
	"etap_SegmentedDonation_TotalNonDeductibleAmount__c": "etap_SD_TotalNonDeductibleAmount__c",
}

var FieldLabelSubstitutions = map[string]string{
	"ETap: RecurringGift: RecurringGiftScheduleRef":     "ETap: RecurringGift: ScheduleRef",
	"ETap: SegmentedDonation: LastModifiedDate":         "ETap: SegmentedDonation: LastModDate",
	"ETap: SegmentedDonation: TotalNonDeductibleAmount": "ETap: Seg. Donat.: TotalNonDeductibleAmt",
}

var SectionLabelSubstitutions = map[string]string{
	"ETapestry: Org. Info": "ETapestry: Organization Information",
}

func OverrideFieldNameValue(name string) string {
	return name
}

func OverrideFieldLabelValue(label string) string {
	replacements := map[string]string{}
	deletions := []string{}

	for _, deletion := range deletions {
		label = strings.ReplaceAll(label, deletion, "")
	}
	for k, v := range replacements {
		label = strings.ReplaceAll(label, k, v)
	}
	return label
}

func CustomFieldsForObjectType(sot salesforce.ObjectType) []*sfmetadata.CustomField {
	fields := []*sfmetadata.CustomField{}
	// This is where you can create custom fields for your migration needs.
	return fields
}

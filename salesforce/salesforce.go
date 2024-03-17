package salesforce

import (
	"fmt"
)

type ObjectType string

const (
	ObjectType_Account               ObjectType = "Account"
	ObjectType_AccountSoftCredit     ObjectType = "AccountSoftCredit"
	ObjectType_AdditionalContext     ObjectType = "EtapAdditionalContext"
	ObjectType_Affiliation           ObjectType = "Affiliation"
	ObjectType_Campaign              ObjectType = "Campaign"
	ObjectType_Contact               ObjectType = "Contact"
	ObjectType_ContentDocumentLink   ObjectType = "ContentDocumentLink"
	ObjectType_ContentVersion        ObjectType = "ContentVersion"
	ObjectType_GeneralAccountingUnit ObjectType = "GeneralAccountingUnit"
	ObjectType_GAUAllocation         ObjectType = "GAUAllocation"
	ObjectType_Opportunity           ObjectType = "Opportunity"
	ObjectType_Payment               ObjectType = "Payment"
	ObjectType_PartialSoftCredit     ObjectType = "PartialSoftCredit"
	ObjectType_RecurringDonation     ObjectType = "RecurringDonation"
	ObjectType_Relationship          ObjectType = "Relationship"
	ObjectType_Task                  ObjectType = "Task"
)

var ObjectTypes = []ObjectType{
	ObjectType_Account,
	ObjectType_AccountSoftCredit,
	ObjectType_AdditionalContext,
	ObjectType_Affiliation,
	ObjectType_Campaign,
	ObjectType_Contact,
	ObjectType_ContentDocumentLink,
	ObjectType_ContentVersion,
	ObjectType_GeneralAccountingUnit,
	ObjectType_GAUAllocation,
	ObjectType_Opportunity,
	ObjectType_Payment,
	ObjectType_PartialSoftCredit,
	ObjectType_RecurringDonation,
	ObjectType_Relationship,
	ObjectType_Task,
}

func (ot ObjectType) SalesforceName() (string, error) {
	switch ot {
	case ObjectType_Account:
		return "Account", nil
	case ObjectType_AccountSoftCredit:
		return "npsp__Account_Soft_Credit__c", nil
	case ObjectType_Affiliation:
		return "npe5__Affiliation__c", nil
	case ObjectType_AdditionalContext:
		return "etap_AdditionalContext__c", nil
	case ObjectType_Campaign:
		return "Campaign", nil
	case ObjectType_Contact:
		return "Contact", nil
	case ObjectType_ContentDocumentLink:
		return "ContentDocumentLink", nil
	case ObjectType_ContentVersion:
		return "ContentVersion", nil
	case ObjectType_GeneralAccountingUnit:
		return "npsp__General_Accounting_Unit__c", nil
	case ObjectType_GAUAllocation:
		return "npsp__Allocation__c", nil
	case ObjectType_Opportunity:
		return "Opportunity", nil
	case ObjectType_Payment:
		return "npe01__OppPayment__c", nil
	case ObjectType_PartialSoftCredit:
		return "npsp__Partial_Soft_Credit__c", nil
	case ObjectType_RecurringDonation:
		return "npe03__Recurring_Donation__c", nil
	case ObjectType_Relationship:
		return "npe4__Relationship__c", nil
	case ObjectType_Task:
		return "Task", nil
	}
	return "", fmt.Errorf("unknown object type: %s", ot)
}

func (ot ObjectType) SalesforceNameForFieldCreation() (string, error) {
	if ot == ObjectType_Task {
		return "Activity", nil
	}
	return ot.SalesforceName()
}

const MultiObjectExternalFieldKey = "etap_MultiObject_EtapRef__c"

func (ot ObjectType) SalesforceObjectExternalFieldKey() (string, error) {
	switch ot {
	case ObjectType_Account:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_AccountSoftCredit:
		return "etap_SoftCredit_Ref__c", nil
	case ObjectType_AdditionalContext:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_Affiliation:
		return "etap_Relationship_Ref__c", nil
	case ObjectType_Campaign:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_ContentDocumentLink:
		return "Id", nil
	case ObjectType_ContentVersion:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_Contact:
		return "etap_Account_Ref__c", nil
	case ObjectType_GeneralAccountingUnit:
		return "etap_Fund_Ref__c", nil
	case ObjectType_GAUAllocation:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_Opportunity:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_Payment:
		return MultiObjectExternalFieldKey, nil
	case ObjectType_PartialSoftCredit:
		return "etap_SoftCredit_Ref__c", nil
	case ObjectType_RecurringDonation:
		return "etap_RecurringGiftSchedule_Ref__c", nil
	case ObjectType_Relationship:
		return "etap_Relationship_Ref__c", nil
	case ObjectType_Task:
		return MultiObjectExternalFieldKey, nil
	}
	return "", fmt.Errorf("unknown object type for salesforce-object-external-field-key: %s", ot)
}

/*
Finding this list (below) is a bit nontrivial. I recommend using this request:

func (c *Client) ListAllLayouts(sot salesforce.ObjectType, fields []*sfmetadata.CustomField) error {
	_, err := c.metadataClient.ListMetadata([]*metaforce.ListMetadataQuery{{
		Type: "Layout",
	}})
	if err != nil {
		return fmt.Errorf("reading metadata: %w", err)
	}
	return nil
}

then look in the logged response to find each of these.
*/

func (ot ObjectType) SalesforceLayoutNeedingFilesButton() ([]string, error) {
	switch ot {
	case ObjectType_Account, ObjectType_Affiliation, ObjectType_Campaign,
		ObjectType_GeneralAccountingUnit, ObjectType_Contact, ObjectType_GAUAllocation,
		ObjectType_Opportunity, ObjectType_RecurringDonation, ObjectType_AdditionalContext,
		ObjectType_Relationship, ObjectType_ContentDocumentLink, ObjectType_ContentVersion:
		return []string{}, nil
	case ObjectType_Payment:
		return []string{"Payment"}, nil
	case ObjectType_Task:
		return []string{"Task-Task Layout"}, nil
	case ObjectType_AccountSoftCredit:
		return []string{"Account Soft Credit"}, nil
	case ObjectType_PartialSoftCredit:
		return []string{"Partial PartialSoft Credit"}, nil
	}
	return nil, fmt.Errorf("unknown object type for salesforce-layout-names: %s", ot)
}

func (ot ObjectType) SalesforceLayoutNeedingManualIntervention() ([]string, error) {
	switch ot {
	// For these types, no action is needed because we handle them as flexi-page inserts
	case ObjectType_Account, ObjectType_Affiliation, ObjectType_Campaign,
		ObjectType_GeneralAccountingUnit, ObjectType_Contact,
		ObjectType_Opportunity, ObjectType_Payment, ObjectType_RecurringDonation,
		ObjectType_Relationship, ObjectType_ContentDocumentLink, ObjectType_ContentVersion:
		return []string{}, nil
	case ObjectType_AdditionalContext:
		return []string{"etap_AdditionalContext__c-AdditionalContext Layout"}, nil
	case ObjectType_GAUAllocation:
		return []string{"npsp__Allocation__c-Allocation Layout"}, nil
	case ObjectType_AccountSoftCredit:
		return []string{"npsp__Account_Soft_Credit__c-Account Soft Credit Layout"}, nil
	case ObjectType_PartialSoftCredit:
		return []string{"npsp__Partial_Soft_Credit__c-Partial Soft Credit Layout"}, nil
	case ObjectType_Task:
		return []string{"Task-Task Layout"}, nil
	}
	return nil, fmt.Errorf("unknown object type for salesforce-layout-names: %s", ot)
}

func (ot ObjectType) SalesforceFlexiPageNames() ([]string, error) {
	switch ot {
	case ObjectType_Account:
		return []string{"NPSP_Account_Record_Page"}, nil
	case ObjectType_AccountSoftCredit:
		return []string{}, nil
	case ObjectType_AdditionalContext:
		return []string{}, nil
	case ObjectType_Affiliation:
		return []string{"NPSP_Affiliation_Record_Page"}, nil
	case ObjectType_Campaign:
		return []string{"NPSP_Campaign_Record_Page"}, nil
	case ObjectType_Contact:
		return []string{"NPSP_Contact_Record_Page"}, nil
	case ObjectType_ContentDocumentLink:
		return []string{}, nil
	case ObjectType_ContentVersion:
		return []string{}, nil
	case ObjectType_GeneralAccountingUnit:
		return []string{"NPSP_General_Accounting_Unit"}, nil
	case ObjectType_GAUAllocation:
		return []string{"NPSP_GAU_Allocation"}, nil
	case ObjectType_Opportunity:
		return []string{"npsp__NPSP_Opportunity_Record_Page", "NPSP_Opportunity_Record_Page"}, nil
	case ObjectType_Payment:
		return []string{"NPSP_Payment"}, nil
	case ObjectType_PartialSoftCredit:
		return nil, nil
	case ObjectType_RecurringDonation:
		return []string{"NPSP_Recurring_Donation"}, nil
	case ObjectType_Relationship:
		return []string{"NPSP_Relationship_Record_Page"}, nil
	case ObjectType_Task:
		return []string{}, nil
	}
	return nil, fmt.Errorf("unknown object type for salesforce-flexi-page-names: %s", ot)
}

func (sot ObjectType) HasCampaign() bool {
	switch sot {
	case ObjectType_Opportunity, ObjectType_RecurringDonation, ObjectType_Campaign:
		return true
	}
	return false
}

func (sot ObjectType) HasFund() bool {
	switch sot {
	case ObjectType_Opportunity, ObjectType_RecurringDonation:
		return true
	}
	return false
}

func (sot ObjectType) IsCustomToMigration() bool {
	switch sot {
	case ObjectType_AdditionalContext:
		return true
	}
	return false
}

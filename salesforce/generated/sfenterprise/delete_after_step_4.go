package sfenterprise

// This file is a placeholder to allow for the packages to build.
// It should be deleted as soon as step_04 is completed.

type ID string

type Account struct {
	Etap_MultiObject_EtapRef__c *string
}

type Npsp__Account_Soft_Credit__c struct {
	Etap_SoftCredit_Ref__c *string
}

type Etap_AdditionalContext__c struct {
	Name *string
}

type Npe5__Affiliation__c struct {
	Etap_Relationship_Ref__c *string
}

type Campaign struct {
	Etap_MultiObject_EtapRef__c *string
}

type Contact struct {
	Etap_Account_Ref__c *string
}

type ContentDocumentLink struct {
	ContentDocumentId *ID
	LinkedEntityId    *ID
}

type ContentVersion struct {
	Etap_MultiObject_EtapRef__c *string
}

type Npsp__General_Accounting_Unit__c struct {
	Etap_Fund_Ref__c *string
}

type Opportunity struct {
	Etap_MultiObject_EtapRef__c *string
}

type Npe01__OppPayment__c struct {
	Etap_Payment_Ref__c *string
}

type Npe03__Recurring_Donation__c struct {
	Etap_RecurringGiftSchedule_Ref__c *string
}

type Npsp__Partial_Soft_Credit__c struct {
	Etap_SoftCredit_Ref__c *string
}

type Npe4__Relationship__c struct {
	Etap_Relationship_Ref__c *string
}

type Task struct {
	Etap_MultiObject_EtapRef__c *string
}

type Npsp__Allocation__c struct {
	Etap_MultiObject_EtapRef__c *string
}

type Npo02__Household__c struct{}

type ContentNote struct{}

type Task_Subject_ string

type RecordType struct {
	Id          *ID
	Name        *string
	SobjectType *string
}

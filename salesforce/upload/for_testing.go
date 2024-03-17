package upload

import (
	"strconv"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
)

func GetUploaderForLocalValidation() (*Uploader, error) {
	u := &Uploader{}
	u.NumThreads = 1
	u.MaxErrors = 1
	u.Failed = map[string]bool{}
	u.Succeeded = map[string]bool{}
	u.IDMap = map[string]string{}
	u.Wetrun = false
	u.DoShuffles = false
	u.client = &fakeClient{}
	return u, nil
}

type fakeClient struct {
	id int
}

func (u *fakeClient) nextID(prefix string) (string, error) {
	u.id++
	return prefix + "-" + strconv.Itoa(u.id), nil
}
func (u *fakeClient) UpsertCampaign(*sfenterprise.Campaign) (string, error) {
	return u.nextID("campaign")
}
func (u *fakeClient) UpsertRelationship(*sfenterprise.Npe4__Relationship__c) (string, error) {
	return u.nextID("relationship")
}
func (u *fakeClient) UpsertAffiliation(*sfenterprise.Npe5__Affiliation__c) (string, error) {
	return u.nextID("affiliation")
}
func (u *fakeClient) UpsertContact(*sfenterprise.Contact) (string, error) {
	return u.nextID("contact")
}
func (u *fakeClient) UpsertAccount(*sfenterprise.Account) (string, error) {
	return u.nextID("account")
}
func (u *fakeClient) UpsertGeneralAccountingUnit(*sfenterprise.Npsp__General_Accounting_Unit__c) (string, error) {
	return u.nextID("gau")
}
func (u *fakeClient) UpsertGAUAllocation(*sfenterprise.Npsp__Allocation__c) (string, error) {
	return u.nextID("gauallocation")
}
func (u *fakeClient) UpsertOpportunity(*sfenterprise.Opportunity) (string, error) {
	return u.nextID("opportunity")
}
func (u *fakeClient) UpsertPayment(*sfenterprise.Npe01__OppPayment__c) (string, error) {
	return u.nextID("payment")
}
func (u *fakeClient) UpsertRecurringDonation(*sfenterprise.Npe03__Recurring_Donation__c) (string, error) {
	return u.nextID("recurringdonation")
}
func (u *fakeClient) UpsertTask(*sfenterprise.Task) (string, error) {
	return u.nextID("task")
}
func (u *fakeClient) UpsertPartialSoftCredit(*sfenterprise.Npsp__Partial_Soft_Credit__c) (string, error) {
	return u.nextID("partialsoftcredit")
}
func (u *fakeClient) UpsertAccountSoftCredit(*sfenterprise.Npsp__Account_Soft_Credit__c) (string, error) {
	return u.nextID("accountsoftcredit")
}
func (u *fakeClient) UpsertAdditionalContext(*sfenterprise.Etap_AdditionalContext__c) (string, error) {
	return u.nextID("additionalcontext")
}
func (u *fakeClient) UpsertContentVersion(*sfenterprise.ContentVersion) (string, error) {
	return u.nextID("contentversion")
}
func (u *fakeClient) UpsertContentDocumentLink(*sfenterprise.ContentDocumentLink) (string, error) {
	return u.nextID("contentdocumentlink")
}
func (u *fakeClient) LookupContentDocumentByVersion(id sfenterprise.ID) (sfenterprise.ID, error) {
	return id + "-contentdocument", nil
}

package upload

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Silicon-Ally/etap2sf/conv/conversion"
	esfutils "github.com/Silicon-Ally/etap2sf/salesforce/clients/enterprise/utils"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
	"github.com/Silicon-Ally/etap2sf/utils"
)

type client interface {
	UpsertCampaign(*sfenterprise.Campaign) (string, error)
	UpsertRelationship(*sfenterprise.Npe4__Relationship__c) (string, error)
	UpsertRecurringDonation(*sfenterprise.Npe03__Recurring_Donation__c) (string, error)
	UpsertAffiliation(*sfenterprise.Npe5__Affiliation__c) (string, error)
	UpsertPayment(*sfenterprise.Npe01__OppPayment__c) (string, error)
	UpsertContact(*sfenterprise.Contact) (string, error)
	UpsertAccount(*sfenterprise.Account) (string, error)
	UpsertGeneralAccountingUnit(*sfenterprise.Npsp__General_Accounting_Unit__c) (string, error)
	UpsertGAUAllocation(*sfenterprise.Npsp__Allocation__c) (string, error)
	UpsertTask(*sfenterprise.Task) (string, error)
	UpsertPartialSoftCredit(*sfenterprise.Npsp__Partial_Soft_Credit__c) (string, error)
	UpsertAccountSoftCredit(*sfenterprise.Npsp__Account_Soft_Credit__c) (string, error)
	UpsertOpportunity(*sfenterprise.Opportunity) (string, error)
	UpsertAdditionalContext(*sfenterprise.Etap_AdditionalContext__c) (string, error)
	UpsertContentVersion(*sfenterprise.ContentVersion) (string, error)
	UpsertContentDocumentLink(*sfenterprise.ContentDocumentLink) (string, error)
	LookupContentDocumentByVersion(sfenterprise.ID) (sfenterprise.ID, error)
}

type Uploader struct {
	NumThreads       int
	MaxErrors        int
	client           client
	Failed           map[string]bool
	Succeeded        map[string]bool
	IDMap            map[string]string
	Wetrun           bool
	DoShuffles       bool
	Todo             int
	name             string
	stateChangeMutex sync.Mutex
	anyThreadDead    bool
	cleanups         []func() error
}

var uploaderMemoPath = filepath.Join(utils.ProjectRoot(), "data", "uploader.json")

func GetOrCreateUploader(doShuffles bool) (*Uploader, error) {
	u := &Uploader{}
	data, err := os.ReadFile(uploaderMemoPath)
	if os.IsNotExist(err) {
		u.NumThreads = 1
		u.MaxErrors = 1
		u.Failed = map[string]bool{}
		u.Succeeded = map[string]bool{}
		u.IDMap = map[string]string{}
		u.Wetrun = true
		u.DoShuffles = doShuffles
	} else if err != nil {
		return nil, fmt.Errorf("reading uploader from %s: %w", uploaderMemoPath, err)
	} else {
		if err := json.Unmarshal(data, u); err != nil {
			return nil, fmt.Errorf("unmarshalling uploader: %w", err)
		}
	}
	client, err := esfutils.NewSandboxClient()
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}
	undoFn, err := client.DisableNPSPRelationshipTriggers()
	if err != nil {
		return nil, fmt.Errorf("disabling npsp triggers: %w", err)
	}
	u.cleanups = append(u.cleanups, undoFn)
	u.client = client
	return u, nil
}

func save(u *Uploader) error {
	if !u.Wetrun {
		return nil
	}
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling uploader: %w", err)
	}
	if err := os.WriteFile(uploaderMemoPath, data, 0777); err != nil {
		return fmt.Errorf("writing uploader to %s: %w", uploaderMemoPath, err)
	}
	return nil
}

func cleanup(u *Uploader) {
	for _, fn := range u.cleanups {
		if err := fn(); err != nil {
			fmt.Printf("error cleaning up: %v", err)
		}
	}
}

func (u *Uploader) Upload(output *conversion.Output) error {
	defer cleanup(u)
	defer func() {
		if err := save(u); err != nil {
			log.Printf("failed to save output: %v", err)
		}
	}()
	if err := u.uploadCampaigns(output); err != nil {
		return fmt.Errorf("uploading campaigns: %w", err)
	}
	if err := u.uploadGAUs(output); err != nil {
		return fmt.Errorf("uploading gaus: %w", err)
	}
	if err := u.uploadAccounts(output); err != nil {
		return fmt.Errorf("uploading accounts: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInContacts(u.IDMap)); err != nil {
		return fmt.Errorf("replacing contact ids: %w", err)
	}
	if err := u.uploadContacts(output); err != nil {
		return fmt.Errorf("uploading contacts: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInRelationships(u.IDMap)); err != nil {
		return fmt.Errorf("replacing relationship ids: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInAffiliations(u.IDMap)); err != nil {
		return fmt.Errorf("replacing affiliation ids: %w", err)
	}
	if err := u.uploadRelationships(output); err != nil {
		return fmt.Errorf("uploading relationships: %w", err)
	}
	if err := u.uploadAffiliations(output); err != nil {
		return fmt.Errorf("uploading affiliations: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInRecurringDonations(u.IDMap)); err != nil {
		return fmt.Errorf("replacing recurring donation ids: %w", err)
	}
	if err := u.uploadRecurringDonations(output); err != nil {
		return fmt.Errorf("uploading recurring donations: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInOpportunities(u.IDMap)); err != nil {
		return fmt.Errorf("replacing opportunity ids: %w", err)
	}
	if err := u.uploadOpportunities(output); err != nil {
		return fmt.Errorf("uploading opportunities: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInPayments(u.IDMap)); err != nil {
		return fmt.Errorf("replacing payment ids: %w", err)
	}
	if err := u.uploadPayments(output); err != nil {
		return fmt.Errorf("uploading payments: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInGAUAllocations(u.IDMap)); err != nil {
		return fmt.Errorf("replacing gau allocation ids: %w", err)
	}
	if err := u.uploadGAUAllocations(output); err != nil {
		return fmt.Errorf("uploading gau allocations: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInPartialSoftCredits(u.IDMap)); err != nil {
		return fmt.Errorf("replacing partial soft credit ids: %w", err)
	}
	if err := u.uploadPartialSoftCredits(output); err != nil {
		return fmt.Errorf("uploading partial soft credits: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInAccountSoftCredits(u.IDMap)); err != nil {
		return fmt.Errorf("replacing account soft credit ids: %w", err)
	}
	if err := u.uploadAccountSoftCredits(output); err != nil {
		return fmt.Errorf("uploading account soft credits: %w", err)
	}
	if err := u.uploadAdditionalContexts(output); err != nil {
		return fmt.Errorf("uploading additional contexts: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInTasks(u.IDMap)); err != nil {
		return fmt.Errorf("replacing task ids: %w", err)
	}
	if err := u.uploadTasks(output); err != nil {
		return fmt.Errorf("uploading tasks: %w", err)
	}
	if err := u.uploadContentVersions(output); err != nil {
		return fmt.Errorf("uploading content versions: %w", err)
	}
	if err := handleErrors(output.ReplaceAllIDsInContentDocumentLinks(u.IDMap)); err != nil {
		return fmt.Errorf("replacing content document link ids: %w", err)
	}
	if err := u.lookupContentDocumentIDsFromContentVersionIDs(output); err != nil {
		return fmt.Errorf("looking up content document ids from content version ids: %w", err)
	}
	if err := u.uploadContentDocumentLinks(output); err != nil {
		return fmt.Errorf("uploading content document links: %w", err)
	}
	fmt.Printf("Your upload has completed successfully. Congratulations!")
	return nil
}

func runErrorsOnly[T any](u *Uploader, ts []T, idFn func(t T) string, fn func(t T) (string, error)) error {
	u.anyThreadDead = false
	retained := []T{}
	for _, t := range ts {
		tt := t
		id := idFn(t)
		if failed, ok := u.Failed[id]; ok && failed {
			retained = append(retained, tt)
		}
	}
	if len(retained) == 0 {
		return fmt.Errorf("no errors to run errors only on")
	}
	for i, t := range retained {
		tt := t
		fmt.Printf("RUNNING ERRORS ONLY %d/%d\n", i, len(retained))
		if errs := runThread(u, []T{tt}, idFn, fn, false); len(errs) > 0 {
			return fmt.Errorf("running errors only: %w", handleErrors(errs))
		}
	}
	return nil
}

func run[T any](u *Uploader, ts []T, idFn func(t T) string, fn func(t T) (string, error), hard bool) error {
	u.anyThreadDead = false
	defer func() {
		if r := recover(); r != nil {
			if err := save(u); err != nil {
				log.Printf("failed to save output: %v", err)
			}
			panic(r)
		} else {
			if err := save(u); err != nil {
				log.Printf("failed to save output: %v", err)
			}
		}
	}()
	if len(ts) == 0 {
		return nil
	}
	if u.DoShuffles {
		// We shuffle here because otherwise we're likely to run into lock contention in salesforce
		// because opportunities are looked up per-user.
		utils.Shuffle(ts)
	}
	u.name = fmt.Sprintf("%T", ts[0])
	retained := []T{}
	for _, t := range ts {
		id := idFn(t)
		if !u.Succeeded[id] || hard {
			retained = append(retained, t)
		}
	}
	sort.Slice(retained, func(i, j int) bool {
		idI := idFn(retained[i])
		idJ := idFn(retained[j])
		if u.Failed[idI] != u.Failed[idJ] {
			return u.Failed[idI] && !u.Failed[idJ]
		}
		// We don't have a stable sort for a reason - our entity keys are highly correlated
		// so inserting objects (like relationships) can cause a lot of lock contention if
		// doing so multi-threaded in a stable way.
		return false
	})
	u.Todo = len(retained)

	split := splitByNumThreads(retained, u.NumThreads)
	errorsChan := make(chan []error)
	for _, ts := range split {
		go func(ts []T) {
			errorsChan <- runThread(u, ts, idFn, fn, hard)
		}(ts)
	}
	errors := []error{}
	for range split {
		errors = append(errors, <-errorsChan...)
	}
	if len(errors) > 0 {
		if err := runErrorsOnly(u, retained, idFn, fn); err != nil {
			return fmt.Errorf("running errors only: %w", err)
		}
		return run(u, retained, idFn, fn, hard)
	}
	fmt.Printf("Done with %s (%d retained, resulted in %d errors)\n", u.name, len(retained), len(errors))
	return handleErrors(errors)
}

func handleErrors(errors []error) error {
	if len(errors) > 0 {
		filePath, _ := utils.WriteErrorsToTempFile(errors, "upload-errors")
		return fmt.Errorf("step yielded %d errors:\n\tall errors: %s\n\tfirst error: %w", len(errors), filePath, errors[0])
	}
	return nil
}

func runThread[T any](u *Uploader, ts []T, idFn func(t T) string, fn func(t T) (string, error), hard bool) []error {
	maxErrorsPerThread := u.MaxErrors / u.NumThreads
	if maxErrorsPerThread < 1 {
		return []error{fmt.Errorf("max errors per thread must be at least 1")}
	}
	errors := []error{}
	for _, t := range ts {
		if len(errors) >= maxErrorsPerThread {
			u.anyThreadDead = true
			return errors
		}
		if u.anyThreadDead {
			return errors
		}
		id := idFn(t)
		u.stateChangeMutex.Lock()
		dupe, ok := u.IDMap[id]
		u.stateChangeMutex.Unlock()
		if !hard {
			if ok {
				errors = append(errors, fmt.Errorf("duplicate ref %s (already uploaded as %s)", id, dupe))
				u.failed(idFn(t))
			}
		}
		resultID, err := fn(t)
		if err != nil {
			errors = append(errors, err)
			u.failed(idFn(t))
		} else {
			u.succeeded(idFn(t), resultID)
		}
	}
	return errors
}

func (u *Uploader) succeeded(id, resultID string) {
	u.stateChangeMutex.Lock()
	defer u.stateChangeMutex.Unlock()
	u.Succeeded[id] = true
	delete(u.Failed, id)
	u.Todo--
	u.IDMap[id] = resultID
	u.report()
	if len(u.Succeeded)%25 != 0 {
		return
	}
	if err := save(u); err != nil {
		log.Printf("failed to save uploader: %v", err)
	}
}

func (u *Uploader) failed(id string) {
	u.stateChangeMutex.Lock()
	defer u.stateChangeMutex.Unlock()
	u.Failed[id] = true
	delete(u.Succeeded, id)
	u.Todo--
	u.report()
}

func (u *Uploader) report() {
	if u.Wetrun {
		fmt.Printf("%s: Success %d, Failed %d, Remaining %d\n", u.name, len(u.Succeeded), len(u.Failed), u.Todo)
	}
}

func (u *Uploader) lookupContentDocumentIDsFromContentVersionIDs(output *conversion.Output) error {
	return run(
		u,
		output.ContentDocumentLinks,
		func(c *sfenterprise.ContentDocumentLink) string { return string(*c.ContentDocumentId) + "-CDL" },
		func(cdl *sfenterprise.ContentDocumentLink) (string, error) {
			contentVersionId := *cdl.ContentDocumentId
			contentDocumentId, err := u.client.LookupContentDocumentByVersion(contentVersionId)
			if err != nil {
				return "", err
			}
			cdl.ContentDocumentId = &contentDocumentId
			return string(*cdl.ContentDocumentId), nil
		},
		true, /* This always needs to be done from scratch */
	)
}

func splitByNumThreads[T any](ts []T, n int) [][]T {
	result := [][]T{}
	for i := 0; i < n; i++ {
		result = append(result, []T{})
	}
	for i, t := range ts {
		result[i%n] = append(result[i%n], t)
	}
	return result
}

func (u *Uploader) uploadCampaigns(output *conversion.Output) error {
	return run(
		u,
		output.Campaigns,
		func(c *sfenterprise.Campaign) string { return *c.Etap_MultiObject_EtapRef__c },
		u.client.UpsertCampaign,
		false)
}
func (u *Uploader) uploadGAUs(output *conversion.Output) error {
	return run(
		u,
		output.GeneralAccountingUnits,
		func(c *sfenterprise.Npsp__General_Accounting_Unit__c) string { return *c.Etap_Fund_Ref__c },
		u.client.UpsertGeneralAccountingUnit,
		false)
}
func (u *Uploader) uploadAccounts(output *conversion.Output) error {
	return run(
		u,
		output.Accounts,
		func(a *sfenterprise.Account) string { return *a.Etap_MultiObject_EtapRef__c },
		u.client.UpsertAccount,
		false)
}
func (u *Uploader) uploadContacts(output *conversion.Output) error {
	return run(
		u,
		output.Contacts,
		func(a *sfenterprise.Contact) string { return *a.Etap_Account_Ref__c },
		u.client.UpsertContact,
		false)
}
func (u *Uploader) uploadRelationships(output *conversion.Output) error {
	return run(
		u,
		output.Relationships,
		// No longer need differentiation here - we add a 1-of-2 2-of-2 suffix to differentiate
		func(a *sfenterprise.Npe4__Relationship__c) string { return *a.Etap_Relationship_Ref__c },
		u.client.UpsertRelationship,
		false)
}
func (u *Uploader) uploadOpportunities(output *conversion.Output) error {
	return run(
		u,
		output.Opportunities,
		func(a *sfenterprise.Opportunity) string { return *a.Etap_MultiObject_EtapRef__c },
		u.client.UpsertOpportunity,
		false)
}
func (u *Uploader) uploadAffiliations(output *conversion.Output) error {
	return run(
		u,
		output.Affiliations,
		func(a *sfenterprise.Npe5__Affiliation__c) string { return *a.Etap_Relationship_Ref__c },
		u.client.UpsertAffiliation,
		false)
}
func (u *Uploader) uploadRecurringDonations(output *conversion.Output) error {
	return run(
		u,
		output.RecurringDonations,
		func(a *sfenterprise.Npe03__Recurring_Donation__c) string { return *a.Etap_RecurringGiftSchedule_Ref__c },
		u.client.UpsertRecurringDonation,
		false)
}
func (u *Uploader) uploadPayments(output *conversion.Output) error {
	return run(
		u,
		output.Payments,
		func(a *sfenterprise.Npe01__OppPayment__c) string { return *a.Etap_Payment_Ref__c },
		u.client.UpsertPayment,
		false)
}
func (u *Uploader) uploadGAUAllocations(output *conversion.Output) error {
	return run(
		u,
		output.GAUAllocations,
		func(a *sfenterprise.Npsp__Allocation__c) string { return *a.Etap_MultiObject_EtapRef__c },
		u.client.UpsertGAUAllocation,
		false)
}
func (u *Uploader) uploadTasks(output *conversion.Output) error {
	return run(
		u,
		output.Tasks,
		func(t *sfenterprise.Task) string { return *t.Etap_MultiObject_EtapRef__c },
		u.client.UpsertTask,
		false)
}
func (u *Uploader) uploadPartialSoftCredits(output *conversion.Output) error {
	return run(
		u,
		output.PartialSoftCredits,
		func(t *sfenterprise.Npsp__Partial_Soft_Credit__c) string { return *t.Etap_SoftCredit_Ref__c },
		u.client.UpsertPartialSoftCredit,
		false)
}
func (u *Uploader) uploadAccountSoftCredits(output *conversion.Output) error {
	return run(
		u,
		output.AccountSoftCredits,
		func(t *sfenterprise.Npsp__Account_Soft_Credit__c) string { return *t.Etap_SoftCredit_Ref__c },
		u.client.UpsertAccountSoftCredit,
		false)
}
func (u *Uploader) uploadAdditionalContexts(output *conversion.Output) error {
	return run(
		u,
		output.AdditionalContexts,
		func(t *sfenterprise.Etap_AdditionalContext__c) string { return *t.Name },
		u.client.UpsertAdditionalContext,
		false)
}
func (u *Uploader) uploadContentVersions(output *conversion.Output) error {
	return run(
		u,
		output.ContentVersions,
		func(t *sfenterprise.ContentVersion) string { return *t.Etap_MultiObject_EtapRef__c },
		u.client.UpsertContentVersion,
		false)
}
func (u *Uploader) uploadContentDocumentLinks(output *conversion.Output) error {
	return run(
		u,
		output.ContentDocumentLinks,
		func(t *sfenterprise.ContentDocumentLink) string {
			return (string(*t.ContentDocumentId) + string(*t.LinkedEntityId))
		},
		u.client.UpsertContentDocumentLink,
		false)
}

package conversion

import (
	"errors"
)

var err = errors.New("delete the file 'delete_me_after_step_12.go', and remove build tags from the remainder of the `conversio` package to continue")
var errs = []error{err}

func (o *Output) ReplaceAllIDsInContacts(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInRelationships(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInAffiliations(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInOpportunities(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInPayments(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInRecurringDonations(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInGAUAllocations(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInPartialSoftCredits(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInAccountSoftCredits(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInTasks(idMap map[string]string) []error { return errs }

func (o *Output) ReplaceAllIDsInContentDocumentLinks(idMap map[string]string) []error { return errs }

func (i *Input) Convert() (*Output, error) { return nil, err }

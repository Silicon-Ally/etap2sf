//go:build ignore_until_step_12

package conversion

import (
	"fmt"
	"strings"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
)

const prefix = "NeedIdHaveRef:"
const noreplPrefix = "NoRepl:"

func noReplacementForId(i *sfenterprise.ID) (*sfenterprise.ID, error) {
	if i == nil || *i == "" {
		return nil, fmt.Errorf("ref is nil or empty")
	}
	id := sfenterprise.ID(fmt.Sprintf("%s%s", noreplPrefix, *i))
	return &id, nil
}

func idPlaceholderForRef(ref *string) (*sfenterprise.ID, error) {
	if ref == nil || *ref == "" {
		return nil, fmt.Errorf("ref is nil or empty")
	}
	id := sfenterprise.ID(fmt.Sprintf("%s%s", prefix, *ref))
	return &id, nil
}

func needsRepl(ref string) bool {
	return strings.HasPrefix(ref, prefix) || strings.HasPrefix(ref, noreplPrefix)
}

func getRepl(og *sfenterprise.ID, idMap map[string]string) (*sfenterprise.ID, error) {
	if og == nil {
		return nil, fmt.Errorf("original is nil")
	}
	original := string(*og)
	if !needsRepl(original) {
		return nil, fmt.Errorf("expected to need replacement, but was %q", original)
	}
	if strings.HasPrefix(original, noreplPrefix) {
		original = strings.TrimPrefix(original, noreplPrefix)
		if !strings.HasPrefix(original, prefix) {
			return ptr(sfenterprise.ID(original)), nil
		}
	}
	ref := strings.TrimPrefix(original, prefix)
	newId, ok := idMap[ref]
	if !ok {
		return nil, fmt.Errorf("couldn't find id for ref %q", ref)
	}
	if newId == "" {
		return nil, fmt.Errorf("new id is empty")
	}
	sfid := sfenterprise.ID(newId)
	return &sfid, nil
}

func (o *Output) ReplaceAllIDsInContacts(idMap map[string]string) []error {
	return replaceAllIDs(o.Contacts, idMap, func(c *sfenterprise.Contact) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			c.AccountId,
		}
	})
}

func (o *Output) ReplaceAllIDsInRelationships(idMap map[string]string) []error {
	return replaceAllIDs(o.Relationships, idMap, func(r *sfenterprise.Npe4__Relationship__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			r.Npe4__Contact__c,
			r.Npe4__RelatedContact__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInAffiliations(idMap map[string]string) []error {
	return replaceAllIDs(o.Affiliations, idMap, func(a *sfenterprise.Npe5__Affiliation__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			a.Npe5__Contact__c,
			a.Npe5__Organization__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInOpportunities(idMap map[string]string) []error {
	return replaceAllIDs(o.Opportunities, idMap, func(o *sfenterprise.Opportunity) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			o.AccountId,
			o.ContactId,
			o.CampaignId,
			o.Npe03__Recurring_Donation__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInPayments(idMap map[string]string) []error {
	return replaceAllIDs(o.Payments, idMap, func(p *sfenterprise.Npe01__OppPayment__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			p.Npe01__Opportunity__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInRecurringDonations(idMap map[string]string) []error {
	return replaceAllIDs(o.RecurringDonations, idMap, func(rd *sfenterprise.Npe03__Recurring_Donation__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			rd.Npe03__Organization__c,
			rd.Npe03__Contact__c,
			rd.Npe03__Recurring_Donation_Campaign__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInGAUAllocations(idMap map[string]string) []error {
	return replaceAllIDs(o.GAUAllocations, idMap, func(a *sfenterprise.Npsp__Allocation__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			a.Npsp__Campaign__c,
			a.Npsp__Opportunity__c,
			a.Npsp__Recurring_Donation__c,
			a.Npsp__General_Accounting_Unit__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInPartialSoftCredits(idMap map[string]string) []error {
	return replaceAllIDs(o.PartialSoftCredits, idMap, func(a *sfenterprise.Npsp__Partial_Soft_Credit__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			a.Npsp__Opportunity__c,
			a.Npsp__Contact__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInAccountSoftCredits(idMap map[string]string) []error {
	return replaceAllIDs(o.AccountSoftCredits, idMap, func(a *sfenterprise.Npsp__Account_Soft_Credit__c) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			a.Npsp__Opportunity__c,
			a.Npsp__Account__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInTasks(idMap map[string]string) []error {
	return replaceAllIDs(o.Tasks, idMap, func(a *sfenterprise.Task) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			a.WhoId,
			a.WhatId,
			a.AccountId,
			a.Etap_AdditionalContextForRecord__c,
		}
	})
}

func (o *Output) ReplaceAllIDsInContentDocumentLinks(idMap map[string]string) []error {
	return replaceAllIDs(o.ContentDocumentLinks, idMap, func(a *sfenterprise.ContentDocumentLink) []*sfenterprise.ID {
		return []*sfenterprise.ID{
			a.ContentDocumentId,
			a.LinkedEntityId,
		}
	})
}

func replaceAllIDs[T any](ts []*T, idMap map[string]string, fieldsFn func(*T) []*sfenterprise.ID) []error {
	errors := []error{}
	for _, t := range ts {
		fields := fieldsFn(t)
		for i := range fieldsFn(t) {
			if fields[i] == nil {
				continue
			}
			if repl, err := getRepl(fields[i], idMap); err != nil {
				errors = append(errors, fmt.Errorf("getting replacement for %T's %s: %w (potential solution: is it possible that the pointer used for the placeholder is shared between populatable-entities)", t, *fields[i], err))
			} else {
				*fields[i] = *repl
			}
		}
	}
	return errors
}

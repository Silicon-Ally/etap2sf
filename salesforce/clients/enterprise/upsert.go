package client

import (
	"fmt"
	"strings"

	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
	"github.com/Silicon-Ally/etap2sf/utils"
	"github.com/tzmfreedom/go-soapforce"
)

func (c *Client) UpsertRelationship(relationship *sfenterprise.Npe4__Relationship__c) (string, error) {
	return c.upsert(salesforce.ObjectType_Relationship, relationship)
}

func (c *Client) UpsertAffiliation(affiliation *sfenterprise.Npe5__Affiliation__c) (string, error) {
	return c.upsert(salesforce.ObjectType_Affiliation, affiliation)
}

func (c *Client) UpsertContact(contact *sfenterprise.Contact) (string, error) {
	return c.upsert(salesforce.ObjectType_Contact, contact)
}

func (c *Client) UpsertAccount(account *sfenterprise.Account) (string, error) {
	return c.upsert(salesforce.ObjectType_Account, account)
}

func (c *Client) UpsertCampaign(campaign *sfenterprise.Campaign) (string, error) {
	return c.upsert(salesforce.ObjectType_Campaign, campaign)
}

func (c *Client) UpsertGeneralAccountingUnit(gau *sfenterprise.Npsp__General_Accounting_Unit__c) (string, error) {
	return c.upsert(salesforce.ObjectType_GeneralAccountingUnit, gau)
}

func (c *Client) UpsertGAUAllocation(gaua *sfenterprise.Npsp__Allocation__c) (string, error) {
	return c.upsert(salesforce.ObjectType_GAUAllocation, gaua)
}

func (c *Client) UpsertOpportunity(o *sfenterprise.Opportunity) (string, error) {
	return c.upsert(salesforce.ObjectType_Opportunity, o)
}

func (c *Client) UpsertPayment(p *sfenterprise.Npe01__OppPayment__c) (string, error) {
	return c.upsert(salesforce.ObjectType_Payment, p)
}

func (c *Client) UpsertRecurringDonation(rd *sfenterprise.Npe03__Recurring_Donation__c) (string, error) {
	return c.upsert(salesforce.ObjectType_RecurringDonation, rd)
}

func (c *Client) UpsertTask(task *sfenterprise.Task) (string, error) {
	return c.upsert(salesforce.ObjectType_Task, task)
}

func (c *Client) UpsertPartialSoftCredit(psc *sfenterprise.Npsp__Partial_Soft_Credit__c) (string, error) {
	return c.upsert(salesforce.ObjectType_PartialSoftCredit, psc)
}

func (c *Client) UpsertAccountSoftCredit(psc *sfenterprise.Npsp__Account_Soft_Credit__c) (string, error) {
	return c.upsert(salesforce.ObjectType_AccountSoftCredit, psc)
}

func (c *Client) UpsertAdditionalContext(ac *sfenterprise.Etap_AdditionalContext__c) (string, error) {
	return c.upsert(salesforce.ObjectType_AdditionalContext, ac)
}

func (c *Client) UpsertContentVersion(cv *sfenterprise.ContentVersion) (string, error) {
	return c.upsert(salesforce.ObjectType_ContentVersion, cv)
}

func (c *Client) UpsertContentDocumentLink(value *sfenterprise.ContentDocumentLink) (string, error) {
	ots, err := salesforce.ObjectType_ContentDocumentLink.SalesforceName()
	if err != nil {
		return "", fmt.Errorf("getting salesforce name: %w", err)
	}
	fields, err := StructToFieldsMap(value)
	if err != nil {
		return "", fmt.Errorf("converting struct to map: %w", err)
	}
	sobj := &soapforce.SObject{
		Type:   ots,
		Fields: fields,
	}
	results, err := c.gc.EnterpriseClient.Create([]*soapforce.SObject{sobj})
	if err != nil {
		return "", fmt.Errorf("generally creating cdl: %w", err)
	}
	if len(results) != 1 {
		return "", fmt.Errorf("expected 1 result, got %d", len(results))
	}
	result := results[0]
	if len(result.Errors) > 0 {
		return "", fmt.Errorf("errors found in response - see %+v", result.Errors)
	}
	if !result.Success {
		return "", fmt.Errorf("failed to upsert field - see response")
	}
	if result.Id == "" {
		return "", fmt.Errorf("id is empty")
	}
	return result.Id, nil
}

func (c *Client) upsert(sot salesforce.ObjectType, value any) (string, error) {
	return c.upsertWithOptionalRetry(sot, value, true, false)
}

func (c *Client) upsertWithOptionalRetry(sot salesforce.ObjectType, value any, optionalRetry bool, deleteCreates bool) (string, error) {
	ots, err := sot.SalesforceName()
	if err != nil {
		return "", fmt.Errorf("getting salesforce name: %w", err)
	}
	fieldKey, err := sot.SalesforceObjectExternalFieldKey()
	if err != nil {
		return "", fmt.Errorf("getting salesforce object external field key: %w", err)
	}
	fields, err := StructToFieldsMap(value)
	if err != nil {
		return "", fmt.Errorf("converting struct to map: %w", err)
	}
	if deleteCreates {
		delete(fields, "CreatedDate")
		delete(fields, "CreatedById")
		delete(fields, "LastModifiedDate")
		delete(fields, "LastModifiedById")
		delete(fields, "CompletedDateTime")
		delete(fields, "npe4__Contact__c")
		delete(fields, "npe4__RelatedContact__c")
		delete(fields, "npe5__Organization__c")
		delete(fields, "npe5__Contact__c")
		delete(fields, "ContactId")
	}
	key, ok := fields[fieldKey]
	if !ok {
		return "", fmt.Errorf("field key %s not found in fields", fieldKey)
	}
	keyAsStr, ok := key.(string)
	if !ok {
		return "", fmt.Errorf("field key %s is not a string, it is a %T", fieldKey, key)
	}
	if keyAsStr == "" {
		return "", fmt.Errorf("field key %s is empty", fieldKey)
	}
	sobj := &soapforce.SObject{
		Type:   ots,
		Fields: fields,
	}
	results, err := c.gc.EnterpriseClient.Upsert([]*soapforce.SObject{sobj}, fieldKey)
	if err != nil {
		return "", fmt.Errorf("generally upserting %s: %w", sot, err)
	}
	if len(results) != 1 {
		return "", fmt.Errorf("expected 1 result, got %d", len(results))
	}
	tmpFile, err := utils.WriteValueToTempJSONFile(results, fmt.Sprintf("upsert-%s-error", sot))
	if err != nil {
		return "", fmt.Errorf("writing response message: %w", err)
	}
	result := results[0]
	if len(result.Errors) > 0 {
		err0 := result.Errors[0]
		if *err0.StatusCode == "DUPLICATE_EXTERNAL_ID" {
			m := err0.Message
			m = m[strings.Index(m, "[")+1 : strings.Index(m, "]")]
			splits := strings.Split(m, ",")
			for i, split := range splits {
				splits[i] = strings.TrimSpace(split)
			}
			resp, err := c.gc.EnterpriseClient.Delete(splits)
			if err != nil {
				return "", fmt.Errorf("trying to delete duplicates: %w", err)
			}
			for _, r := range resp {
				if !r.Success && len(r.Errors) > 0 && *r.Errors[0].StatusCode != "ENTITY_IS_DELETED" {
					return "", fmt.Errorf("error deleting duplicate: %s", r.Errors[0].Message)
				}
			}
			return c.upsertWithOptionalRetry(sot, value, false, false)
		}
		if *err0.StatusCode == "INVALID_FIELD_FOR_INSERT_UPDATE" {
			if strings.Contains(err0.Message, "PathOnClient") {
				return result.Id, nil
			}
			/*
					if strings.Contains(err0.Message, "Unable to create/update fields") && (strings.Contains(err0.Message, "CreatedDate") || strings.Contains(err0.Message, "LastModifiedDate")) {

					results, err := c.gc.EnterpriseClient.Create([]*soapforce.SObject{sobj})
					if err != nil {
						return "", fmt.Errorf("generally creating %s: %w", sot, err)
					}
					if results[0].Success && results[0].Id != "" {
						return results[0].Id, nil
					}
				return c.upsertWithOptionalRetry(sot, value, false, true)

			*/
			if optionalRetry {
				return c.upsertWithOptionalRetry(sot, value, false, true)
			}
		}
		return "", fmt.Errorf("errors found in response - see %s", tmpFile)
	}
	if !result.Success {
		return "", fmt.Errorf("failed to upsert field - see %s", tmpFile)
	}
	if result.Id == "" {
		return "", fmt.Errorf("id is empty - see %s", tmpFile)
	}
	return result.Id, nil
}

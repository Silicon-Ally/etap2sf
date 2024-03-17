package client

import (
	"fmt"
	"strings"

	"github.com/Silicon-Ally/etap2sf/salesforce"
)

func (c *Client) deleteAllIDs(ids []string) error {
	batches := [][]string{}
	for len(ids) > 200 {
		batches = append(batches, ids[:200])
		ids = ids[200:]
	}
	if len(ids) > 0 {
		batches = append(batches, ids)
	}
	for i, batch := range batches {
		resp, err := c.gc.EnterpriseClient.Delete(batch)
		if err != nil {
			return fmt.Errorf("deleting batch %d/%d: %w", i+1, len(batches), err)
		}
		for _, result := range resp {
			if !result.Success {
				asStr := fmt.Sprintf("%+v", result.Errors[0])
				if strings.Contains(asStr, "entity is deleted") {
					continue
				}
				return fmt.Errorf("delete failed: %+v", result.Errors[0])
			}
		}
	}
	return nil
}

func (c *Client) getAllIDs(sot salesforce.ObjectType) ([]string, error) {
	sn, err := sot.SalesforceName()
	if err != nil {
		return nil, fmt.Errorf("getting salesforce name: %w", err)
	}
	ids := []string{}
	query := "SELECT Id FROM " + sn
	response, err := c.gc.EnterpriseClient.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying: %w", err)
	}
	for _, record := range response.Records {
		if record.Id != "068Ho00000M6HMsIAN" && record.Id != "001Ho000017vEpGIAU" {
			ids = append(ids, record.Id)
		}
	}
	done := response.Done
	cursor := response.QueryLocator
	for !done {
		response, err := c.gc.EnterpriseClient.QueryMore(cursor)
		if err != nil {
			return nil, fmt.Errorf("query more: %w", err)
		}
		for _, record := range response.Records {
			ids = append(ids, record.Id)
		}
		done = response.Done
		cursor = response.QueryLocator
	}
	return ids, nil
}

func (c *Client) deleteAll(sot salesforce.ObjectType) error {
	ids, err := c.getAllIDs(sot)
	if err != nil {
		return fmt.Errorf("getting all ids: %w", err)
	}
	return c.deleteAllIDs(ids)
}

func (c *Client) DeleteAllExistingRelationships() error {
	ids, err := c.getAllIDs(salesforce.ObjectType_Relationship)
	if err != nil {
		return fmt.Errorf("getting all ids: %w", err)
	}
	return c.deleteAllIDs(ids)
}

func (c *Client) DeleteRelationshipsNotCreatedThroughETap() error {
	ids := []string{}
	query := "SELECT Id, Etap_Relationship_Ref__c FROM Npe4__Relationship__c WHERE Etap_Relationship_Ref__c = NULL"
	response, err := c.gc.EnterpriseClient.Query(query)
	if err != nil {
		return fmt.Errorf("querying relationships: %w", err)
	}
	for _, record := range response.Records {
		ids = append(ids, record.Id)
	}
	done := response.Done
	cursor := response.QueryLocator

	for !done {
		response, err := c.gc.EnterpriseClient.QueryMore(cursor)
		if err != nil {
			return fmt.Errorf("querying relationships: %w", err)
		}
		for _, record := range response.Records {
			ids = append(ids, record.Id)
		}
		done = response.Done
		cursor = response.QueryLocator
	}
	return c.deleteAllIDs(ids)
}

func (c *Client) DeleteAll() error {
	order := []salesforce.ObjectType{
		// salesforce.ObjectType_ContentDocumentLink,
		// salesforce.ObjectType_ContentVersion,
		salesforce.ObjectType_Task,
		salesforce.ObjectType_AdditionalContext,
		salesforce.ObjectType_AccountSoftCredit,
		salesforce.ObjectType_PartialSoftCredit,
		salesforce.ObjectType_GAUAllocation,
		salesforce.ObjectType_Payment,
		salesforce.ObjectType_Opportunity,
		salesforce.ObjectType_RecurringDonation,
		salesforce.ObjectType_Affiliation,
		salesforce.ObjectType_Relationship,
		salesforce.ObjectType_Contact,
		salesforce.ObjectType_Account,
		salesforce.ObjectType_GAUAllocation,
		salesforce.ObjectType_Campaign,
	}

	for _, sot := range order {
		fmt.Printf("Starting to Delete %q...", sot)
		if err := c.deleteAll(sot); err != nil {
			return fmt.Errorf("deleting %q: %w", sot, err)
		}
	}
	return nil
}

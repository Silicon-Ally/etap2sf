package client

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tzmfreedom/go-soapforce"
)

func (c *Client) FixNPSPTriggers() error {
	triggers, err := c.getAllNPSPTriggers()
	if err != nil {
		return fmt.Errorf("getting all npsp triggers: %w", err)
	}

	groupedByClass := map[string][]*npspTrigger{}
	for _, trigger := range triggers {
		groupedByClass[trigger.Class] = append(groupedByClass[trigger.Class], trigger)
	}

	sortByCreatedAt := func(triggers []*npspTrigger) {
		sort.Slice(triggers, func(i, j int) bool {
			return triggers[i].CreatedAt < triggers[j].CreatedAt
		})
	}

	toDelete := []*npspTrigger{}
	for _, triggers := range groupedByClass {
		sortByCreatedAt(triggers)
		toDelete = append(toDelete, triggers[1:]...)
	}

	fmt.Printf("to delete: %+v", toDelete)

	maxN := 100
	batches := [][]string{}
	for len(toDelete) > maxN {
		ids := []string{}
		for _, trigger := range toDelete[:maxN] {
			ids = append(ids, trigger.ID)
		}
		batches = append(batches, ids)
		toDelete = toDelete[maxN:]
	}
	ids := []string{}
	for _, trigger := range toDelete {
		ids = append(ids, trigger.ID)
	}
	batches = append(batches, ids)
	for _, batch := range batches {
		resp, err := c.gc.EnterpriseClient.Delete(batch)
		if err != nil {
			return fmt.Errorf("deleting triggers: %w", err)
		}
		for _, result := range resp {
			if !result.Success {
				asStr := fmt.Sprintf("%+v", result.Errors[0])
				if strings.Contains(asStr, "entity is deleted") {
					continue
				}
				return fmt.Errorf("deleting triggers failed for result: %+v", result.Errors[0])
			}
		}
	}
	return nil
}

func (c *Client) DisableNPSPRelationshipTriggers() (func() error, error) {
	triggers, err := c.getAllNPSPTriggers()
	if err != nil {
		return nil, fmt.Errorf("getting all npsp triggers: %w", err)
	}
	toDisable := []*npspTrigger{}
	for _, trigger := range triggers {
		if strings.Contains(trigger.Class, "Relationship") {
			toDisable = append(toDisable, trigger)
		}
	}
	if err := c.setNPSPTriggerActivityState(toDisable, false); err != nil {
		return nil, fmt.Errorf("setting npsp trigger activity state to inactive: %w", err)
	}
	undoFn := func() error {
		return c.setNPSPTriggerActivityState(toDisable, true)
	}
	return undoFn, nil
}

func (c *Client) setNPSPTriggerActivityState(triggers []*npspTrigger, active bool) error {
	sobjs := []*soapforce.SObject{}
	for _, trigger := range triggers {
		fields := map[string]interface{}{"npsp__Active__c": fmt.Sprintf("%t", active)}
		sobjs = append(sobjs, &soapforce.SObject{
			Id:     trigger.ID,
			Fields: fields,
			Type:   "npsp__Trigger_Handler__c",
		})
	}
	resp, err := c.gc.EnterpriseClient.Update(sobjs)
	if err != nil {
		return fmt.Errorf("updating npsp triggers: %w", err)
	}
	for _, result := range resp {
		if result.Success {
			continue
		}
		if *result.Errors[0].StatusCode == "ENTITY_IS_DELETED" {
			continue
		}
		return fmt.Errorf("updating npsp triggers failed for result: %+v", result.Errors[0])
	}
	return nil
}

type npspTrigger struct {
	ID        string
	Name      string
	Class     string
	CreatedAt string
}

func (c *Client) getAllNPSPTriggers() ([]*npspTrigger, error) {
	response, err := c.gc.EnterpriseClient.QueryAll("SELECT Id, Name, npsp__Class__c, CreatedDate FROM npsp__Trigger_Handler__c")
	if err != nil {
		return nil, fmt.Errorf("querying npsp triggers: %w", err)
	}

	triggers := []*npspTrigger{}
	for _, trigger := range response.Records {
		triggers = append(triggers, &npspTrigger{
			ID:        trigger.Id,
			Name:      trigger.Fields["Name"].(string),
			Class:     trigger.Fields["npsp__Class__c"].(string),
			CreatedAt: trigger.Fields["CreatedDate"].(string),
		})
	}
	return triggers, nil
}

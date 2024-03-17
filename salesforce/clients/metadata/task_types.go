package client

import (
	"encoding/xml"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
)

func (c *Client) AddCustomTaskTypes(valuesToAdd []string) error {
	type ReadResult struct {
		Records []*sfmetadata.CustomField `xml:"records,omitempty"`
	}

	type ReadMetadataResponse struct {
		XMLName xml.Name `xml:"http://soap.sforce.com/2006/04/metadata readMetadataResponse"`

		Result *ReadResult `xml:"result,omitempty"`
	}

	for _, fullName := range []string{"Task.Subject", "Task.Type"} {
		resp := &ReadMetadataResponse{}
		err := c.gc.MetadataClient.ReadMetadataInto("CustomField", []string{fullName}, resp)
		if err != nil {
			return fmt.Errorf("reading metadata: %w", err)
		}
		if resp.Result == nil || len(resp.Result.Records) == 0 {
			return fmt.Errorf("no fields found with name %q", fullName)
		}
		if len(resp.Result.Records) > 1 {
			return fmt.Errorf("multiple fields found with name %q", fullName)
		}
		if resp.Result.Records[0] == nil {
			return fmt.Errorf("no field found with name %q", fullName)
		}
		field := resp.Result.Records[0]
		if field.ValueSet == nil || field.ValueSet.ValueSetDefinition == nil || len(field.ValueSet.ValueSetDefinition.Value) == 0 {
			return fmt.Errorf("something is wrong with field %q - it appears to be empty", fullName)
		}
		existingValues := map[string]bool{}
		newValues := map[string]bool{}
		for _, v := range field.ValueSet.ValueSetDefinition.Value {
			existingValues[v.Metadata.FullName] = true
		}
		for _, v := range valuesToAdd {
			if existingValues[v] {
				continue
			}
			newValues[v] = true
		}

		for v := range newValues {
			field.ValueSet.ValueSetDefinition.Value = append(field.ValueSet.ValueSetDefinition.Value, &sfmetadata.CustomValue{
				Description: fmt.Sprintf("Relationship of type %s, automatically ported in from eTapestry", v),
				IsActive:    true,
				Metadata: &sfmetadata.Metadata{
					FullName: v,
				},
			})
		}
		field.FullName = relationshipTypeFieldName
		if err := c.UpsertCustomField(salesforce.ObjectType_Task, field); err != nil {
			return fmt.Errorf("upserting custom field: %w", err)
		}
	}
	return nil
}

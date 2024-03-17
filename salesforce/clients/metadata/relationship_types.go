package client

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const relationshipTypeFieldName = "npe4__Type__c"

func (c *Client) AddRelationshipTypesToPicklist(valuesToAdd []string) error {
	sot := salesforce.ObjectType_Relationship
	sots, err := sot.SalesforceName()
	if err != nil {
		return fmt.Errorf("getting salesforce name for relationship: %w", err)
	}
	fullName := sots + "." + relationshipTypeFieldName

	type ReadResult struct {
		Records []*sfmetadata.CustomField `xml:"records,omitempty"`
	}

	type ReadMetadataResponse struct {
		XMLName xml.Name `xml:"http://soap.sforce.com/2006/04/metadata readMetadataResponse"`

		Result *ReadResult `xml:"result,omitempty"`
	}

	resp := &ReadMetadataResponse{}
	err = c.gc.MetadataClient.ReadMetadataInto("CustomField", []string{fullName}, resp)
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
		existingValues[standardizeRelationshipName(v.Metadata.FullName)] = true
	}
	for _, v := range valuesToAdd {
		if existingValues[v] {
			continue
		}
		newValues[standardizeRelationshipName(v)] = true
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
	if err := c.UpsertCustomField(salesforce.ObjectType_Relationship, field); err != nil {
		return fmt.Errorf("upserting custom field: %w", err)
	}
	return nil
}

func standardizeRelationshipName(s string) string {
	toTitle := cases.Title(language.English)

	s = strings.ReplaceAll(s, "-", " ")
	splits := strings.Split(s, " ")
	for i := range splits {
		splits[i] = toTitle.String(splits[i])
	}
	s = strings.Join(splits, " ")
	s = strings.ReplaceAll(s, " In Law", "-in-Law")
	s = strings.ReplaceAll(s, " Or ", " or ")
	s = strings.ReplaceAll(s, " Of ", " of ")
	return s
}

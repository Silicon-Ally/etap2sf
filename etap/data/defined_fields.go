package data

import (
	"encoding/json"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var definedFields []*generated.DefinedField

func GetDefinedFields() ([]*generated.DefinedField, error) {
	if definedFields != nil {
		return definedFields, nil
	}
	data, err := utils.MemoizeOperation("etap-defined-fields.json", doGetDefinedFields)
	if err != nil {
		return nil, fmt.Errorf("failed to get defined field data: %v", err)
	}
	result := []*generated.DefinedField{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaign data: %v", err)
	}
	refs := map[string]bool{}
	for _, a := range result {
		if a.Ref == nil || *a.Ref == "" {
			return nil, fmt.Errorf("defined field %s has no ref", *a.Ref)
		}
		if _, ok := refs[*a.Ref]; ok {
			return nil, fmt.Errorf("duplicate defined field ref: %s", *a.Ref)
		}
		refs[*a.Ref] = true
	}
	definedFields = result
	return result, nil
}

func doGetDefinedFields() ([]byte, error) {
	return client.WithClient(func(c *client.Client) ([]byte, error) {
		dfs, err := c.GetAllDefinedFields()
		if err != nil {
			return nil, fmt.Errorf("failed to get defined fields: %v", err)
		}
		result, err := json.MarshalIndent(dfs, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal defined fields: %v", err)
		}
		return result, nil
	})
}

package data

import (
	"encoding/json"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var funds []*generated.Fund

func GetFunds() ([]*generated.Fund, error) {
	if funds != nil {
		return funds, nil
	}
	fData, err := utils.MemoizeOperation("etap-funds.json", doGetFundData)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund data: %v", err)
	}
	result := []*generated.Fund{}
	if err := json.Unmarshal(fData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fund data: %v", err)
	}
	refs := map[string]bool{}
	for _, f := range result {
		if f.Ref == nil || *f.Ref == "" {
			return nil, fmt.Errorf("account %s has no ref", *f.Ref)
		}
		if _, ok := refs[*f.Ref]; ok {
			return nil, fmt.Errorf("duplicate account ref: %s", *f.Ref)
		}
		refs[*f.Ref] = true
	}
	funds = result
	return result, nil
}

func doGetFundData() ([]byte, error) {
	return client.WithClient(func(c *client.Client) ([]byte, error) {
		funds, err := c.GetAllFunds()
		if err != nil {
			return nil, fmt.Errorf("failed to get funds : %v", err)
		}
		result, err := json.MarshalIndent(funds, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal funds: %v", err)
		}
		return result, nil
	})
}

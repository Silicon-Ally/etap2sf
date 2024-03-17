package data

import (
	"encoding/json"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var approaches []string

func GetApproaches() ([]string, error) {
	if approaches != nil {
		return approaches, nil
	}
	aData, err := utils.MemoizeOperation("etap-approaches.json", doGetApproachData)
	if err != nil {
		return nil, fmt.Errorf("failed to get approach data: %v", err)
	}
	result := []string{}
	if err := json.Unmarshal(aData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal approach data: %v", err)
	}
	uniq := map[string]bool{}
	for _, a := range result {
		if _, ok := uniq[a]; ok {
			return nil, fmt.Errorf("duplicate approach: %s", a)
		}
		uniq[a] = true
	}
	approaches = result
	return result, nil
}

func doGetApproachData() ([]byte, error) {
	return client.WithClient(func(c *client.Client) ([]byte, error) {
		approaches, err := c.GetAllApproaches()
		if err != nil {
			return nil, fmt.Errorf("failed to get approaches: %v", err)
		}
		result, err := json.MarshalIndent(approaches, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal approaches: %v", err)
		}
		return result, nil
	})
}

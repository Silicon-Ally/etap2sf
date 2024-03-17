package data

import (
	"encoding/json"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var campaigns []string

func GetCampaigns() ([]string, error) {
	if campaigns != nil {
		return campaigns, nil
	}
	cData, err := utils.MemoizeOperation("etap-campaigns.json", doGetCampaignData)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign data: %v", err)
	}
	result := []string{}
	if err := json.Unmarshal(cData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaign data: %v", err)
	}
	uniq := map[string]bool{}
	for _, a := range result {
		if _, ok := uniq[a]; ok {
			return nil, fmt.Errorf("duplicate campaign: %s", a)
		}
		uniq[a] = true
	}
	campaigns = result
	return result, nil
}

func doGetCampaignData() ([]byte, error) {
	return client.WithClient(func(c *client.Client) ([]byte, error) {
		campaigns, err := c.GetAllCampaigns()
		if err != nil {
			return nil, fmt.Errorf("failed to get campaigns: %v", err)
		}
		result, err := json.MarshalIndent(campaigns, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal campaigns: %v", err)
		}
		return result, nil
	})
}

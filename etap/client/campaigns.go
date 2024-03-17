package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetAllCampaigns() ([]string, error) {
	request := struct {
		M generated.OperationMessagingService_getCampaigns `xml:"tns:getCampaigns"`
	}{
		generated.OperationMessagingService_getCampaigns{
			Boolean_1: ptr(true), // Include disabled
		},
	}
	result := overrides.CampaignsBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetCampaigns", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return result.M.Result.Items, nil
}

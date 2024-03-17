package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetAllFunds() ([]*generated.Fund, error) {
	request := struct {
		M generated.OperationMessagingService_getFundObjects `xml:"tns:getFundObjects"`
	}{
		generated.OperationMessagingService_getFundObjects{
			Boolean_1: ptr(true), // Include disabled
		},
	}
	result := overrides.FundObjectsBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetFunds", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return result.M.Result.Items, nil
}

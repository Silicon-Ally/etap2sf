package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetAllApproaches() ([]string, error) {
	request := struct {
		M generated.OperationMessagingService_getApproaches `xml:"tns:getApproaches"`
	}{
		generated.OperationMessagingService_getApproaches{
			Boolean_1: ptr(true), // Include disabled
		},
	}
	result := overrides.ApproachesBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetApproaches", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return result.M.Result.Items, nil
}

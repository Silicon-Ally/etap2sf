package client

import (
	"fmt"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetDefinedFields(start, count int) ([]*generated.DefinedField, error) {
	request := struct {
		M generated.OperationMessagingService_getDefinedFields `xml:"tns:getDefinedFields"`
	}{
		generated.OperationMessagingService_getDefinedFields{
			PagedDefinedFieldsRequest_1: &generated.PagedDefinedFieldsRequest{
				ClearCache:            ptr(false),
				Start:                 ptr(start),
				Count:                 ptr(count),
				IncludeDisabledFields: ptr(true),
				IncludeDisabledValues: ptr(true),
			},
		},
	}
	result := overrides.DefinedFieldsBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetDefinedFields", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return ptrs(result.M.Result.Data.Items), nil
}

func (c *Client) GetAllDefinedFields() ([]*generated.DefinedField, error) {
	dfs := []*generated.DefinedField{}
	start := 0
	count := 100
	for {
		page, err := c.GetDefinedFields(start, count)
		if err != nil {
			return nil, fmt.Errorf("failed to get defined fields: %v", err)
		}
		dfs = append(dfs, page...)
		start += count
		if len(page) < count {
			break
		}
		fmt.Printf("Successfully Processed %d DefinedFields\n", start)
		time.Sleep(time.Second * 1)
	}
	return dfs, nil
}

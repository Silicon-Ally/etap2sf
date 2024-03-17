package client

import (
	"fmt"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetRelationships(account *generated.Account, start, count int) ([]*generated.Relationship, error) {
	request := struct {
		M generated.OperationMessagingService_getRelationships `xml:"tns:getRelationships"`
	}{
		generated.OperationMessagingService_getRelationships{
			PagedRelationshipsRequest_1: &generated.PagedRelationshipsRequest{
				ClearCache: ptr(false),
				Start:      ptr(start),
				Count:      ptr(count),
				AccountRef: account.Ref,
			},
		},
	}
	result := overrides.RelationshipsBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetRelationships", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return ptrs(result.M.Result.Data.Items), nil
}

func (c *Client) GetAllRelationships(account *generated.Account) ([]*generated.Relationship, error) {
	rs := []*generated.Relationship{}
	start := 0
	count := 100
	for {
		page, err := c.GetRelationships(account, start, count)
		if err != nil {
			return nil, fmt.Errorf("failed to get relationships: %v", err)
		}
		rs = append(rs, page...)
		start += count
		if len(page) < count {
			break
		}
		fmt.Printf("Successfully Processed %d Relationships\n", start)
		time.Sleep(time.Second * 1)
	}
	return rs, nil
}

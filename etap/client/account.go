package client

import (
	"fmt"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetAccounts(queryName string, start int, count int) ([]*generated.Account, error) {
	request := struct {
		M generated.OperationMessagingService_getExistingQueryResults `xml:"tns:getExistingQueryResults"`
	}{
		generated.OperationMessagingService_getExistingQueryResults{
			PagedExistingQueryResultsRequest_1: &generated.PagedExistingQueryResultsRequest{
				ClearCache: ptr(false),
				Start:      ptr(start),
				Count:      ptr(count),
				Query:      ptr(queryName),
			},
		},
	}
	result := overrides.AccountBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetExistingQueryResults", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return ptrs(result.M.Result.Data.Items), nil
}

func (c *Client) GetAllAccounts(queryName string) ([]*generated.Account, error) {
	accounts := []*generated.Account{}
	start := 0
	count := 100
	for {
		page, err := c.GetAccounts(queryName, start, count)
		if err != nil {
			return nil, fmt.Errorf("failed to get accounts: %v", err)
		}
		accounts = append(accounts, page...)
		start += count
		if len(page) < count {
			break
		}
		fmt.Printf("Successfully Processed %d Accounts\n", start)
		time.Sleep(time.Second * 5)
	}
	return accounts, nil
}

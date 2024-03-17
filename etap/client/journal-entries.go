package client

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
)

func (c *Client) GetJournalEntries(account *generated.Account, start, count int) ([]*overrides.JournalEntry, error) {
	request := struct {
		M            generated.OperationMessagingService_getJournalEntries `xml:"tns:getJournalEntries"`
		XmlnsSOAPENC string                                                `xml:"xmlns:soapenc,attr"`
		XmlnsXSD     string                                                `xml:"xmlns:xsd,attr"`
	}{
		M: generated.OperationMessagingService_getJournalEntries{
			PagedJournalEntriesRequest_1: &generated.PagedJournalEntriesRequest{
				ClearCache: ptr(false),
				Start:      ptr(start),
				Count:      ptr(count),
				AccountRef: account.Ref,
				Types: &generated.ArrayOfint{
					ArrayType: xml.Attr{
						Name:  xml.Name{Space: "http://schemas.xmlsoap.org/soap/encoding/", Local: "arrayType"},
						Value: "xsd:int[17]",
					},
					Items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 18, 19},
				},
			},
		},
		XmlnsSOAPENC: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsXSD:     "http://www.w3.org/2001/XMLSchema",
	}
	result := overrides.JournalEntryBody{}
	if err := generated.RoundTripWithAction(c.ms, "GetJournalEntries", request, &result); err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	if c.err != nil {
		return nil, fmt.Errorf("fault code error: %v", c.err)
	}
	return ptrs(result.M.Result.Data.Items), nil
}

func (c *Client) GetAllJournalEntries(account *generated.Account) ([]*overrides.JournalEntry, error) {
	jes := []*overrides.JournalEntry{}
	start := 0
	count := 100
	for {
		page, err := c.GetJournalEntries(account, start, count)
		if err != nil {
			return nil, fmt.Errorf("failed to get journal entries: %v", err)
		}
		jes = append(jes, page...)
		start += count
		if len(page) < count {
			break
		}
		fmt.Printf("Successfully Processed %d Journal Entries\n", start)
		time.Sleep(time.Second * 1)
	}
	return jes, nil
}

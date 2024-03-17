package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
)

func (c *Client) GetAccountRecordTypes() (orgID, hhID sfenterprise.ID, err error) {
	resp, err := c.gc.EnterpriseClient.QueryAll("SELECT Id, Name, SobjectType FROM RecordType")
	if err != nil {
		err = fmt.Errorf("querying content document by version: %w", err)
		return
	}
	for _, r := range resp.Records {
		if r.Fields["SobjectType"] == "Account" && r.Fields["Name"] == "Organization" {
			orgID = sfenterprise.ID(r.Id)
		}
		if r.Fields["SobjectType"] == "Account" && r.Fields["Name"] == "Household Account" {
			hhID = sfenterprise.ID(r.Id)
		}
	}
	if orgID == "" {
		err = fmt.Errorf("could not find organization record type")
		return
	}
	if hhID == "" {
		err = fmt.Errorf("could not find household record type")
		return
	}
	return
}

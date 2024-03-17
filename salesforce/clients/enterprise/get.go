package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
)

func (c *Client) LookupContentDocumentByVersion(id sfenterprise.ID) (sfenterprise.ID, error) {
	response, err := c.gc.EnterpriseClient.QueryAll("SELECT ContentDocumentId FROM ContentVersion WHERE Id = '" + string(id) + "'")
	if err != nil {
		return "", fmt.Errorf("querying content document by version: %w", err)
	}
	ids := []string{}
	for _, record := range response.Records {
		cdID := record.Fields["ContentDocumentId"]
		if cdID == nil || cdID.(string) == "" {
			return "", fmt.Errorf("no content document found for version %q", id)
		}
		ids = append(ids, cdID.(string))
	}
	if len(ids) == 0 {
		return "", fmt.Errorf("no content document found for version %q", id)
	}
	if len(ids) > 1 {
		return "", fmt.Errorf("multiple content documents found for version %q", id)
	}
	return sfenterprise.ID(ids[0]), nil
}

func (c *Client) LookupUserByEmail(email string) (sfenterprise.ID, error) {
	response, err := c.gc.EnterpriseClient.QueryAll("SELECT Id FROM User WHERE Email = '" + email + "'")
	if err != nil {
		return "", fmt.Errorf("querying user by email: %w", err)
	}
	ids := []string{}
	for _, record := range response.Records {
		cdID := record.Id
		if cdID == "" {
			return "", fmt.Errorf("no id found for user with email %q", email)
		}
		ids = append(ids, cdID)
	}
	if len(ids) == 0 {
		return "", fmt.Errorf("no user found with email %q", email)
	}
	if len(ids) > 1 {
		return "", fmt.Errorf("multiple users found for email %q", email)
	}
	return sfenterprise.ID(ids[0]), nil
}

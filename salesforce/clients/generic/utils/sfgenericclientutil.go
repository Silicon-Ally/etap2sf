package utils

import (
	"fmt"

	genericclient "github.com/Silicon-Ally/etap2sf/salesforce/clients/generic"
	"github.com/Silicon-Ally/etap2sf/secrets"
)

func NewGenericClient() (*genericclient.Client, error) {
	connConfig, err := secrets.GetSalesforceConnectionConfig()
	if err != nil {
		return nil, fmt.Errorf("getting salesforce sandbox credentials: %w", err)
	}
	client, err := genericclient.New(&genericclient.Config{
		APIVersion: "58.0",
		Debug:      true,
		ConnConfig: connConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}
	return client, nil
}

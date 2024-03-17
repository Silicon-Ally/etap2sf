package utils

import (
	"fmt"

	client "github.com/Silicon-Ally/etap2sf/salesforce/clients/enterprise"
	"github.com/Silicon-Ally/etap2sf/secrets"
)

func NewSandboxClient() (*client.Client, error) {
	connConfig, err := secrets.GetSalesforceConnectionConfig()
	if err != nil {
		return nil, fmt.Errorf("getting salesforce sandbox credentials: %w", err)
	}
	client, err := client.New(&client.Config{
		APIVersion: "58.0",
		Debug:      true,
		ConnConfig: connConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}
	return client, nil
}

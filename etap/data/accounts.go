package data

import (
	"encoding/json"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var accounts []*generated.Account

func GetAccounts() ([]*generated.Account, error) {
	if accounts != nil {
		return accounts, nil
	}
	aData, err := utils.MemoizeOperation("etap-accounts.json", doGetAccountData)
	if err != nil {
		return nil, fmt.Errorf("failed to get account data: %v", err)
	}
	result := []*generated.Account{}
	if err := json.Unmarshal(aData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account data: %v", err)
	}
	accounts = result
	return result, nil
}

// In order for this to work, you'll need to first create a query for all accounts you want.
const GetAllAccountsQuery = "Folder::Query"
const originalQuery = "Folder::Query"

func doGetAccountData() ([]byte, error) {
	if GetAllAccountsQuery == originalQuery {
		return nil, fmt.Errorf(`

This step likely failed because there isn't a query in eTapestry to export all of your accounts.

To do this, create a new report query in eTapestry which exports all of your accounts (only the ID is needed).

Then, update the GetAllAccountsQuery constant in etap/data/accounts.go to use the name of your new query, noting that the folder will need to be specified before the ::.

If you have any problems with this, please file a bug.

		`)
	}
	return client.WithClient(func(c *client.Client) ([]byte, error) {
		accounts, err := c.GetAllAccounts(GetAllAccountsQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to get accounts: %v", err)
		}
		result, err := json.MarshalIndent(accounts, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal accounts: %v", err)
		}
		return result, nil
	})
}

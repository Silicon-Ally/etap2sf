package data

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/client"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/utils"
	"github.com/google/go-cmp/cmp"
)

var relationships []*generated.Relationship

func GetRelationships() ([]*generated.Relationship, error) {
	if relationships != nil {
		return relationships, nil
	}
	data, err := utils.MemoizeOperation("etap-relationships.json", doGetRelationshipData)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship data: %v", err)
	}
	result := []*generated.Relationship{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal relationship data: %v", err)
	}
	unique := map[string]*generated.Relationship{}
	for _, r := range result {
		if r.Ref == nil || *r.Ref == "" {
			return nil, fmt.Errorf("relationship %s has no ref", *r.Account1Ref+*r.Account2Ref)
		}
		if otherRelationship, ok := unique[*r.Ref]; ok {
			if diff := cmp.Diff(otherRelationship, r); diff != "" {
				return nil, fmt.Errorf("unexpected diff in relationships with the same ref:\n%s", diff)
			}
		}
		unique[*r.Ref] = r
	}
	// Determistic for testing
	uids := []string{}
	for uid := range unique {
		uids = append(uids, uid)
	}
	sort.Strings(uids)
	result = []*generated.Relationship{}
	for _, uid := range uids {
		result = append(result, unique[uid])
	}
	relationships = result
	return result, nil
}

func doGetRelationshipData() ([]byte, error) {
	accounts, err := GetAccounts()
	if err != nil {
		return nil, fmt.Errorf("getting accounts: %w", err)
	}

	return client.WithClient(func(c *client.Client) ([]byte, error) {
		getRelationships := func(account *generated.Account) ([]*generated.Relationship, error) {
			fileName := fmt.Sprintf("relationships/%s.json", *account.Ref)
			rData, err := utils.MemoizeOperation(fileName, func() ([]byte, error) {
				rs, err := c.GetAllRelationships(account)
				if err != nil {
					return nil, fmt.Errorf("getting relationships for %q: %w", *account.Ref, err)
				}
				data, err := json.MarshalIndent(rs, "", "  ")
				if err != nil {
					return nil, fmt.Errorf("marshaling relationships: %w", err)
				}
				// Don't hammer the server
				time.Sleep(100 * time.Millisecond)
				return data, nil
			})
			if err != nil {
				return nil, fmt.Errorf("memoizing operation: %w", err)
			}
			rs := []*generated.Relationship{}
			if err := json.Unmarshal(rData, &rs); err != nil {
				return nil, fmt.Errorf("unmarshaling relationships: %w", err)
			}
			return rs, nil
		}

		rss := []*generated.Relationship{}
		for i, account := range accounts {
			rs, err := getRelationships(account)
			if err != nil {
				return nil, fmt.Errorf("getting relationships: %w", err)
			}
			rss = append(rss, rs...)
			fmt.Printf("Relationships for account %d/%d: %s: #=%d\n", i+1, len(accounts), *account.Name, len(rs))
		}

		result, err := json.MarshalIndent(rss, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal relationships: %v", err)
		}

		fmt.Printf("Completed downloading relationship data!\n")
		return result, nil
	})
}

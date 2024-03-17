package create_sf_fields

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Silicon-Ally/etap2sf/utils"
)

type errorz struct {
	Errors []string
	Ids    map[string]bool
}

var errsPath = filepath.Join(utils.ProjectRoot(), "data", "salesforce-create-fields-errors.json")

func readErrors() (*errorz, error) {
	data, err := os.ReadFile(errsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &errorz{Ids: map[string]bool{}}, nil
		}
		return nil, fmt.Errorf("reading errors: %w", err)
	}
	result := &errorz{Ids: map[string]bool{}}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("unmarshaling errors: %w", err)
	}
	return result, nil
}

func (e *errorz) Write() error {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling errors: %w", err)
	}
	if err := os.WriteFile(errsPath, data, 0644); err != nil {
		return fmt.Errorf("writing errors: %w", err)
	}
	return nil
}

func (e *errorz) sortIdx(id string) int {
	b, ok := e.Ids[id]
	if !ok {
		return 1
	}
	if b {
		return 0
	}
	return 2
}

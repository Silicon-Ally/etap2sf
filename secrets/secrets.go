package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Silicon-Ally/etap2sf/utils"
)

func GetETapAPIKey() (string, error) {
	return read("etapestry-api-key.txt")
}

func GetETapDBName() (string, error) {
	return read("etapestry-db-name.txt")
}

type SalesforceConnectionConfig struct {
	Username      string
	Password      string
	SecurityToken string
	LoginURL      string
}

func (scc *SalesforceConnectionConfig) GetUsername() string      { return scc.Username }
func (scc *SalesforceConnectionConfig) GetPassword() string      { return scc.Password }
func (scc *SalesforceConnectionConfig) GetSecurityToken() string { return scc.SecurityToken }
func (scc *SalesforceConnectionConfig) GetLoginURL() string      { return scc.LoginURL }

func GetSalesforceConnectionConfig() (*SalesforceConnectionConfig, error) {
	fileName := "salesforce-connection-config.txt"
	all, err := read(fileName)
	if err != nil {
		return nil, err
	}
	split := strings.Split(strings.TrimSpace(all), "\n")
	if len(split) != 4 {
		return nil, fmt.Errorf("expected 4 lines in %s, got %d", fileName, len(split))
	}
	return &SalesforceConnectionConfig{
		Username:      split[0],
		Password:      split[1],
		SecurityToken: split[2],
		LoginURL:      split[3],
	}, nil
}

func read(fileName string) (string, error) {
	filePath := filepath.Join(utils.ProjectRoot(), "secrets", fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read secrets/%s: %v", fileName, err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("secrets/%s exists, but was empty", fileName)
	}
	return strings.TrimSpace(string(data)), nil
}

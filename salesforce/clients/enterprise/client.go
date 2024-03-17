package client

import (
	"fmt"

	genericclient "github.com/Silicon-Ally/etap2sf/salesforce/clients/generic"
)

type Client struct {
	gc *genericclient.Client
}

type ConnConfig interface {
	GetUsername() string
	GetPassword() string
	GetSecurityToken() string
	GetLoginURL() string
}

type Config struct {
	ConnConfig ConnConfig
	APIVersion string
	Debug      bool
}

func New(c *Config) (*Client, error) {
	gc, err := genericclient.New(&genericclient.Config{
		APIVersion: c.APIVersion,
		Debug:      c.Debug,
		ConnConfig: c.ConnConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("creating generic client: %w", err)
	}
	return &Client{gc: gc}, nil
}

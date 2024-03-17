package genericclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/tzmfreedom/go-metaforce"
	"github.com/tzmfreedom/go-soapforce"
)

type Client struct {
	EnterpriseClient   *soapforce.Client
	MetadataClient     *metaforce.Client
	MetadataSOAPClient *soapforce.SOAPClient
	IDMap              map[string]string
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
	username := c.ConnConfig.GetUsername()
	password := c.ConnConfig.GetPassword()
	securityToken := c.ConnConfig.GetSecurityToken()
	loginURL := c.ConnConfig.GetLoginURL()
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if securityToken == "" {
		return nil, fmt.Errorf("security token is required")
	}
	if loginURL == "" {
		return nil, fmt.Errorf("login url is required")
	}
	if c.APIVersion == "" {
		return nil, fmt.Errorf("api version is required")
	}

	m := metaforce.NewClient()
	m.SetApiVersion(c.APIVersion)
	m.SetDebug(c.Debug)
	m.SetLoginUrl(loginURL)

	if err := m.Login(username, password+securityToken); err != nil {
		return nil, fmt.Errorf("failed to login to metaforce client: %w", err)
	}

	e := soapforce.NewClient()
	e.SetApiVersion(c.APIVersion)
	e.SetDebug(c.Debug)
	e.SetLoginUrl(loginURL)
	if _, err := e.Login(username, password+securityToken); err != nil {
		return nil, fmt.Errorf("failed to login to soapforce client: %w", err)
	}

	s := soapforce.NewSOAPClient(fmt.Sprintf("https://login.salesforce.com/services/Soap/u/%s", c.APIVersion), true, nil)

	return &Client{
		MetadataClient:     m,
		EnterpriseClient:   e,
		MetadataSOAPClient: s,
	}, nil
}

func (c *Client) DownloadMetadataWSDL(toFilePath string) error {
	if err := os.MkdirAll(filepath.Dir(toFilePath), 0644); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	u, err := url.Parse(c.MetadataClient.GetServerURL())
	if err != nil {
		return fmt.Errorf("parsing server url: %w", err)
	}
	if u == nil {
		return fmt.Errorf("no server url parsed for %q", c.MetadataClient.GetServerURL())
	}
	u.Path = "/services/wsdl/metadata"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to format request: %w", err)
	}
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
	req.Header.Set("Accept-Language", `en-US,en;q=0.9`)
	req.Header.Set("Cache-Control", `max-age=0`)
	req.Header.Set("Connection", `keep-alive`)
	req.Header.Set("Cookie", fmt.Sprintf(`cleared-onetrust-cookies=; sid=%s;`, c.MetadataClient.GetSessionID()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to issue request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dat, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response code %d for metadata download request: %q", resp.StatusCode, string(dat))
	}

	f, err := os.Open(toFilePath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("failed to download WSDL to file: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close WSDL download file: %w", err)
	}

	return nil
}

func (c *Client) DownloadEnterpriseWSDL(toFilePath string) error {
	if err := os.MkdirAll(filepath.Dir(toFilePath), 0644); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	u, err := url.Parse(c.MetadataClient.GetServerURL())
	if err != nil {
		return fmt.Errorf("parsing server url: %w", err)
	}
	if u == nil {
		return fmt.Errorf("no server url parsed for %q", c.MetadataClient.GetServerURL())
	}
	u.Path = "/soap/wsdl.jsp?type=*"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to format request: %w", err)
	}
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
	req.Header.Set("Accept-Language", `en-US,en;q=0.9`)
	req.Header.Set("Cache-Control", `max-age=0`)
	req.Header.Set("Connection", `keep-alive`)
	req.Header.Set("Cookie", fmt.Sprintf(`cleared-onetrust-cookies=; sid=%s;`, c.MetadataClient.GetSessionID()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to issue request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dat, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response code %d for metadata download request: %q", resp.StatusCode, string(dat))
	}

	f, err := os.Open(toFilePath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("failed to download WSDL to file: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close WSDL download file: %w", err)
	}

	return nil
}

func (c *Client) ReadMetadataInto(request *metaforce.ReadMetadata, response any) error {
	err := c.MetadataSOAPClient.Call(request, response, &soapforce.ResponseSOAPHeader{})
	if err != nil {
		return err
	}
	return nil
}

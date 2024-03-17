package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/secrets"
	"github.com/Silicon-Ally/etap2sf/utils"
	"github.com/fiorix/wsdl2go/soap"
	"go.uber.org/multierr"
)

type Client struct {
	ms      generated.MessagingService
	cookies []*http.Cookie
	err     error
}

const InitialUrl string = "https://sna.etapestry.com/v3messaging/service?WSDL"

func NewClient(dbName, secretAPIKey string) (*Client, error) {
	url := InitialUrl
	redirectURL, cookies, err := attemptLogin(url, dbName, secretAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform initial login: %v", err)
	}
	if redirectURL != "" {
		url = redirectURL
		fmt.Printf("redirected to %s\n", url)
		redirectURL, cookies, err = attemptLogin(url, dbName, secretAPIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to perform secondary login: %v", err)
		}
		if redirectURL != "" {
			return nil, fmt.Errorf("redirected login tried to redirect again, to %s", redirectURL)
		}
	}
	if cookies == nil {
		return nil, fmt.Errorf("failed to get JSESSIONID cookie")
	}

	c := &Client{cookies: cookies}
	sc := &soap.Client{
		URL:       url,
		Namespace: generated.Namespace,
		Pre: func(r *http.Request) {
			if err := addRequiredSOAPEncodingStyle(r); err != nil {
				log.Printf("failed to add required SOAP encoding style: %v", err)
			}
			for _, cookie := range c.cookies {
				r.AddCookie(cookie)
			}
			body, _ := io.ReadAll(r.Body)

			// Restore the io.ReadCloser to its original state
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Use the body content
			file, _ := utils.WriteBytesToTempFile(body, "etap-request")
			log.Printf("\n%s request: %s ", r.URL.Path, file)
		},
		Post: func(resp *http.Response) {
			body, _ := io.ReadAll(resp.Body)
			if fileName, err := utils.WriteBytesToTempFile(body, "etap-response"); err != nil {
				c.err = err
			} else {
				fmt.Printf("response: %s\n", fileName)
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(body))

			if err := checkForFaultCode(resp); err != nil {
				c.err = err
			}
			if err := resolveXmlHrefsHttpResponse(resp); err != nil {
				c.err = err
			}
		},
	}
	c.ms = generated.NewMessagingService(sc)

	return c, nil
}

func WithClient(fn func(c *Client) ([]byte, error)) (dat []byte, rErr error) {
	secretAPIKey, err := secrets.GetETapAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to read API key: %v", err)
	}

	dbName, err := secrets.GetETapDBName()
	if err != nil {
		return nil, fmt.Errorf("failed to read db name: %v", err)
	}

	c, err := NewClient(dbName, secretAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer func() {
		if err := c.Logout(); err != nil {
			rErr = multierr.Append(rErr, fmt.Errorf("failed to logout: %w", err))
		}
	}()

	return fn(c)
}

func attemptLogin(url string, dbName string, secretAPIKey string) (string, []*http.Cookie, error) {
	var cookies []*http.Cookie
	var e error
	sc := &soap.Client{
		URL:       url,
		Namespace: generated.Namespace,
		Pre: func(r *http.Request) {
			if err := addRequiredSOAPEncodingStyle(r); err != nil {
				log.Printf("failed to add required SOAP encoding style: %v", err)
			}
		},
		Post: func(resp *http.Response) {
			if err := checkForFaultCode(resp); err != nil {
				e = err
			}
			if len(resp.Cookies()) > 0 {
				cookies = resp.Cookies()
			}
		},
	}
	ms := generated.NewMessagingService(sc)

	redirectURL, err := ms.ApiKeyLogin(dbName, secretAPIKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to perform initial login: %v", err)
	}
	if e != nil {
		return "", nil, fmt.Errorf("fault code error: %v", e)
	}
	if redirectURL != "" {
		return redirectURL, nil, nil
	}
	if len(cookies) == 0 {
		return "", nil, fmt.Errorf("failed to get JSESSIONID cookie")
	}
	return "", cookies, nil
}

func (c *Client) Logout() error {
	return c.ms.Logout()
}

func ptr[T any](v T) *T {
	return &v
}

func ptrs[T any](in []T) []*T {
	out := make([]*T, len(in))
	for i, t := range in {
		out[i] = ptr(t)
	}
	return out
}

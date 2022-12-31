// Package membership provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/do87/oapi-codegen version v0.5.0 DO NOT EDIT.
package membership

import (
	"net/url"
	"strings"

	common "github.com/SchwarzIT/community-stackit-go-client/internal/common"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/membership/v2.0/generated/membership"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/membership/v2.0/generated/private"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/membership/v2.0/generated/public"
)

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// list of connected client services
	Membership *membership.Client
	Public     *public.Client
	Private    *private.Client

	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client common.Client
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a factory client
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}

	client.Membership = membership.NewClient(server, client.Client)
	client.Public = public.NewClient(server, client.Client)
	client.Private = private.NewClient(server, client.Client)

	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer common.Client) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	Client *Client

	// list of connected client services
	Membership *membership.ClientWithResponses
	Public     *public.ClientWithResponses
	Private    *private.ClientWithResponses
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}

	cwr := &ClientWithResponses{Client: client}
	cwr.Membership = membership.NewClientWithResponses(server, client.Client)
	cwr.Public = public.NewClientWithResponses(server, client.Client)
	cwr.Private = private.NewClientWithResponses(server, client.Client)

	return cwr, nil
}

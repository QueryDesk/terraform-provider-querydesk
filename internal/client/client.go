//go:generate go run github.com/Khan/genqlient
package client

import (
	"fmt"
	"net/http"

	_ "github.com/Khan/genqlient/generate"
	"github.com/Khan/genqlient/graphql"
	_ "github.com/suessflorian/gqlfetch"
)

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-api-key", t.key)
	return t.wrapped.RoundTrip(req)
}

func NewClient(host *string, apiKey *string) (*graphql.Client, error) {
	httpClient := http.Client{
		Transport: &authedTransport{
			key:     *apiKey,
			wrapped: http.DefaultTransport,
		},
	}

	c := graphql.NewClient(fmt.Sprintf("%s/graphql", *host), &httpClient)

	return &c, nil
}

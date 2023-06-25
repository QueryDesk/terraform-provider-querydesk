//go:generate go run github.com/Khan/genqlient
package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
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

func NewGraphQLReq(ctx context.Context, client graphql.Client) *GraphQLReq {
	return &GraphQLReq{
		Context: ctx,
		Client:  client,
	}
}

//go:generate mockery --name GraphQLClient
type GraphQLClient interface {
	GetDatabase(id string) (*getDatabaseResponse, error)
	CreateDatabase(input CreateDatabaseInput) (*createDatabaseResponse, error)
	UpdateDatabase(id string, input UpdateDatabaseInput) (*updateDatabaseResponse, error)
	DeleteDatabase(id string) (*deleteDatabaseResponse, error)
}

type GraphQLReq struct {
	Context context.Context
	Client  graphql.Client
}

func (c GraphQLReq) GetDatabase(id string) (*getDatabaseResponse, error) {
	return getDatabase(c.Context, c.Client, id)
}

func (c GraphQLReq) CreateDatabase(input CreateDatabaseInput) (*createDatabaseResponse, error) {
	return createDatabase(c.Context, c.Client, input)
}

func (c GraphQLReq) UpdateDatabase(id string, input UpdateDatabaseInput) (*updateDatabaseResponse, error) {
	return updateDatabase(c.Context, c.Client, id, input)
}

func (c GraphQLReq) DeleteDatabase(id string) (*deleteDatabaseResponse, error) {
	return deleteDatabase(c.Context, c.Client, id)
}

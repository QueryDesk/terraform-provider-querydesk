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

//go:generate mockery --name GraphQLClient
type GraphQLClient interface {
	GetDatabase(ctx context.Context, id string) (*getDatabaseResponse, error)
	CreateDatabase(ctx context.Context, input CreateDatabaseInput) (*createDatabaseResponse, error)
	UpdateDatabase(ctx context.Context, id string, input UpdateDatabaseInput) (*updateDatabaseResponse, error)
	DeleteDatabase(ctx context.Context, id string) (*deleteDatabaseResponse, error)
}

type GraphQLReq struct {
	Client graphql.Client
}

func (c GraphQLReq) GetDatabase(ctx context.Context, id string) (*getDatabaseResponse, error) {
	return getDatabase(ctx, c.Client, id)
}

func (c GraphQLReq) CreateDatabase(ctx context.Context, input CreateDatabaseInput) (*createDatabaseResponse, error) {
	return createDatabase(ctx, c.Client, input)
}

func (c GraphQLReq) UpdateDatabase(ctx context.Context, id string, input UpdateDatabaseInput) (*updateDatabaseResponse, error) {
	return updateDatabase(ctx, c.Client, id, input)
}

func (c GraphQLReq) DeleteDatabase(ctx context.Context, id string) (*deleteDatabaseResponse, error) {
	return deleteDatabase(ctx, c.Client, id)
}

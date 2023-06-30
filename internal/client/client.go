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

type GetDatabaseResponse = getDatabaseResponse
type GetDatabaseDatabase = getDatabaseDatabase

type CreateDatabaseResponse = createDatabaseResponse
type CreateDatabaseCreateDatabaseCreateDatabaseResult = createDatabaseCreateDatabaseCreateDatabaseResult
type CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase = createDatabaseCreateDatabaseCreateDatabaseResultResultDatabase

type UpdateDatabaseResponse = updateDatabaseResponse
type UpdateDatabaseUpdateDatabaseUpdateDatabaseResult = updateDatabaseUpdateDatabaseUpdateDatabaseResult
type UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase = updateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase

type DeleteDatabaseResponse = deleteDatabaseResponse
type DeleteDatabaseDeleteDatabaseDeleteDatabaseResult = deleteDatabaseDeleteDatabaseDeleteDatabaseResult
type DeleteDatabaseDeleteDatabaseDeleteDatabaseResultResultDatabase = deleteDatabaseDeleteDatabaseDeleteDatabaseResultResultDatabase

type GetCredentialResponse = getCredentialResponse
type GetCredentialCredential = getCredentialCredential
type GetCredentialCredentialDatabase = getCredentialCredentialDatabase

type CreateCredentialResponse = createCredentialResponse
type CreateCredentialCreateCredentialCreateCredentialResult = createCredentialCreateCredentialCreateCredentialResult
type CreateCredentialCreateCredentialCreateCredentialResultResultCredential = createCredentialCreateCredentialCreateCredentialResultResultCredential

type UpdateCredentialResponse = updateCredentialResponse
type UpdateCredentialUpdateCredentialUpdateCredentialResult = updateCredentialUpdateCredentialUpdateCredentialResult
type UpdateCredentialUpdateCredentialUpdateCredentialResultResultCredential = updateCredentialUpdateCredentialUpdateCredentialResultResultCredential

type DeleteCredentialResponse = deleteCredentialResponse
type DeleteCredentialDeleteCredentialDeleteCredentialResult = deleteCredentialDeleteCredentialDeleteCredentialResult
type DeleteCredentialDeleteCredentialDeleteCredentialResultResultCredential = deleteCredentialDeleteCredentialDeleteCredentialResultResultCredential

//go:generate go run github.com/vektra/mockery/v2 --name GraphQLClient
type GraphQLClient interface {
	GetDatabase(ctx context.Context, id string) (*GetDatabaseResponse, error)
	CreateDatabase(ctx context.Context, input CreateDatabaseInput) (*CreateDatabaseResponse, error)
	UpdateDatabase(ctx context.Context, id string, input UpdateDatabaseInput) (*UpdateDatabaseResponse, error)
	DeleteDatabase(ctx context.Context, id string) (*DeleteDatabaseResponse, error)
	GetCredential(ctx context.Context, id string) (*GetCredentialResponse, error)
	CreateCredential(ctx context.Context, input CreateCredentialInput) (*CreateCredentialResponse, error)
	UpdateCredential(ctx context.Context, id string, input UpdateCredentialInput) (*UpdateCredentialResponse, error)
	DeleteCredential(ctx context.Context, id string) (*DeleteCredentialResponse, error)
}

type GraphQLReq struct {
	Client graphql.Client
}

func (c GraphQLReq) GetDatabase(ctx context.Context, id string) (*GetDatabaseResponse, error) {
	return getDatabase(ctx, c.Client, id)
}

func (c GraphQLReq) CreateDatabase(ctx context.Context, input CreateDatabaseInput) (*CreateDatabaseResponse, error) {
	return createDatabase(ctx, c.Client, input)
}

func (c GraphQLReq) UpdateDatabase(ctx context.Context, id string, input UpdateDatabaseInput) (*UpdateDatabaseResponse, error) {
	return updateDatabase(ctx, c.Client, id, input)
}

func (c GraphQLReq) DeleteDatabase(ctx context.Context, id string) (*DeleteDatabaseResponse, error) {
	return deleteDatabase(ctx, c.Client, id)
}

func (c GraphQLReq) GetCredential(ctx context.Context, id string) (*GetCredentialResponse, error) {
	return getCredential(ctx, c.Client, id)
}

func (c GraphQLReq) CreateCredential(ctx context.Context, input CreateCredentialInput) (*CreateCredentialResponse, error) {
	return createCredential(ctx, c.Client, input)
}

func (c GraphQLReq) UpdateCredential(ctx context.Context, id string, input UpdateCredentialInput) (*UpdateCredentialResponse, error) {
	return updateCredential(ctx, c.Client, id, input)
}

func (c GraphQLReq) DeleteCredential(ctx context.Context, id string) (*DeleteCredentialResponse, error) {
	return deleteCredential(ctx, c.Client, id)
}

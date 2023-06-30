// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"terraform-provider-querydesk/internal/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccDatabaseResource(t *testing.T) {
	mockClient := client.NewMockGraphQLClient(t)

	const dbId = "db_12345"

	mockClient.EXPECT().CreateDatabase(
		mock.Anything,
		client.CreateDatabaseInput{Name: "one", Adapter: "postgres", Hostname: "localhost", Database: "mydb", Ssl: false, RestrictAccess: true, Cacertfile: "", Keyfile: "", Certfile: "", AgentId: ""},
	).Return(&client.CreateDatabaseResponse{
		CreateDatabase: client.CreateDatabaseCreateDatabaseCreateDatabaseResult{
			Result: client.CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase{
				Id: dbId,
			},
			Errors: nil,
		},
	}, nil)

	mockClient.EXPECT().GetDatabase(
		mock.Anything,
		dbId,
	).Return(&client.GetDatabaseResponse{
		Database: client.GetDatabaseDatabase{
			Id:             dbId,
			Name:           "one",
			Adapter:        "postgres",
			Hostname:       "localhost",
			Database:       "mydb",
			Ssl:            false,
			RestrictAccess: true,
		},
	}, nil).Times(3)

	mockClient.EXPECT().UpdateDatabase(
		mock.Anything,
		dbId,
		client.UpdateDatabaseInput{Name: "two", Adapter: "postgres", Hostname: "localhost", Database: "mydb", Ssl: false, RestrictAccess: true, Cacertfile: "", Keyfile: "", Certfile: "", AgentId: ""},
	).Return(&client.UpdateDatabaseResponse{
		UpdateDatabase: client.UpdateDatabaseUpdateDatabaseUpdateDatabaseResult{
			Result: client.UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase{
				Id: dbId,
			},
			Errors: nil,
		},
	}, nil)

	mockClient.EXPECT().GetDatabase(
		mock.Anything,
		dbId,
	).Return(&client.GetDatabaseResponse{
		Database: client.GetDatabaseDatabase{
			Id:             dbId,
			Name:           "two",
			Adapter:        "postgres",
			Hostname:       "localhost",
			Database:       "mydb",
			Ssl:            false,
			RestrictAccess: true,
		},
	}, nil)

	mockClient.EXPECT().DeleteDatabase(
		mock.Anything,
		dbId,
	).Return(&client.DeleteDatabaseResponse{
		DeleteDatabase: client.DeleteDatabaseDeleteDatabaseDeleteDatabaseResult{
			Result: client.DeleteDatabaseDeleteDatabaseDeleteDatabaseResultResultDatabase{
				Id: dbId,
			},
			Errors: nil,
		},
	}, nil)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(mockClient),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("querydesk_database.test", "id", dbId),
					resource.TestCheckResourceAttr("querydesk_database.test", "name", "one"),
					resource.TestCheckResourceAttr("querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("querydesk_database.test", "restrict_access", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "querydesk_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("querydesk_database.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatabaseResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "querydesk_database" "test" {
  name     = %[1]q
  adapter  = "postgres"
  hostname = "localhost"
  database = "mydb"
}
`, name)
}

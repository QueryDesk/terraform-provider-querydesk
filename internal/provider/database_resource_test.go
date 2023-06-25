// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccDatabaseResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("querydesk_database.test", "name", "one"),
					resource.TestCheckResourceAttr("querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("querydesk_database.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "querydesk_database.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"name", "defaulted"},
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
	return fmt.Sprintf(`
resource "querydesk_database" "test" {
  name     = %[1]q
  adapter  = "postgres"
  hostname = "localhost"
  database = "mydb"
}
`, name)
}

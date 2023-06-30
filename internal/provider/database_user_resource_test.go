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

func TestAccDatabaseUserResource(t *testing.T) {
	mockClient := client.NewMockGraphQLClient(t)

	const dbId = "db_12345"
	const credId = "crd_12345"

	mockClient.EXPECT().CreateCredential(
		mock.Anything,
		client.CreateCredentialInput{
			DatabaseId:      dbId,
			Description:     "",
			Password:        "postgres",
			ReviewsRequired: 0,
			Username:        "postgres",
		},
	).Return(&client.CreateCredentialResponse{
		CreateCredential: client.CreateCredentialCreateCredentialCreateCredentialResult{
			Result: client.CreateCredentialCreateCredentialCreateCredentialResultResultCredential{
				Id: credId,
			},
			Errors: nil,
		},
	}, nil)

	mockClient.EXPECT().GetCredential(
		mock.Anything,
		credId,
	).Return(&client.GetCredentialResponse{
		Credential: client.GetCredentialCredential{
			Id:              credId,
			Description:     "",
			ReviewsRequired: 0,
			Username:        "postgres",
			Database: client.GetCredentialCredentialDatabase{
				Id: dbId,
			},
		},
	}, nil).Times(3)

	mockClient.EXPECT().UpdateCredential(
		mock.Anything,
		credId,
		client.UpdateCredentialInput{
			Description:     "",
			NewPassword:     "postgres",
			ReviewsRequired: 0,
			Username:        "other_user",
		},
	).Return(&client.UpdateCredentialResponse{
		UpdateCredential: client.UpdateCredentialUpdateCredentialUpdateCredentialResult{
			Result: client.UpdateCredentialUpdateCredentialUpdateCredentialResultResultCredential{
				Id: credId,
			},
			Errors: nil,
		},
	}, nil)

	mockClient.EXPECT().GetCredential(
		mock.Anything,
		credId,
	).Return(&client.GetCredentialResponse{
		Credential: client.GetCredentialCredential{
			Id:              credId,
			Description:     "",
			ReviewsRequired: 0,
			Username:        "other_user",
			Database: client.GetCredentialCredentialDatabase{
				Id: dbId,
			},
		},
	}, nil)

	mockClient.EXPECT().DeleteCredential(
		mock.Anything,
		credId,
	).Return(&client.DeleteCredentialResponse{
		DeleteCredential: client.DeleteCredentialDeleteCredentialDeleteCredentialResult{
			Result: client.DeleteCredentialDeleteCredentialDeleteCredentialResultResultCredential{
				Id: credId,
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
				Config: testAccDatabaseUserResourceConfig("postgres"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("querydesk_database_user.test", "id", credId),
					resource.TestCheckResourceAttr("querydesk_database_user.test", "username", "postgres"),
					resource.TestCheckResourceAttr("querydesk_database_user.test", "password", "postgres"),
					resource.TestCheckResourceAttr("querydesk_database_user.test", "reviews_required", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "querydesk_database_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccDatabaseUserResourceConfig("other_user"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("querydesk_database_user.test", "username", "other_user"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatabaseUserResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "querydesk_database_user" "test" {
	database_id      = "db_12345"
  username         = %[1]q
  password 		     = "postgres"
  reviews_required = 0
}
`, name)
}

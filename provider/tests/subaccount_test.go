package tests

import (
	"context"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSubAccount(t *testing.T) {
	testUser := "acc-test-tf-provider-user"
	resourceName := "chaossearch_sub_account.sub-account"
	resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSubAccountConfig(testUser),
				Check: resource.ComposeTestCheckFunc(
					testAccSubAccountExists(resourceName, testUser),
				),
			},
		},
	})
}

func testAccSubAccountConfig(username string) string {
	return fmt.Sprintf(`
		%s
	    resource "chaossearch_sub_account" "sub-account" {
		  username  = "%s"
		  full_name = "Test User"
		  password  = "%s"
		  hocon     = tolist(["override.Services.worker.quota=50"])
	    }
	`, testAccProviderConfigBlock(), username, uuid.New().String())
}

func testAccSubAccountExists(resourceName, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("Sub Account Email is not set")
		}

		providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
		csClient := providerMeta.CSClient
		ctx := context.Background()
		response, err := csClient.ListUsers(ctx, providerMeta.Token)
		if err != nil {
			return err
		}

		subaccount, err := resources.SortAndGetSubAccount(response.Users[0].SubAccounts, username)
		if err != nil {
			return err
		}

		if subaccount.Username != res.Primary.ID {
			return fmt.Errorf("SubAccount not found.")
		}

		return nil
	}
}

package tests

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserGroup(t *testing.T) {
	resourceName := "chaossearch_user_group.user_group_create"
	groupName := generateName("acc-test-tf-provider-ug")
	resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccUserGroupExists(resourceName),
				),
			},
		},
	})
}

func testAccUserGroupConfig(groupName string) string {
	return fmt.Sprintf(`
		%s
	    resource "chaossearch_user_group" "user_group_create" {
			user_groups {
			  name = "%s"
			  permissions {
				  effect    = "Allow"
				  actions    = ["*"]
				  resources = ["*"]
				  version   = "1.2"
			  }
			}
		  }
	`, testAccProviderConfigBlock(), groupName)
}

func testAccUserGroupExists(resourceName string) resource.TestCheckFunc {
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
		req := &client.ReadUserGroupRequest{
			AuthToken: providerMeta.Token,
			ID:        res.Primary.ID,
		}

		userGroup, err := csClient.ReadUserGroup(ctx, req)
		if err != nil {
			return err
		}

		if userGroup.ID != res.Primary.ID {
			return fmt.Errorf("User Group not found.")
		}

		return nil
	}
}

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

func TestAccDestinations(t *testing.T) {
	defer resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccDestDestroy,
		Steps: []resource.TestStep{
			testDestStep(
				testAccDestConf,
				"chaossearch_destination.dest",
				generateName("acc-test-tf-dest"),
			),
			testDestStep(
				testAccDestCstmConf,
				"chaossearch_destination.dest-custom",
				generateName("acc-test-tf-dest-custom"),
			),
			testDestStep(
				testAccDestHookConf,
				"chaossearch_destination.dest-hook",
				generateName("acc-test-tf-dest-hook"),
			),
		},
	})
	t.Parallel()
}

func testDestStep(config func(string) string, rsrcName, objName string) resource.TestStep {
	return resource.TestStep{
		Config: config(objName),
		Check: resource.ComposeTestCheckFunc(
			testAccDestExists(rsrcName, objName),
			resource.TestCheckResourceAttr(rsrcName, "name", objName),
		),
	}
}

func testAccDestExists(rsrcName, objName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[rsrcName]
		if !ok {
			return fmt.Errorf("Not found: %s", rsrcName)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("Destination ID is not set")
		}

		providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
		csClient := providerMeta.CSClient
		ctx := context.Background()
		resp, err := csClient.ReadDestination(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        objName,
		})

		if err != nil {
			return err
		}

		found := false
		for _, dest := range resp.Destinations {
			if dest.Id == res.Primary.ID {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("Destination not found.")
		}

		return nil
	}
}

func testAccDestDestroy(s *terraform.State) error {
	providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
	csClient := providerMeta.CSClient
	ctx := context.Background()

	for _, res := range s.RootModule().Resources {
		resp, err := csClient.ReadDestination(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        res.Primary.ID,
		})

		if err != nil {
			return err
		}

		found := false
		for _, dest := range resp.Destinations {
			if dest.Id == res.Primary.ID {
				found = true
			}
		}

		if found {
			return fmt.Errorf("Destination (%s) still exists", res.Primary.ID)
		}
	}

	return nil
}

func testAccDestConf(name string) string {
	return fmt.Sprintf(`
	%s
	resource "chaossearch_destination" "dest" {
		name = "%s"
		type = "slack"
		slack {
		  url = "http://slack.com"
		}
	  }
	`, testAccProviderConfigBlock(), name)
}

func testAccDestCstmConf(name string) string {
	return fmt.Sprintf(`
	%s
	resource "chaossearch_destination" "dest-custom" {
		name = "%s"
		type = "custom_webhook"
		custom_webhook {
		  url = "http://test.com"
		}
	  }
	`, testAccProviderConfigBlock(), name)
}

func testAccDestHookConf(name string) string {
	return fmt.Sprintf(`
	%s
	resource "chaossearch_destination" "dest-hook" {
		name = "%s"
		type = "custom_webhook"
		custom_webhook {
		  scheme = "HTTPS"
		  host = "test.com"
		  path = "/api/test"
		  port = "8080"
		  method = "POST"
		  query_params = jsonencode({
			"test": "value"
		  })
		  header_params = jsonencode({
			"Content-Type": "application/json"
		  })
		}
	  }
	`, testAccProviderConfigBlock(), name)
}

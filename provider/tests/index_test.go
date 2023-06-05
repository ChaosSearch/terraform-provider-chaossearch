package tests

import (
	"context"
	"cs-tf-provider/provider/models"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIndex(t *testing.T) {
	defer resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			testIndexStep(
				testAccIndexConfig,
				"chaossearch_index_model.model",
				generateName("acc-test-tf-provider-index"),
			),
		},
	})
	t.Parallel()
}

func testIndexStep(config func(string) string, rsrcName, objName string) resource.TestStep {
	return resource.TestStep{
		Config: testAccIndexConfig(objName),
		Check: resource.ComposeTestCheckFunc(
			testAccIndexExists(rsrcName, objName),
			resource.TestCheckResourceAttr(rsrcName, "bucket_name", objName),
		),
	}
}

func testAccIndexConfig(bucket string) string {
	return fmt.Sprintf(`
		%s
		resource "chaossearch_index_model" "model" {
			bucket_name = "%s"
			model_mode = 0
			options {
				delete_enabled = true
			}
			depends_on = [
			  chaossearch_object_group.csv-og
			]
		  }
	`, testAccObjectGroupConfigCSV(bucket), bucket)
}

func testAccIndexExists(resourceName, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("Object Group ID is not set")
		}

		providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
		csClient := providerMeta.CSClient
		ctx := context.Background()

		response, err := csClient.ReadIndexModel(
			ctx,
			bucketName,
			providerMeta.Token,
		)

		if err != nil {
			return err
		}

		if response.Contents == nil {
			return fmt.Errorf("Index Content not found.")
		}

		return nil
	}
}

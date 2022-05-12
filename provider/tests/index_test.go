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
	resourceName := "chaossearch_object_group.create-object-group"
	bucketName := "acc-test-tf-provider-index"
	resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIndexConfig(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccIndexExists(resourceName, bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName),
				),
			},
		},
	})
}

func testAccIndexConfig(bucket string) string {
	return fmt.Sprintf(`
		%s
		resource "chaossearch_index_model" "model" {
			bucket_name = "%s"
			model_mode = 0
			delete_enabled = true
			depends_on = [
			  chaossearch_object_group.create-object-group
			]
		  }
	`, testAccObjectGroupConfig(bucket), bucket)
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

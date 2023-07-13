package tests

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const csvSource string = "chaossearch-tf-provider-test"
const jsonSource string = "chaossearch-self-service-users"

func TestAccObjectGroup(t *testing.T) {
	defer resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccObjectGroupDestroy,
		Steps: []resource.TestStep{
			testOGStep(
				testAccObjectGroupConfigCSV,
				"chaossearch_object_group.csv-og",
				generateName("acc-test-tf-csv-og"),
			),
			testOGStep(
				testAccObjectGroupConfigJSON,
				"chaossearch_object_group.json-og",
				generateName("acc-test-tf-json-og"),
			),
		},
	})
	t.Parallel()
}

func testOGStep(config func(string) string, rsrcName, objName string) resource.TestStep {
	return resource.TestStep{
		Config: config(objName),
		Check: resource.ComposeTestCheckFunc(
			testAccObjectGroupExists(rsrcName, objName),
			resource.TestCheckResourceAttr(rsrcName, "bucket", objName),
		),
	}
}

func testAccObjectGroupConfigCSV(bucket string) string {
	return fmt.Sprintf(`
		%s
		resource "chaossearch_object_group" "csv-og" {
		  bucket = "%s"
		  source = "%s"
		  format {
			type             = "CSV"
			column_delimiter = ","
			row_delimiter    = "\\n"
			header_row       = true
		  }
		  index_retention {
			overall = -1
		  }
		  options {
			col_types = jsonencode({
			  "Period": "Timeval"
			})
		  }
		  filter {
			field  = "key"
			prefix = "ec"
		  }
		  filter {
			field = "key"
			regex = ".*"
		  }
		}
	`, testAccProviderConfigBlock(), bucket, csvSource)
}

func testAccObjectGroupConfigJSON(bucket string) string {
	return fmt.Sprintf(`
		%s
		resource "chaossearch_object_group" "json-og" {
		  bucket = "%s"
		  source = "%s"
		  format {
			type                = "JSON"
			array_flatten_depth = -1
			field_selection     = jsonencode([{
				"type": "blacklist",
				"excludes": [
					"email",
					"hq_phone"
				]
			}])
		  }
		  index_retention {
			overall = -1
		  }
		  options {
			col_types = jsonencode({
			  "createddate.value": "Timeval"
			})
		  }
		  filter {
			field = "key"
			regex = ".*"
		  }
		}
		
	`, testAccProviderConfigBlock(), bucket, jsonSource)
}

func testAccObjectGroupExists(resourceName, bucketName string) resource.TestCheckFunc {
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
		response, err := csClient.ReadObjGroup(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        bucketName,
		})

		if err != nil {
			return err
		}

		if response.ID != res.Primary.ID {
			return fmt.Errorf("Object Group not found.")
		}

		return nil
	}
}

func testAccObjectGroupDestroy(s *terraform.State) error {
	providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
	csClient := providerMeta.CSClient
	ctx := context.Background()

	for _, res := range s.RootModule().Resources {
		response, err := csClient.ReadObjGroup(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        res.Primary.ID,
		})

		if err == nil {
			if response != nil && response.Bucket == res.Primary.ID {
				return fmt.Errorf("Object Group (%s) still exists.", res.Primary.ID)
			}
		}

		if !strings.Contains(err.Error(), "NotFound") {
			return err
		}
	}

	return nil
}

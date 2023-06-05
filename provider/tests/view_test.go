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

func TestAccView(t *testing.T) {
	bucketName := generateName("acc-test-tf-provider-view-og")
	defer resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			testViewStep(
				testAccViewConfig,
				"chaossearch_view.view",
				bucketName,
				generateName("acc-test-tf-provider-view"),
			),
			testViewStep(
				testAccViewPredsConfig,
				"chaossearch_view.view-preds",
				bucketName,
				generateName("acc-test-tf-provider-view-preds"),
			),
			testViewStep(
				testAccViewTransformsConfig,
				"chaossearch_view.view-transforms",
				bucketName,
				generateName("acc-test-tf-provider-view-transforms"),
			),
		},
	})
	t.Parallel()
}

func testViewStep(config func(string, string) string, rsrcName, srcName, objName string) resource.TestStep {
	return resource.TestStep{
		Config: config(objName, srcName),
		Check: resource.ComposeTestCheckFunc(
			testAccViewExists(rsrcName, srcName),
		),
	}
}

func testAccViewConfig(viewName, bucketName string) string {
	return fmt.Sprintf(`
		%s
	    resource "chaossearch_view" "view" {
			bucket           = "%s"
			case_insensitive = false
			index_pattern    = ".*"
			index_retention  = -1
			overwrite        = true
			sources          = ["%s"]
			time_field_name  = "timestamp"
			filter {
			  predicate {
				type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
				pred {
				  type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
				  field = "cs_partition_key_0"
				  query = "Test"
				  state {
					type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
				  }
				}
			  }
			}
			depends_on = [
				chaossearch_index_model.model
			]
		  }
	`, testAccIndexConfig(bucketName), viewName, bucketName)
}

func testAccViewPredsConfig(viewName, bucketName string) string {
	return fmt.Sprintf(`
	%s
	resource "chaossearch_view" "view-preds" {
		bucket           = "%s"
		case_insensitive = false
		index_pattern    = ".*"
		index_retention  = -1
		overwrite        = true
		sources          = ["%s"]
		time_field_name  = "@timestamp"
		filter {
		  predicate {
			type  = "chaossumo.query.NIRFrontend.Request.Predicate.Or"
			preds = [
				jsonencode(
				  {
					"_type" = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch",
					"field" = "Subject",
					"query" = "Test"
					"state" = {
					  "_type" = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
					},
				  }
				),
				jsonencode(
				  {
					"_type" = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch",
					"field" = "Series_title_1",
					"query" = "Test2"
					"state" = {
					  "_type" = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
					},
				  }
				)
			  ]
		  }
		}
		depends_on = [
			chaossearch_index_model.model
		]
	  }
`, testAccIndexConfig(bucketName), viewName, bucketName)
}

func testAccViewTransformsConfig(viewName, bucketName string) string {
	return fmt.Sprintf(`
	%s
	resource "chaossearch_view" "view-transforms" {
		bucket           = "%s"
		case_insensitive = false
		index_pattern    = ".*"
		index_retention  = -1
		overwrite        = true
		sources          = ["%s"]
		time_field_name  = "Period"
		filter {
		  predicate {
			type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
			pred {
			  type  = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
			  field = "STATUS"
			  query = "*F*"
			  state {
				type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
			  }
			}
		  }
		}
		transforms = [
		  jsonencode({
			"_type": "MaterializeRegexTransform",
			"inputField": "Data_value",
			"pattern": "(\\d+)\\.(\\d+)"
			"outputFields": [
			  {
				"name": "Whole",
				"type": "NUMBER"
			  },
			  {
				"name": "Decimal",
				"type": "NUMBER"
			  }
			]
		  }),
		]
		depends_on = [
		  chaossearch_index_model.model
		]
	  }
	`, testAccIndexConfig(bucketName), viewName, bucketName)
}

func testAccViewExists(resourceName, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("View ID is not set")
		}

		providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
		csClient := providerMeta.CSClient
		ctx := context.Background()
		response, err := csClient.ReadView(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        bucketName,
		})

		if err != nil {
			return err
		}

		if response.ID == "" && response.Bucket != bucketName {
			return fmt.Errorf("View not found")
		}

		return nil
	}
}

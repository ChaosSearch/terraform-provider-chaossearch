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

func TestAccMonitors(t *testing.T) {
	defer resource.Test(t, resource.TestCase{
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccMonitorDestroy,
		Steps: []resource.TestStep{
			testMonitorStep(
				testAccMonitorConf,
				"chaossearch_monitor.monitor",
				generateName("acc-test-tf-monitor"),
			),
		},
	})
	t.Parallel()
}

func testMonitorStep(config func(string) string, rsrcName, objName string) resource.TestStep {
	return resource.TestStep{
		Config: config(objName),
		Check: resource.ComposeTestCheckFunc(
			testAccMonitorExists(rsrcName),
			resource.TestCheckResourceAttr(rsrcName, "name", objName),
		),
	}
}

func testAccMonitorExists(rsrcName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[rsrcName]
		if !ok {
			return fmt.Errorf("Not found: %s", rsrcName)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("Monitor ID is not set")
		}

		providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
		csClient := providerMeta.CSClient
		ctx := context.Background()
		resp, err := csClient.ReadMonitor(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        res.Primary.ID,
		})

		if err != nil {
			return err
		}

		if !resp.Ok {
			return fmt.Errorf("Monitor not found during exists check")
		}

		if len(resp.Resp.Triggers) == 0 {
			return fmt.Errorf("Triggers found empty during exists check")
		}

		for _, trigger := range resp.Resp.Triggers {
			if len(trigger.Actions) == 0 {
				return fmt.Errorf("Actions found empty for trigger '%s'", trigger.Name)
			}
		}

		return nil
	}
}

func testAccMonitorDestroy(s *terraform.State) error {
	providerMeta := testAccProvider.Meta().(*models.ProviderMeta)
	csClient := providerMeta.CSClient
	ctx := context.Background()

	for _, res := range s.RootModule().Resources {
		resp, err := csClient.ReadMonitor(ctx, &client.BasicRequest{
			AuthToken: providerMeta.Token,
			Id:        res.Primary.ID,
		})

		if err != nil {
			return err
		}

		if resp.Ok {
			return fmt.Errorf("Monitor found during destroy check")
		}
	}

	return nil
}

func testAccMonitorConf(name string) string {
	bucketName := generateName("acc-test-tf-provider-view-og")
	viewName := generateName("acc-test-tf-provider-view")
	return fmt.Sprintf(`
	%s
	resource "chaossearch_destination" "dest" {
		name = "test-dest"
		type = "slack"
		slack {
		  url = "http://slack.com"
		}
	  }

	
	resource "chaossearch_monitor" "monitor" {
		name = "%s"
		type = "monitor"
		enabled = true
		depends_on = [
			chaossearch_destination.dest,
			chaossearch_view.view,
		]
		schedule {
			period {
				interval = 1
				unit = "MINUTES"
			}
		}
		inputs {
			search {
				indices = [
					chaossearch_view.view.bucket,
				]
				query = jsonencode({
					"size": 0,
					"aggregations": {
						"when": {
							"avg": {
								"field": "Magnitude"
							},
							"meta": null
						}
					},
					"query": {
						"bool": {
							"filter": [
								{
									"range": {
										"Period": {
											"gte": "{{period_end}}||-1h",
											"lte": "{{period_end}}",
											"format": "epoch_millis"
										}
						  			}
					  			}
				  			]
			  			}
					}
		  		})
			}
	  	}
	  	triggers {
			name = "test-trigger"
			severity = "1"
			condition {
		  		script {
					lang = "painless"
					source = "ctx.results[0].hits.total.value > 1000"
		  		}
			}
			actions {
		  		name = "test-action"
		  		destination_id = chaossearch_destination.dest.id
		  		subject_template {
					lang = "mustache"
					source = "Monitor {{ctx.monitor.name}} Triggered"
		  		}
		  		message_template {
					lang = "mustache"
					source = "Some message"
		  		}
		  		throttle_enabled = true
		  		throttle {
					value = 10
					unit = "MINUTES"
		  		}
			}
	  	}
	}
	`, testAccViewConfig(viewName, bucketName), name)
}

# Monitor Resource

> Creating a monitor in ChaosSearch allows you to specify a particular condition or event for which you want to be alerted. You can define a monitor to watch by using an extraction query or a visual graph.

Check out the _Alerting Overview_ documentation here: [Alerting Overview](https://docs.chaossearch.io/docs/alerting-overview)

## Example Usage
```hcl
resource "chaossearch_monitor" "monitor" {
  name = "tf-provider-monitor"
  type = "monitor"
  enabled = true
  schedule {
    period {
      interval = 1
      unit = "MINUTES"
    }
  }
  inputs {
    search {
      indices = ["example-view-bucket-name"]
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
    name = "tf-provider-trigger"
    severity = "1"
    condition {
      script {
        lang = "painless"
        source = "ctx.results[0].hits.total.value > 1000"
      }
    }
    actions {
      name = "tf-provider-action"
      destination_id = "example-destination-id"
      subject_template {
        lang = "mustache"
        source = "Monitor {{ctx.monitor.name}} Triggered"
      }
      message_template {
        lang = "mustache"
        source = "Monitor {{ctx.monitor.name}} just entered alert status. Please investigate the issue.\n- Trigger: {{ctx.trigger.name}}\n- Severity: {{ctx.trigger.severity}}\n- Period start: {{ctx.periodStart}}\n- Period end: {{ctx.periodEnd}}"
      }
      throttle_enabled = true
      throttle {
        value = 10
        unit = "MIN"
      }
    }
  }
}
```

**Note** Multiple trigger blocks and action blocks can be defined in one resource
**Note** All fields are required except for `throttle_enabled` and `throttle`
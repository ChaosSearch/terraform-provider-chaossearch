terraform {
  required_providers {
    chaossearch = {
      version = "~> 1.0.10"
      source  = "chaossearch/chaossearch"
    }
  }
}

provider "chaossearch" {
  login {}
}

resource "chaossearch_user_group" "user_group" {
  name = "provider_test"
  permissions = jsonencode([
    {
      "Version"   = "1.0",
      "Effect"    = "Allow",
      "Actions"   = ["ui:analytics"],
      "Resources" = ["*"],
    },
    {
      "Version"   = "1.0",
      "Effect"    = "Allow",
      "Actions"   = ["ui:storage"]
      "Resources" = ["*"],
      "Condition" = {
        "Conditions" = [
          {
            "Equals"     = {
              "chaos:document/attributes.title" = ""
            },
            "Like"       = {
              "chaos:document/attributes.title" = ""
            },
            "NotEquals"  = {
              "chaos:document/attributes.title" = ""
            },
            "StartsWith" = {
              "chaos:document/attributes.title" = "test"
            },
          }
        ]
      }
    }
  ])
}

resource "chaossearch_sub_account" "sub-account" {
  username  = "provider_test"
  full_name = "provider_test"
  password  = "1234"
  group_ids = [chaossearch_user_group.user_group.id]
  hocon     = ["override.Services.worker.quota=50"]
}

resource "chaossearch_object_group" "create-object-group" {
  bucket = "tf-provider"
  source = "chaossearch-tf-provider-test"
  //live_events = "test"
  format {
    type             = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = true
  }
  index_retention {
    overall       = -1
  }

  options {
    #compression = "GZIP"
    col_types = jsonencode({
      "Period": "Timeval"
    })
  }
  // Filter options:
  filter {
    field = "key"
    prefix = "ec"
  }
  filter {
    field = "key"
    regex = ".*"
  }
  /*
  filter {
    field = "storageClass"
    equals = "STANDARD"
  }
  */
}

resource "chaossearch_object_group" "selection-og" {
  bucket = "tf-provider-selections"
  source = "chaossearch-tf-provider-test"
  format {
    type             = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = true
    field_selection = jsonencode([{
        "excludes": [
          "data",
          "bigobject"
        ],
        "type": "blacklist"
    }])
    array_selection = jsonencode([{
      "includes": [
        "object.ids",
      ],
      "type": "whitelist"
    }])
    vertical_selection = jsonencode([{
      "include": true,
       "patterns": [
        "^line\\.level$",
        "^attrs.version$",
        "^timestamp$",
        "^line\\.meta\\.[^\\.]*$",
        "^host$",
        "^line\\.correlation_id$",
        "^sourcetype$",
        "^line\\.message$",
        "^message$",
        "^source$",
        "^_rawJson$"
      ],
      "type": "regex"
    }])
  }
  index_retention {
    overall       = -1
  }

  options {
    col_types = jsonencode({
      "TimeStamp": "Timeval"
    })
    col_renames = jsonencode({
      "TimeStamp": "Period"
    })
    col_selection = jsonencode([{
      "includes": [
        "object.ids",
      ],
      "type": "whitelist"
    }])
  }
}

resource "chaossearch_index_model" "model" {
  bucket_name = "tf-provider"
  model_mode  = 0
  delete_enabled = true
  depends_on  = [
    chaossearch_object_group.create-object-group
  ]
}

resource "chaossearch_view" "view-pred" {
  bucket           = "tf-provider-view"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
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
  depends_on = [
    chaossearch_index_model.model
  ]
}

resource "chaossearch_view" "view-preds" {
  bucket           = "tf-provider-view-preds"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
  time_field_name  = "Period"
  filter {
    predicate {
      type = "chaossumo.query.NIRFrontend.Request.Predicate.Or"
      preds = [
        jsonencode(
          {
            "state": {
              "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
            },
            "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch",
            "field": "Subject",
            "query": "subject"
          }
        ),
        jsonencode(
          {
            "state": {
              "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
            },
            "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch",
            "field": "Series_title_1",
            "query": "title"
          }
        )
      ]
    }
  }
  depends_on = [
    chaossearch_index_model.model
  ]
}

resource "chaossearch_view" "view-transforms" {
  bucket           = "tf-provider-view-transforms"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
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

resource "chaossearch_destination" "dest" {
  name = "tf-provider-destination"
  type = "slack"
  slack {
    url = "http://slack.com"
  }
}

resource "chaossearch_destination" "dest_custom" {
  name = "tf-provider-destination-custom"
  type = "custom_webhook"
  custom_webhook {
    url = "http://test.com"
  }
}

resource "chaossearch_destination" "dest_custom_host" {
  name = "tf-provider-destination-custom-host"
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

resource "chaossearch_monitor" "monitor" {
  name = "tf-provider-monitor"
  type = "monitor"
  enabled = true
  depends_on = [
    chaossearch_destination.dest,
    chaossearch_view.view-pred,
    chaossearch_view.view-preds
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
        chaossearch_view.view-pred.bucket,
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
  triggers { // Can have multiple triggers
    name = "tf-provider-trigger"
    severity = "1"
    condition {
      script {
        lang = "painless"
        source = "ctx.results[0].hits.total.value > 1000"
      }
    }
    actions { // Can have multiple actions
      name = "tf-provider-action"
      destination_id = chaossearch_destination.dest.id
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

data "chaossearch_retrieve_sub_accounts" "sub_accounts" {}

output "object_group_retrieve_sub_accounts" {
  value = data.chaossearch_retrieve_sub_accounts.sub_accounts
}

data "chaossearch_retrieve_object_group" "object-group" {
  bucket     = "tf-provider"
  depends_on = [
    chaossearch_object_group.create-object-group
  ]
}

output "object_group" {
  value = data.chaossearch_retrieve_object_group.object-group
}

data "chaossearch_retrieve_object_groups" "object-groups" {}

output "object-groups" {
  value = data.chaossearch_retrieve_object_groups.object-groups
}

data "chaossearch_retrieve_view" "retrieve_view" {
  bucket     = "tf-provider-view"
  depends_on = [
    chaossearch_view.view-pred
  ]
}

output "view" {
  value = data.chaossearch_retrieve_view.retrieve_view
}

data "chaossearch_retrieve_views" "views" {}

output "views" {
  value = data.chaossearch_retrieve_views.views
}

data "chaossearch_retrieve_groups" "user_groups"{}

output "user_groups" {
  value = data.chaossearch_retrieve_groups.user_groups
}

data "chaossearch_retrieve_user_group" "user_group"{
  id         = chaossearch_user_group.user_group.id
  depends_on = [
    chaossearch_user_group.user_group
  ]
}
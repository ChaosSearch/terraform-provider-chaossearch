terraform {
  required_providers {
    chaossearch = {
      version = "0.1.0"
      source = "chaos/chaossearch"
    }
  }
}

provider "chaossearch" {
    url               = "http://example.com"
    access_key_id     = "chaossearch_access_key_id"
    secret_access_key = "chaossearch_secret_access_key"
    region            = "eu-west-1"
}

resource "chaossearch_object_group" "my-object-group" {
  name = "my-object-group"
  source_bucket = "<s3 bucket name>"
  live_events_sqs_arn ="arn:aws:sqs:sqs_sqs"

  filter_json = jsonencode({
    AND = [
      {
        field = "key"
        regex = ".*"
      }
    ]
  })

  compression = "gzip"
  format = "JSON"

  partition_by = "<regex>"
  array_flatten_depth = -1

  keep_original = true

  column_selection {
    type = "whitelist"
    includes = [
      "host",
      "source",
    ]
  }
}

resource "chaossearch_indexing_state" "my-object-group" {
  object_group_name = chaossearch_object_group.my-object-group.name
  active = true
}

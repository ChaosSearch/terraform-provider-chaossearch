terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.2"
      source = "chaos/chaossearch"
    }
  }
}
provider "chaossearch" {
    url               = "https://ap-south-1-aeternum.chaossearch.io"
    access_key_id     = "LCE8T6HRFGJI3ZKBGMGD"
    secret_access_key = "r5MEYkYntYvXqRSBMK6SFLQfPw7hHRQ0v5cqlkIk"
    region            = "ap-south-1"
    login  {
      user_name = "service_user@chaossearch.com"
      password = "thisIsAnEx@mple1!"
      parent_user_id = "be4aeb53-21d5-4902-862c-9c9a17ad6675"
    }

}


resource "chaossearch_view" "chaossearch-create-view" {
  bucket="dinesh-tf-004"
  index_pattern=".*"
  filter_json=""
  //array_flatten_depth =-1 
  case_insensitive=false
  index_retention =-1
  transforms=[]
  sources=["test-object-group"]
}


#  resource "chaossearch_object_group" "my-object-group" {
#    name = "dines-object-group-003"
#    source_bucket = "chaos-test-data-aps1"
#    live_events_sqs_arn ="arn:aws:sqs:sqs_sqs"

#    filter_json = jsonencode({
#      AND = [
#        {
#          field = "key"
#          regex = ".*"
#        }
#      ]
#    })

#    compression = "gzip"
#    format = "JSON"

#    partition_by = ""
#    array_flatten_depth = -1

#    keep_original = true

#    column_selection {
#      type = "whitelist"
#      includes = [
#        "host",
#        "source",
#      ]
#    }
#  }

# resource "chaossearch_indexing_state" "my-object-group" {
#   object_group_name = chaossearch_object_group.my-object-group.name
#   active = true
# }


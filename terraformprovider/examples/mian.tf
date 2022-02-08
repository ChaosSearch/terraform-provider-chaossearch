terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.19"
      source = "chaos/chaossearch"
    }
  }
}

provider "chaossearch" {
    url               = "https://ap-south-1-aeternum.chaossearch.io"
    access_key_id     = "LCE8T6HRFGJI3ZKBGMGD"
    secret_access_key = "r5MEYkYntYvXqRSBMK6SFLQfPw7hHRQ0v5cqlkIk"
    region            = "ap-south-1"

}


resource "chaossearch_view" "chaossearch-create-view" {
  name="dinesh-tf-pro-test-9909"
  bucket="my-bucket-123"
  index_pattern=".*"
 filter_json=""
 //array_flatten_depth =-1 
  case_insensitive=false
index_retention=-1
transforms=[""]
sources=[]
}

# resource "chaossearch_object_group" "my-object-group" {
#   name = "my-object-group-dinesh-1"
#   source_bucket = "<s3 bucket name>"
#   live_events_sqs_arn ="arn:aws:sqs:sqs_sqs"

#   filter_json = jsonencode({
#     AND = [
#       {
#         field = "key"
#         regex = ".*"
#       }
#     ]
#   })

#   compression = "gzip"
#   format = "JSON"

#   partition_by = "<regex>"
#   array_flatten_depth = -1

#   keep_original = true

#   column_selection {
#     type = "whitelist"
#     includes = [
#       "host",
#       "source",
#     ]
#   }
# }

# resource "chaossearch_indexing_state" "my-object-group" {
#   object_group_name = chaossearch_object_group.my-object-group.name
#   active = true
# }


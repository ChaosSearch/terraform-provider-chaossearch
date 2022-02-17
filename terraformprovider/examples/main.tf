terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source = "chaossearch/chaossearch"
    }
  }
}
provider "chaossearch" {
  url               = "https://ap-south-1-aeternum.chaossearch.io"
  access_key_id     = "LCE8T6HRFGJI3ZKBGMGD"
  secret_access_key = "r5MEYkYntYvXqRSBMK6SFLQfPw7hHRQ0v5cqlkIk"
  region            = "ap-south-1"
  login {
    user_name      = "service_user@chaossearch.com"
    password       = "thisIsAnEx@mple1!"
    parent_user_id = "be4aeb53-21d5-4902-862c-9c9a17ad6675"
  }

}


# resource "chaossearch_view" "chaossearch-create-view" {
#   bucket="nibras-tf-005"
#   index_pattern=".*"
#   filter_json=""
#   //array_flatten_depth =-1
#   case_insensitive=false
#   index_retention =-1
#   transforms=[]
#   sources=[]
# }


resource "chaossearch_object_group" "my-object-group" {

  bucket = "nibras-og-0103"
  source = "chaos-test-data-aps1"
  format {
    _type            = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = true
  }
  interval {
    mode   = 0
    column = 0
  }
  index_retention {
    for_partition = []
    overall       = -1
  }
  filter {
    obj1 {
      field  = "key"
      prefix = "bluebike"
    }
    obj2 {
      field = "key"
      regex = ".*"
    }
  }
  options {
    ignore_irregular = true
  }
  realtime = false
}

# resource "chaossearch_indexing_state" "my-object-group" {
#   object_group_name = chaossearch_object_group.my-object-group.name
#   active = true
# }


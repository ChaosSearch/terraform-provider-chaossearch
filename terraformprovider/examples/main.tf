terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source  = "chaossearch/chaossearch"
    }
  }
}

provider "chaossearch" {
  url               = var.url
  access_key_id     = var.access_key_id
  secret_access_key = var.secret_access_key
  region            = var.region
  login {

    # Normal User Credentials
#        user_name      = var.user_name
#        password       = var.password
#        parent_user_id = var.parent_user_id

    # ADMIN Credentials
    user_name = var.admin_user_name
    password  = var.admin_password
  }
}
resource "chaossearch_object_group" "my-object-group-test-3" {
  bucket = "test-object-group-005"
  source = "chaos-test-data-aps1"
  format {
    _type            = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = false
  }
  interval {
    mode   = 0
    column = 0
  }
  index_retention {
    for_partition = []
    overall       = 1
  }
  filter {
    prefix_filter {
      field  = "key"
      prefix = "bluebike"
    }
    regex_filter {
      field = "key"
      regex = ".*"
    }
  }
  options {
    ignore_irregular = true
  }
  realtime = true
}

#data "chaossearch_retrieve_user_group" "my-user-group" {
#  user_groups {
#    id = "98928445-865c-4606-a1cc-c2395a2fad13"
#  }
#}
#get user group by id
#data "chaossearch_retrieve_user_group" "my-user-group" {
#  id="38afab15-76e9-40ee-bdff-18fcd5480437"
#}
#
#output "object_group_retrieve_user_group" {
#  value = data.chaossearch_retrieve_user_group.my-user-group
#}




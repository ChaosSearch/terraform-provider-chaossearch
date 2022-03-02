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

resource "chaossearch_user_group" "chaossearch_user_group_crate_test" {
  user_groups {
    id   = "1002215"
    name = "nibras-ug-0031"
    permissions {
      permission {

        effect    = "Allow"
        action    = "2.*"
        resources = "*"
        version   = "1.3"
      }
    }
  }
}

#data "chaossearch_retrieve_user_group" "my-user-group" {
#  user_groups {
#    id = "9436aed9-e994-4dba-a25b-7d950d7f3623"
#  }
#}
##get user group by id
#data "chaossearch_retrieve_user_group" "my-user-group" {
#  id="9436aed9-e994-4dba-a25b-7d950d7f3623"
#}
##
#output "object_group_retrieve_user_group" {
#  value = data.chaossearch_retrieve_user_group.my-user-group
#}




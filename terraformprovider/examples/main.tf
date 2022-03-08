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

resource "chaossearch_view" "chaossearch-update-view-test" {
  bucket           = "test-view-check-2"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = []
  time_field_name  = "@timestamp"
  transforms       = []
  filter {
    predicate {
      _type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
      pred {
        _type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
        field = "cs_partition_key_0"
        query = "*bluebike*"
        state {
          _type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
        }
      }
    }
  }
}



#resource "chaossearch_user_group" "chaossearch_user_group_crate_test" {
#  user_groups {
#    id   = "46b477fa-1bf8-49c8-8130-3ebe9e694421"
#    name = "chathura-delete-0005"
#    permissions {
#      effect    = "Allow"
#      actions   = ["*1"]
#      resources = ["*1"]
#      version   = "11"
#
#    }
#  }
#}

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


#get view group by id
#data "chaossearch_retrieve_object_group" "object-group" {
#  bucket = "Chathura-og-jhjhjhj"
#}
#
#output "object_group" {
#  value = data.chaossearch_retrieve_object_group.object-group
#}


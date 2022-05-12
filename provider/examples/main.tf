terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source  = "chaossearch/chaossearch"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "chaossearch" {
  url               = var.CS_URL
  access_key_id     = var.CS_ACCESS_KEY
  secret_access_key = var.CS_SECRET_KEY
  region            = var.CS_REGION
  login {
    user_name = var.CS_USERNAME
    password  = var.CS_PASSWORD
  }
}

provider "aws" {
  profile = var.AWS_PROFILE
  region  = var.CS_REGION
}


resource "chaossearch_user_group" "chaossearch_user_group_create_test" {
  user_groups {
    name = "provider_test"
    permissions {
        effect    = "Allow"
        actions    = [".*"]
        resources = [".*"]
        version   = "1.4"
      }
  }
}
resource "chaossearch_sub_account" "sub-account" {
  user_info_block {
    username  = "provider_test"
    full_name = "provider_test"
    email     = "provider_test"
  }
  password  = "1234"
  hocon     = ["override.Services.worker.quota=50"]
}

resource "chaossearch_object_group" "create-object-group" {
  bucket = "tf-provider"
  source = "chaossearch-tf-provider-test"
  format {
    type             = "CSV"
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
    regex_filter {
      field = "key"
      regex = ".*"
    }
  }
  options {
    ignore_irregular = true
  }
}

resource "chaossearch_index_model" "model" {
  bucket_name = "tf-provider"
  model_mode = 0
  depends_on = [
    chaossearch_object_group.create-object-group
  ]
}

resource "chaossearch_view" "chaossearch-create-view" {
  bucket           = "tf-provider-view"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
  time_field_name  = "@timestamp"
  filter {
    predicate {
      type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
      pred {
        type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
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

data "chaossearch_retrieve_sub_accounts" "sub_accounts" {}

output "object_group_retrieve_sub_accounts" {
  value = data.chaossearch_retrieve_sub_accounts.sub_accounts
}

data "chaossearch_retrieve_object_group" "object-group" {
  bucket = "tf-provider"
  depends_on = [
    chaossearch_object_group.create-object-group
  ]
}

output "object_group" {
  value = data.chaossearch_retrieve_object_group.object-group
}

data "chaossearch_retrieve_view" "retrieve_view" {
  bucket = "tf-provider-view"
  depends_on = [
    chaossearch_view.chaossearch-create-view
  ]
}

output "view" {
  value = data.chaossearch_retrieve_view.retrieve_view
}

data "chaossearch_retrieve_views" "views" {}

output "views" {
  value = data.chaossearch_retrieve_views.views
}
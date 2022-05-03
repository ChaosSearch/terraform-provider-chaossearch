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
    type            = "CSV"
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
    overall       = 14
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
  realtime = false
}

resource "chaossearch_index_model" "model-1" {
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
    chaossearch_index_model.model-1
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

/*
resource "aws_s3_bucket" "bucket-creation" {
  bucket = "chaossearch-tf-provider-test"
}

resource "aws_s3_bucket_object" "upload-file" {
  bucket = aws_s3_bucket.bucket-creation.id
  key    = "economic-survey-of-manufacturing-dec-2021.csv"
  source = "economic-survey-of-manufacturing-dec-2021.csv"
  etag = filemd5("economic-survey-of-manufacturing-dec-2021.csv")
}

resource "chaossearch_object_group" "create-object-group" {
  bucket = "test-object-group-tera5"
  source = "chaossearch-tf-provider-test"
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
  depends_on = [
    aws_s3_bucket_object.upload-file
  ]
}

resource "chaossearch_index_model" "model-1" {
  bucket_name = "test-object-group-tera5"
  model_mode = -1
}

data "chaossearch_retrieve_object_group" "object-group" {
  bucket = "test-object-group-tera5"
  depends_on = [
    chaossearch_object_group.create-object-group
  ]
}

output "object_group" {
  value = data.chaossearch_retrieve_object_group.object-group
}

data "chaossearch_retrieve_object_groups" "object_groups" {}

output "object_groups" {
  value = data.chaossearch_retrieve_object_groups.object_groups
}

resource "chaossearch_view" "chaossearch-create-view" {
  bucket           = "test-view-tera5"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["test-object-group-tera5"]
  time_field_name  = "@timestamp"
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
  depends_on = [
    chaossearch_index_model.model-1
  ]
}

data "chaossearch_retrieve_view" "retrieve_view" {
  bucket = "test-view-tera5"
  depends_on = [
    chaossearch_view.chaossearch-create-view
  ]
}

output "view" {
  value = data.chaossearch_retrieve_view.retrieve_view
}

resource "chaossearch_view" "chaossearch-update-view-test" {
  bucket           = "test-view-tera5"
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
        field = "cs_partition_key_0112"
        query = "*bluebike*"
        state {
          _type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
        }
      }
    }
  }
  depends_on = [
    chaossearch_view.chaossearch-create-view
  ]
}

data "chaossearch_retrieve_views" "views" {}

output "views" {
  value = data.chaossearch_retrieve_views.views
}

resource "chaossearch_index_metadata" "chaossearch-index-metadata" {
  bucket_names = "test-object-group-tera4"
}
*/

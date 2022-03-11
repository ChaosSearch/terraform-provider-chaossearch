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
  url               = var.url
  access_key_id     = var.access_key_id
  secret_access_key = var.secret_access_key
  region            = var.region
  login {
    user_name = var.admin_user_name
    password  = var.admin_password
  }
}

provider "aws" {
  profile = var.profile
  region  = var.region
}

resource "aws_s3_bucket" "bucket-creation" {
  bucket = "my-tera-test-chaos2"
}

resource "aws_s3_bucket_object" "upload-file" {
  bucket = aws_s3_bucket.bucket-creation.id
  key    = "economic-survey-of-manufacturing-dec-2021.csv"
  source = "economic-survey-of-manufacturing-dec-2021.csv"
  etag = filemd5("economic-survey-of-manufacturing-dec-2021.csv")

}

resource "chaossearch_object_group" "create-object-group" {
  bucket = "test-object-group-tera2"
  source = "my-tera-test-chaos2"
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

resource "chaossearch_index_model" "chaossearch-og-index" {
  bucket_name = "test-object-group-tera2"
  model_mode = 0
  depends_on = [
    chaossearch_object_group.create-object-group
  ]
}

resource "chaossearch_view" "chaossearch-create-view" {
  bucket           = "test-view-tera2"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["test-object-group-tera2"]
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
    chaossearch_index_model.chaossearch-og-index
  ]
}

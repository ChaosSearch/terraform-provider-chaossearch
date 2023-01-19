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

/*
  For the cs provider to init, ensure automation has the following env variables:
    CS_URL 
    CS_ACCESS_KEY
    CS_SECRET_KEY
    CS_REGION
    CS_USERNAME
    CS_PASSWORD
    CS_PARENT_USER_ID (OPTIONAL) 
*/
provider "chaossearch" {
  login {}
}

// AWS infra needed as acc testing source bucket
resource "aws_s3_bucket" "bucket-creation" {
  bucket = "chaossearch-tf-provider-acc-test"
}

resource "aws_s3_bucket_object" "upload-file" {
  bucket = aws_s3_bucket.bucket-creation.id
  key    = "economic-survey-of-manufacturing-dec-2021.csv"
  source = "economic-survey-of-manufacturing-dec-2021.csv"
  etag   = filemd5("economic-survey-of-manufacturing-dec-2021.csv")
}


// ChaosSearch Object Group needed as a static source for view testing
resource "chaossearch_object_group" "create-object-group" {
  bucket = "acc-test-og-source"
  source = "chaossearch-tf-provider-acc-test"
  format {
    _type            = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = false
    field_selection = jsonencode([
      {
        "excludes" : [
          "data",
          "bigobject"
        ],
        "type" : "blacklist"
      }
    ])
    array_selection = jsonencode([
      {
        "excludes" : [
          "object.ids",
        ],
        "type" : "blacklist"
      }
    ])
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

// Generates index data and dataset for view test
resource "chaossearch_index_model" "model-1" {
  bucket_name = "acc-test-og-source"
  model_mode  = 0
}

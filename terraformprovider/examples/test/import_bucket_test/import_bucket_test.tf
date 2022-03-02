terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source  = "chaossearch/chaossearch"
    }
  }
}

variable "url" {}
variable "access_key_id" {}
variable "secret_access_key" {}
variable "region" {}
variable "admin_user_name" {}
variable "admin_password" {}

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

# Import Bucket
resource "chaossearch_import_bucket" "import_bucket1" {
  bucket      = "test-valid-bucket-name"
  hide_bucket = false
}


# Import Bucket
resource "chaossearch_import_bucket" "import_bucket2" {
  bucket      = "test-invalid-bucket-name"
  hide_bucket = false
}

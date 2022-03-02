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

# Create Sub Account
resource "chaossearch_sub_account" "sub-account" {
  user_info_block {
    username  = "nibras1102"
    full_name = "Nibras S"
    email     = "hello@test.com"
  }
  group_ids = ["aaa", "bbb"]
  password  = "1234"
  hocon     = ["override.Services.worker.quota=50"]
}




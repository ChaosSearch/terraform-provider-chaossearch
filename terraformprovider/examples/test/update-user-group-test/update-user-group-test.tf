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
resource "chaossearch_user_group" "chaossearch_user_group_update_test" {
  user_groups {
    id   = "4f20c1ac-538b-48ce-aa6e-1976a066551a"
    name = "update-test"
    permissions {
      effect    = "Allow"
      actions   = ["*."]
      resources = ["*."]
      version   = "1"

    }
  }
}
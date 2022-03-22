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
variable "user_name" {}
variable "password" {}
variable "parent_user_id" {}

provider "chaossearch" {
  url               = var.url
  access_key_id     = var.access_key_id
  secret_access_key = var.secret_access_key
  region            = var.region
  login {
    user_name      = var.user_name
    password       = var.password
    parent_user_id = var.parent_user_id
  }
}


#get all object groups
data "chaossearch_retrieve_object_groups" "first" {

}

output "object_group" {
  value = data.chaossearch_retrieve_object_groups.first
}
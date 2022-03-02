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

#get view group by id
data "chaossearch_retrieve_view" "retrieve-view" {
  bucket = "Chathura-view-10000tt"
}

output "view" {
  value = data.chaossearch_retrieve_view.retrieve-view
}


#without view id
#data "chaossearch_retrieve_view" "retrieve-view-without-id" {
#
#}
#
#output "view-without-id" {
#  value = data.chaossearch_retrieve_view.retrieve-view-without-id
#}

#when view id not exists
#data "chaossearch_retrieve_view" "view-not-found" {
#  bucket="Chat"
#}
#
#output "view-not-found" {
#  value = data.chaossearch_retrieve_view.view-not-found
#}


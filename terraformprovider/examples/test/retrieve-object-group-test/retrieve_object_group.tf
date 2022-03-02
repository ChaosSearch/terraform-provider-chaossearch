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


#get object group by id
data "chaossearch_retrieve_object_group" "object-group" {
  bucket = "c-og-100198"
}

output "object_group" {
  value = data.chaossearch_retrieve_object_group.object-group
}
#
#without object group id
#data "chaossearch_retrieve_object_group" "without-object-group-id" {
#}
#
#output "without_object_group_id" {
#  value = data.chaossearch_retrieve_object_group.without-object-group-id
#}


#when object group id not exists
#data "chaossearch_retrieve_object_group" "object-group-not-found" {
#  bucket="c-og-197"
#}
#
#output "object_group_not_found" {
#  value = data.chaossearch_retrieve_object_group.object-group-not-found
#}
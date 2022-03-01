terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source  = "chaossearch/chaossearch"
    }
  }
}
provider "chaossearch" {
  url               = "https://ap-south-1-aeternum.chaossearch.io"
  access_key_id     = "LCE8T6HRFGJI3ZKBGMGD"
  secret_access_key = "r5MEYkYntYvXqRSBMK6SFLQfPw7hHRQ0v5cqlkIk"
  region            = "ap-south-1"
  login {
    #    user_name      = "service_user@chaossearch.com"
    #    password       = "thisIsAnEx@mple1!"
    #    parent_user_id = "be4aeb53-21d5-4902-862c-9c9a17ad6675"
    user_name = "aeternum@chaossearch.com"
    password  = "ffpossgjjefjefojwfpjwgpwijaofnaconaonouf3n129091e901ie01292309r8jfcnsijvnsfini1j91e09ur0932hjsaakji"
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


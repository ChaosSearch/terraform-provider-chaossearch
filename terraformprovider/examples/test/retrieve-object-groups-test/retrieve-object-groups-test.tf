terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source = "chaossearch/chaossearch"
    }
  }
}
provider "chaossearch" {
  url               = "https://ap-south-1-aeternum.chaossearch.io"
  access_key_id     = "LCE8T6HRFGJI3ZKBGMGD"
  secret_access_key = "r5MEYkYntYvXqRSBMK6SFLQfPw7hHRQ0v5cqlkIk"
  region            = "ap-south-1"
  login  {
    user_name = "service_user@chaossearch.com"
    password = "thisIsAnEx@mple1!"
    parent_user_id = "be4aeb53-21d5-4902-862c-9c9a17ad6675"
  }

}


#get all object groups
data "chaossearch_retrieve_object_groups" "first" {

}

output "object_group" {
  value = data.chaossearch_retrieve_object_groups.first
}
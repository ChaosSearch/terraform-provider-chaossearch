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

# Create Sub Account
resource "chaossearch_sub_account" "sub-account" {
  user_info_block {
    username  = "nibras"
    full_name = "Nibras S"
    email     = "hello@test.com"
  }
  group_ids = ["aaa", "bbb"]
  password  = "1234"
  hocon     = ["override.Services.worker.quota=50"]
}




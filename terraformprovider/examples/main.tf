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

    # Normal User Credentials
    #    user_name      = "service_user@chaossearch.com"
    #    password       = "thisIsAnEx@mple1!"
    #    parent_user_id = "be4aeb53-21d5-4902-862c-9c9a17ad6675"

    # ADMIN Credentials
    user_name = "aeternum@chaossearch.com"
    password  = "ffpossgjjefjefojwfpjwgpwijaofnaconaonouf3n129091e901ie01292309r8jfcnsijvnsfini1j91e09ur0932hjsaakji"
  }

}

#resource "chaossearch_user_group" "chaossearch_user_group-crate" {
#  user_groups {
#    id   = "100044"
#    name = "dinesh-jayddasinghe"
#    permissions {
#      permission {
#
#        effect    = "Allow"
#        action    = "*"
#        resources = "*"
#        version="1.2"
#      }
#    }
#  }
#}
/*
resource "chaossearch_user_group" "chaossearch_user_group-crate" {
  user_groups {
    id   = "7db91912-a3e9-4641-873c-3deccd07484c"
    name = "Foo"
    permissions {
      permission {
        effect    = "Allow"
        action    = "kibana:*"
        resources = "crn:view:::foo-view"
        conditions {
          condition {
            starts_with {
              chaos_document_attributes_title = "foo"
            }
            equals {
              chaos_document_attributes_title = "bar"
            }
            not_equals  {
              chaos_document_attributes_title = "baz"
            }
            like  {
              chaos_document_attributes_title = "foobar"
            }

          }
        }
      }
      permission {
        effect    = "Allow1"
        action    = "kibana1:*"
        resources = "crn:view:::foo-view1"
        conditions {
          condition {
            starts_with {
              chaos_document_attributes_title = "foo1"
            }
            equals {
              chaos_document_attributes_title = "bar1"
            }
            not_equals  {
              chaos_document_attributes_title = "baz1"
            }
            like  {
              chaos_document_attributes_title = "foobar1"
            }

          }
        }
      }
    }
  }

}*/


#create view
#resource "chaossearch_view" "chaossearch-create-view" {
#  bucket = "Chathura-view-update-1"
#  case_insensitive = false
#  index_pattern   = ".*1112223344"
#  index_retention = -1
#  overwrite       = true
#  sources         = []
#  time_field_name = "@timestamp"
#  transforms      = []
#  filter {
#    predicate {
#      _type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
#      pred {
#        _type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
#        field = "cs_partition_key_0"
#        query = "*bluebike*"
#        state {
#          _type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
#        }
#      }
#    }
#  }
#}

# Import Bucket
#resource "chaossearch_import_bucket" "import_bucket" {
#  bucket      = "chaos-tera-test-123"
#  hide_bucket = false
#}

# Create Sub Account
#resource "chaossearch_sub_account" "sub-account" {
#  user_info_block {
#    username = "nibras"
#    full_name = "Nibras S"
#    email = "hello@test.com"
#  }
#  group_ids = ["aaa", "bbb"]
#  password = "1234"
#}



#create object group
#resource "chaossearch_object_group" "my-object-group" {
#
#  bucket = "Chathura-og-update-1"
#  source = "chaos-test-data-aps1"
#  format {
#    _type            = "CSV"
#    column_delimiter = ","
#    row_delimiter    = "\n"
#    header_row       = true
#  }
#  interval {
#    mode   = 0
#    column = 0
#  }
#  index_retention {
#    for_partition = []
#    overall       = -1
#  }
#  filter {
#    prefix_filter {
#      field  = "key"
#      prefix = "bluebike"
#    }
#    regex_filter {
#      field = "key"
#      regex = ".*"
#    }
#  }
#  options {
#    ignore_irregular = true
#  }
#  realtime = false
#}
#
#resource "chaossearch_view" "chaossearch-create-view" {
#bucket = "dinesh-view-09"
#
#  case_insensitive = false
#  index_pattern    = ".*"
#  index_retention  = -1
#  overwrite        = true
#  sources          = []
#  time_field_name  = "@timestamp"
#  transforms       = []
#  filter {
#    predicate {
#      _type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
#      pred {
#        _type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
#        field = "cs_partition_key_0"
#        query = "*bluebike*"
#        state {
#          _type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
#        }
#      }
#    }
#  }
#  depends_on = [
#    chaossearch_object_group.my-object-group
#  ]
#  }
#
#
#resource "chaossearch_object_group" "my-object-group" {
#  bucket = "dinesh-og-09"
#  source = "chaos-test-data-aps1"
#  format {
#    _type            = "CSV"
#    column_delimiter = ","
#    row_delimiter    = "\n"
#    header_row       = true
#  }
#  interval {
#    mode   = 0
#    column = 0
#  }
#  index_retention {
#    for_partition = []
#    overall       = -1
#  }
#  filter {
#    prefix_filter {
#      field  = "key"
#      prefix = "bluebike"
#    }
#    regex_filter {
#      field = "key"
#      regex = ".*"
#    }
#  }
#  options {
#    ignore_irregular = true
#  }
#  depends_on = [
#    chaossearch_object_group.my-object-group
#  ]
#}

#get all object groups
#data "chaossearch_retrieve_object_groups" "first" {
#
#}
#
#output "object_group_retrieve_object_groups" {
#  value = data.chaossearch_retrieve_object_groups.first
#}

#get all views

#data "chaossearch_retrieve_views" "myview" {
#
#}
#
#output "views" {
#  value = data.chaossearch_retrieve_views.myview
#}

//get all sub accounts
#data "chaossearch_retrieve_sub_accounts" "first" {
#
#}
#
#output "object_group_retrieve_object_groups" {
#  value = data.chaossearch_retrieve_sub_accounts.first
#}

//get all user groups
#data "chaossearch_retrieve_groups" "user_groups" {
#
#}
#
#output "chaossearch_retrieve_groups" {
#  value = data.chaossearch_retrieve_groups.user_groups
#}

##get user group by id
data "chaossearch_retrieve_user_group" "my-user-group" {
  id="9436aed9-e994-4dba-a25b-7d950d7f3623"
}

output "object_group_retrieve_user_group" {
  value = data.chaossearch_retrieve_user_group.my-user-group
}



#get object group by id
#data "chaossearch_retrieve_object_group" "my-object-group" {
#  bucket="c-og-100198"
#}
#

#output "object_group_retrieve_object_group" {
#  value = data.chaossearch_retrieve_object_group.my-object-group
#}


#get view  by id
#data "chaossearch_retrieve_view" "my-view1" {
#bucket="c-view-01"
#}
#
#output "object_group_retrieve_view" {
#  value = data.chaossearch_retrieve_view.my-view1
#}

#without view  id
#data "chaossearch_retrieve_view" "my-og2" {
#
#}
#
#output "object_group_retrieve_view_without_id" {
#  value = data.chaossearch_retrieve_view.my-og2
#}


#without object group id
#data "chaossearch_retrieve_object_group" "without-object-group-id" {
#}
#
#output "without_object_group_id" {
#  value = data.chaossearch_retrieve_object_group.without-object-group-id
#}

#resource "chaossearch_object_group" "my-object-group-1" {
#
#  bucket = "nibras-og-0142"
#  source = "chaos-test-data-aps1"
#  format {
#    _type            = "CSV"
#    column_delimiter = ","
#    row_delimiter    = "\n"
#    header_row       = true
#  }
#  interval {
#    mode   = 0
#    column = 0
#  }
#  index_retention {
#    for_partition = []
#    overall       = -1
#  }
#  filter {
#    prefix_filter {
#      field  = "key"
#      prefix = "bluebike"
#    }
#    regex_filter {
#      field = "key"
#      regex = ".*"
#    }
#  }
#  options {
#    ignore_irregular = true
#  }
#  realtime = false
#}
#
#resource "chaossearch_object_group" "my-object-group-2" {
#
#  bucket = "nibras-og-0251"
#  source = "chaos-test-data-aps1"
#  format {
#    _type            = "CSV"
#    column_delimiter = ","
#    row_delimiter    = "\n"
#    header_row       = true
#  }
#  interval {
#    mode   = 0
#    column = 0
#  }
#  index_retention {
#    for_partition = []
#    overall       = -1
#  }
#  filter {
#    prefix_filter {
#      field  = "key"
#      prefix = "bluebike"
#    }
#    regex_filter {
#      field = "key"
#      regex = ".*"
#    }
#  }
#  options {
#    ignore_irregular = true
#  }
#  realtime = false
#}



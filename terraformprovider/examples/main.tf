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

##create view
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



#get object group by id
data "chaossearch_retrieve_object_group" "my-object-group" {
  bucket="c-og-100198"
}

output "object_group_retrieve_object_group" {
  value = data.chaossearch_retrieve_object_group.my-object-group
}


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



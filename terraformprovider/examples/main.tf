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


#resource "chaossearch_view" "chaossearch-create-view" {
#  bucket = "Chathura-view-1000ropoofl99hdfi"
#  case_insensitive = false
#  index_pattern   = ".*"
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
#}

resource "chaossearch_object_group" "my-object-group" {
  bucket = "test-ab-0221"
  source = "chaos-test-data-aps1"
  format {
    _type            = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = true
  }
  interval {
    mode   = 0
    column = 0
  }
  index_retention {
    for_partition = []
    overall       = -1
  }
  filter {
    prefix_filter {
      field  = "key"
      prefix = "bluebike"
    }
    regex_filter {
      field = "key"
      regex = ".*"
    }
  }
  options {
    ignore_irregular = true
  }
  realtime = false
  index_parallelism = 1
  index_retention_value = 2
  target_active_index = 90
  live_events_parallelism = 10
}
#get all object groups
#data "chaossearch_object_groups" "first" {
#
#}
#
#output "object_group" {
#  value = data.chaossearch_object_groups.first
#}

#get all views

#data "chaossearch_views" "myview" {
#
#}
#
#output "views" {
#  value = data.chaossearch_views.myview
#}



#get object group by id
#data "chaossearch_object_group" "my-object-group" {
#  object_group_id="c-og-100197"
#}
#
#output "object_group" {
#  value = data.chaossearch_object_group.my-object-group
#}


//get view group by id
#data "chaossearch_view" "my-view" {
#
#}
#
#output "object_group" {
#  value = data.chaossearch_view.my-view
#}


#without object group id
data "chaossearch_object_group" "without-object-group-id" {
}

output "without_object_group_id" {
  value = data.chaossearch_object_group.without-object-group-id
}



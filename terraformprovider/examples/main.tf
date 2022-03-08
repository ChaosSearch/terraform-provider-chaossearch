terraform {
  required_providers {
    chaossearch = {
      version = "~> 0.1.1"
      source  = "chaossearch/chaossearch"
    }
  }
}

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

resource "chaossearch_view" "chaossearch-create-view-test1" {
  bucket           = "test_view_011"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = []
  time_field_name  = "@timestamp"
  transforms       = []
  filter {
    predicate {
      _type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
      pred {
        _type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
        field = "cs_partition_key_0"
        query = "*bluebike*"
        state {
          _type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
        }
      }
    }
  }
}



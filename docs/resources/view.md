# View Resource

> Create Refinery views to select the data for visualizations and to virtually transform schema details for analytics.

Creates a searchable and visualizable view for your _Index_ data from your _Object Group_

Check out the _View_ documentation here: [View Docs](https://docs.chaossearch.io/docs/refinery-index-views)

## Example Usage
```hcl
resource "chaossearch_view" "chaossearch-create-view" {
  bucket           = "tf-provider-view"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
  time_field_name  = "@timestamp"
  filter {
    predicate {
      type = "chaossumo.query.NIRFrontend.Request.Predicate.Negate"
      pred {
        type = "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
        field = "STATUS"
        query = "*F*"
        state {
          type = "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
        }
      }
    }
  }
}
```

## Argument Reference
* `

## Attribute Reference
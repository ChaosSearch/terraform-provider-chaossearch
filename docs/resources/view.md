# View Resource

> Create Refinery views to select the data for visualizations and to virtually transform schema details for analytics.

Creates a searchable and visualizable view for your _Index_ data from your _Object Group_

Check out the _View_ documentation here: [View Docs](https://docs.chaossearch.io/docs/refinery-index-views)

## Example Usage

Below is an example of making a view with a single `predicate`
```hcl
resource "chaossearch_view" "view" {
  bucket           = "tf-provider-view"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
  time_field_name  = "timestamp"
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

Below is an example of making a view with multiple `predicates`
```hcl
resource "chaossearch_view" "view-preds" {
  bucket           = "tf-provider-view-preds"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
  time_field_name  = "timestamp"
  filter {
    predicate {
      type = "chaossumo.query.NIRFrontend.Request.Predicate.Or"
      preds = [
        jsonencode(
          {
            "state": {
              "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
            },
            "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch",
            "field": "Subject",
            "query": "subject"
          }
        ),
        jsonencode(
          {
            "state": {
              "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
            },
            "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch",
            "field": "Series_title_1",
            "query": "title"
          }
        ),
        jsonencode(
          {
            "preds": [
              {
                "field": "test_id",
                "query": "1",
                "state": {
                  "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
                },
                "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
              },
              {
                "field": "test_id",
                "query": "2",
                "state": {
                  "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
                },
                "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
              }
            ],
            "_type": "chaossumo.query.NIRFrontend.Request.Predicate.Or"
          }
        )
      ]
    }
  }
}
```

Below is an example view with `transforms`
```hcl
resource "chaossearch_view" "view-transforms" {
  bucket           = "tf-provider-view"
  case_insensitive = false
  index_pattern    = ".*"
  index_retention  = -1
  overwrite        = true
  sources          = ["tf-provider"]
  time_field_name  = "timestamp"
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
  transforms = [
    jsonencode({
        "_type": "PartitionKeyTransform"
        "keyPart": 0
        "inputField": "cs_partition_key_0"
    }),
    jsonencode({
      "_type": "MaterializeRegexTransform",
      "inputField": "Data_value",
      "pattern": "(\\d+)\\.(\\d+)"
      "outputFields": [
        {
          "name": "Whole",
          "type": "NUMBER"
        },
        {
          "name": "Decimal",
          "type": "NUMBER"
        }
      ]
    }),
    jsonencode({
      "_type": "MaterializeJQTransform",
      "inputField": "Data_value",
      "queries": ["jq-query"],
      "outputFields": [
        {
          "name": "Whole",
          "type": "NUMBER"
        }
      ]
    }),
    jsonencode({
      "_type": "MaterializeJSONTransform",
      "inputField": "Data_value",
      "paths": ["json-path"],
      "outputFields": [
        {
          "name": "Whole",
          "type": "NUMBER"
        }
      ]
    })
  ]
}
```

## Argument Reference
* `bucket` - **(Required)** The name of the view bucket
* `case-insensitive` - **(Required)** Declares whether or not attributes during view querying are case-sensitive
* `index-pattern`
* `index-retention` - **(Required)** Determines the number of days an indexes will be retained
  * `-1` For indefinite retention
* `overwrite`
* `sources` - **(Required)** The `object groups` used to provide views with data
* `time_field_name` - **(Required)** The data's attribute to be used as a timestamp
* `filter` - **(Optional)** This object houses any applied filtering to the views
  * `predicate` - **(Required)** Houses predicates for filtering
    * `type` - **(Optional)** Indicates the type of relationship the `preds` or `pred` will have
      * Accepted Values include... 
        * And: `chaossumo.query.NIRFrontend.Request.Predicate.AND`
        * Or: `chaossumo.query.NIRFrontend.Request.Predicate.Or`
        * Negate: `chaossumo.query.NIRFrontend.Request.Predicate.Negate`
        * TextMatch: `chaossumo.query.NIRFrontend.Request.Predicate.TextMatch`
        * Range: `chaossumo.query.NIRFrontend.Request.Predicate.Range`
        * Exists: `chaossumo.query.NIRFrontend.Request.Predicate.Exists`
      * `And` and `Or` are primarily used in the case where there are multiple `preds`
        * Note: If used with a single `pred`, API will throw a `JSON Parse` error
    * `pred` - **(Optional)** Used in the case where only one field is being filtered
    * `preds` - **(Optional)** Used in the case where multiple fields are being filtered
      * Takes in an array of json
      * Follows the same structure as `pred`, but also enables the ability to nest more `preds`
* `transforms` - **(Optional)** Takes in an array populated wth a json object for each transform
  * `_type` - **(Required)** Defines the type of transform, different types include:
    * `PartitionKeyTransform` - Takes in `keyPart`
    * `MaterializeRegexTransform` - Takes in `pattern`
    * `MaterializeSJSONTransform` - Takes in `paths` 
    * `VerticalArrayTransform` - Takes in `vertical`
    * `IPFieldTransform`
    * `GeoPointFieldTransform` - Takes in `format`
    * `NestedFieldTransform`
  * `inputField` - **(Required)** The field you are transforming
  * `outputFields` - **(Optional)** Used when transforming one input field to many outputs
    * `name` - **(Required)**
    * `type` - **(Required)**
  * `keyPart` - See `PartitionKeyTransform`, takes an integer
  * `pattern` - See `MaterializeRegexTransform`, takes regex string
  * `paths` - See `MaterializeSJSONTransform`, takes an array of strings
  * `vertical` - See `VerticalArrayTransform`, takes an array of strings
  * `format` - See `GeoPointFieldTransform`, takes a float 
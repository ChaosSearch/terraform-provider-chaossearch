# Object Group Resource

> Object groups organize and associate similar data files in your cloud storage for ChaosSearch indexing.

Creates an _Object Group_ to structure your data for Indexing

Check out the _Object Group_ documentation here: [Object Group Docs](https://docs.chaossearch.io/docs/creating-object-groups)

## Example Usage
```hcl
resource "chaossearch_object_group" "create-object-group" {
  bucket = "tf-provider"
  source = "chaossearch-tf-provider-test"
  live_events_aws = "arn:partition:service:region:account-id:resource-id"
  live_events_gcp {
    project_id = "test_id"
    subscription_id = "some_test_id"
  }
  format {
    type             = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = true

    # The following selection policies can all contain the same json properties (including col_selection)
    field_selection = jsonencode([{
        "excludes": [
          "data",
          "bigobject"
        ],
        "type": "blacklist"
    }])
    array_selection = jsonencode([{
      "includes": [
        "object.ids",
      ],
      "type": "whitelist"
    }])
    vertical_selection = jsonencode([{
      "include": true,
       "patterns": [
        "^line\\.level$",
        "^attrs.version$",
        "^timestamp$",
        "^line\\.meta\\.[^\\.]*$",
        "^host$",
        "^line\\.correlation_id$",
        "^sourcetype$",
        "^line\\.message$",
        "^message$",
        "^source$",
        "^_rawJson$"
      ],
      "type": "regex"
    }])
  }
  index_retention {
    overall       = -1
  }
  filter {
    field = "key"
    prefix = "ec"
  }
  filter {
    field = "key"
    regex = ".*"
  }
  filter {
    field = "storageClass"
    equals = "STANDARD"
  }
  options {
    compression = "GZIP"
    col_types = jsonencode({
      "TimeStamp": "Timeval"
    })
    col_renames = jsonencode({
      "TimeStamp": "Period"
    })
    col_selection = jsonencode({
      "includes": [
        "object.ids",
      ],
      "type": "whitelist"
    })
  }
}
```

**Note:** Update calls are only made when changes are detected for the following values (Requires OG/Index Model to be paused to apply):
* `index_retention.overall`
* `target_active_index`
* `index_parallelism`
* `live_events_parallelism` 

## Argument Reference
* `bucket` - **(Required)** Name of the object group
* `source` - **(Required)** Name of the bucket where your data is stored
* `live_events` - **(Optional)** **DEPRECATED** Renamed as `live_events_aws`
* `live_events_aws` - **(Optional)** The SQS Arn for live event streaming
* `live_events_gcp` - **(Optional)** A config block for defining live events in GCP
  * `project_id` - **(Required)** Your GCP project ID
  * `subscription_id` - **(Required)** Your subscription topic ID
* `format` - **(Optional)** A config block used for file formatting specifics
  * `type` - **(Optional)** Specifies the type of file
  * `column_delimiter` - **(Optional)** Specifies the character for separating columns
  * `row_delimiter` - **(Optional)** Specifies the character for separating rows
  * `header_row` - **(Optional)** Specifies if the file includes a header row
  * `array_flatten_depth` - **(Optional)** How deeply nested arrays should be allowed to get before parsing stops. Defaults to 0. Use `-1` for unlimited
  * `horizontal` - **(Optional)** If true, array fields will be turned into new columns on each flattened message. If false, array fields will be broadcast into multiple flattened rows for each array item.
  * `array_selection` - **(Optional)** A json policy block for selecting array fields
  * `field_selection` - **(Optional)** A json policy block for selecting object fields
* `index_retention` - **(Optional)** Config block for specifying how long an index is retained
  * Only applies on update
  * `overall` - **(Optional)** Takes the amount of days an index is retained
    * use `-1` for an indefinite amount of time
* `filter` - **(Optional)** Config block for housing filtering
  * Note: Make sure that `prefix`, `regex` and `equals` are all broken into their own filter block
  * `field` - **(Required)** What field the filter applies to
    * Can be `key` and `storageClass`
  * `prefix` - **(Optional)** Used with `key` field. The prefix the field must match for the file
  * `regex` - **(Optional)** Used with `key` field. The regex for filtering files 
  * `equals` - **(Optional)** Used with `storageClass` field. Supplies the `storageClass` type of the S3 bucket
    * Can be `STANDARD`, `STANDARD_IA`, `INTELLIGENT_TIERING`, `ONEZONE_IA`, `GLACIER`, `DEEP_ARCHIVE`, `REDUCED_REDUNDANCY`
* `options` - **(Optional)** A config block for housing advanced settings
  * `col_renames` - **(Optional)** A set of key value pairs, key being new name, val being old name
  * `col_types` - **(Optional)** A set of key value pairs, key being field name, val being field type
  * `col_selection` - **(Optional)** A json policy block for selecting column fields 
  * `compression` - **(Optional)** Form of file compression being used
    * Can either be `GZIP` or `SNAPPY`
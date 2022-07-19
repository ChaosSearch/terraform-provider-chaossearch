# Object Group Resource

> Object groups organize and associate similar data files in your cloud storage for ChaosSearch indexing.

Creates an _Object Group_ to structure your data for Indexing

Check out the _Object Group_ documentation here: [Object Group Docs](https://docs.chaossearch.io/docs/creating-object-groups)

## Example Usage
```hcl
resource "chaossearch_object_group" "create-object-group" {
  bucket = "tf-provider"
  source = "chaossearch-tf-provider-test"
  format {
    type             = "CSV"
    column_delimiter = ","
    row_delimiter    = "\n"
    header_row       = true
  }
  index_retention {
    overall       = -1
  }
  filter {
    prefix_filter {
      field = "key"
      prefix = ".*"
    }
    regex_filter {
      field = "key"
      regex = ".*"
    }
  }
}
```

## Argument Reference
* `bucket` - **(Required)** Name of the object group
* `source` - **(Required)** Name of the bucket where your data is stored
* `format` - **(Optional)** A config block used for file formatting specifics
  * `type` - **(Optional)** Specifies the type of file
  * `column_delimiter` - **(Optional)** Specifies the character for separating columns
  * `row_delimiter` - **(Optional)** Specifies the character for separating rows
  * `header_row` - **(Optional)** Specifies if the file includes a header row
* `index_retention` - **(Optional)** Config block for specifying how long an index is retained
  * Only applies on update
  * `overall` - **(Optional)** Takes the amount of days an index is retained
    * use `-1` for an indefinite amount of time
* `filter` - **(Optional)** Config block for housing filtering
  * `prefix_filter`
    * `field` - **(Required)** What field the filter applies to (usually `key`)
    * `prefix` - **(Required)** The prefix the field must match for the file
  * `regex_filter`
    * `field` - **(Required)** What field the filter applies to (usually `key`)
    * `regex` - **(Required)** The regex used for filtering files

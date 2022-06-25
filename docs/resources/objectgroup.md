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
      field = "key"
      prefix = ".*"
    }
    regex_filter {
      field = "key"
      regex = ".*"
    }
  }
  options {
    ignore_irregular = true
  }
}
```

## Argument Reference

## Attribute Reference

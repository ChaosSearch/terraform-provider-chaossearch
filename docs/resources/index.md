# Index Resource

> Index your cloud storage objects to create the ChaosSearch index files for visualizing and analyzing the information in the content.

This enables/starts indexing for a provided _Object Group_

Check out the _Index_ documentation here: [Index Docs](https://docs.chaossearch.io/docs/modeling-your-data)

## Example Usage

```hcl
resource "chaossearch_index_model" "index" {
  bucket_name    = ""
  model_mode     = 0
  delete_enabled = false
  delete_timeout = 0
}
```

## Argument Reference
* `bucket_name` - **(Required)** Name of _Object Group_ to be indexed
* `model_mode` - **(Required)** The mode you wish to set your index to. Values Include:
  * `-1` Restart Indexing
  * `0` Start Indexing
  * `1` Pause Indexing
* `delete_enabled` - **(Optional)** Enables or Disables _Index_ deletion
  * Defaults to `false`
  * **WARNING** Do not put `delete_enabled = true` in to source control.
  * Acts as a safeguard for data loss by disabling users from accidentally deleting indexes on `terraform destroy`
  * You will have to update this value in your `.tfstate` to delete your index
* `delete_timeout` - **(Optional)** Sets a Timeout limit (in seconds) for the _Index_ deletion confirmation
  * Defaults to `0`
  * This does not disable the _Index_ delete call, only times out the _Index_ delete confirmation
  * If this triggers an Exception, you will have to manually confirm _Index_ deletion in your ChaosSearch cluster

## Attribute Reference
* `indexed` - Reflects the status of whether or not an _Object Group_ has been indexed
  * `Boolean`
* `result` - Index Request confirmation
  * `Boolean`
  
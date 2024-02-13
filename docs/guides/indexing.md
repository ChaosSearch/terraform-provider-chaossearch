# Indexing Guide

If you want your object group's indexing to be managed by terraform, be sure to create an `index_model` definition. If not supplied, you will have to manually delete all of the indexes if you wish to tear down your object group. Be sure to define a `depends_on` for the `index_model`, this will ensure the provider does not attempt to start the indexing prior to the OG existing.

```hcl
resource "chaossearch_object_group" "ex-obj-group" {
  bucket = "og-name"
  source = "your-bucket"
  ...
}

resource "chaossearch_index_model" "model" {
  bucket_name = chaossearch_object_group.ex-obj-group.bucket
  model_mode  = 0
  depends_on  = [
    chaossearch_object_group.ex-obj-group
  ]
}
```
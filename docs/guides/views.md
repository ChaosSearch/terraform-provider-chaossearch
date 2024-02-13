# Views Guide

If you wish to define a view in your terraform, you can associate it's sources with object groups with the following example. Note, it would be best practice to define a depends_on between the indexes and object groups with the views. Otherwise the view will throw errors if either 1) The supplied OGs do not exist or 2) Indexes for it have not been generated yet.

```hcl
resource "chaossearch_object_group" "ex-obj-group" {
  ...
}

resource "chaossearch_object_group" "ex-obj-group-2" {
  ...
}

resource "chaossearch_index_model" "model" {
  ...
}

resource "chaossearch_index_model" "model-2" {
  ...
}

resource "chaossearch_view" "ex-view" {
  bucket           = "ex-view"
  sources          = [
    chaossearch_object_group.ex-obj-group.bucket,
    chaossearch_object_group.ex-obj-group-2.bucket
  ]
  ...
  depends_on = [
    chaossearch_index_model.model,
    chaossearch_index_model.model-2,
    chaossearch_object_group.ex-obj-group,
    chaossearch_object_group.ex-obj-group-2
  ]
}
```
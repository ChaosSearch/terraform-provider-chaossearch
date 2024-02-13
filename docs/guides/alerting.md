# Alerting Guide

If you want to manage your alerting via terraform, you can do it with the following example
```hcl
resource "chaossearch_view" "ex-view" {
    ...
}

resource "chaossearch_destination" "ex-dest" {
  ...
}

resource "chaossearch_monitor" "monitor" {
  name = "tf-provider-monitor"
  type = "monitor"
  enabled = true
  depends_on = [
    chaossearch_destination.ex-dest,
    chaossearch_view.ex-view
  ]
  inputs {
    search {
      indices = [
        chaossearch_view.ex-view.bucket,
      ]
      ...
    }
  }
  triggers {
    name = "ex-trigger"
    actions {
      name = "ex-action"
      destination_id = chaossearch_destination.ex-dest.id
      ...
    }
    ...
  }
}
```
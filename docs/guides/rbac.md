# RBAC Guide

If you wish to define a subaccount and map it to a user group, you may do so with the following.

```hcl
resource "chaossearch_user_group" "user-group" {
  name = "example-group-name"
  ...
}

resource "chaossearch_sub_account" "sub-account" {
  username  = ""
  full_name = ""
  password  = ""
  hocon     = ["override.Services.worker.quota=50"]
  group_ids = [chaossearch_user_group.user-group.id]
}
```
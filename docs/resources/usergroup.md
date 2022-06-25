# User Group Resource

>Use groups to configure role-based access controls within the ChaosSearch platform.

Creates a _User Group_ to manage RBAC permissions for _Subaccounts_

Check out the _User Group_ and RBAC documentation here: [RBAC Configuration Docs](https://docs.chaossearch.io/docs/role-based-access-control-configuration)

## Example Usage

```hcl
resource "chaossearch_user_group" "user_group" {
  name = "provider_test"
  permissions = jsonencode([
    {
      "Version"   = "1.0",
      "Effect"    = "Allow",
      "Actions"   = ["*"]
      "Resources" = ["*"],
      "Condition" = {
        "Conditions" = [
          {
            "Equals"     = {
              "chaos:document/attributes.title" = ""
            },
            "Like"       = {
              "chaos:document/attributes.title" = ""
            },
            "NotEquals"  = {
              "chaos:document/attributes.title" = ""
            },
            "StartsWith" = {
              "chaos:document/attributes.title" = "test"
            },
          }
        ]
      }
    }
  ])
}
```

## Argument Reference
* `name` - **(Required)** Name of the _User Group_
* `permissions` **(Optional)** RBAC permissions for _User Group_
  * Takes in json as permissions body
  * Follows the same structure as what's produced by our _User Group_ permission's Block Creator
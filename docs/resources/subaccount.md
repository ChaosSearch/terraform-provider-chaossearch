# Subaccount Resource

> ChaosSearch includes built-in authentication and authorization controls so that account administrators can create and manage static subaccounts and groups of access roles. There are user interfaces as well as API endpoints for managing subaccounts and groups.

Creates a _Subaccount_ 

Check out the _Subaccount_ documentation here: [Subaccount Docs](https://docs.chaossearch.io/docs/subaccount-users)

## Example Usage
```hcl
resource "chaossearch_sub_account" "sub-account" {
  username  = ""
  full_name = ""
  password  = ""
  hocon     = ["override.Services.worker.quota=50"]
}
```

## Argument Reference
* `username` - (Required) Username for _Subaccount_ login
* `password` - (Required) Password for _Subaccount_ login
* `full_name` - (Required) The name of the _Subaccount_
* `group_ids` - (Optional) A list of _User Group_ Ids to associate this account with
* `hocon` - (Optional) A list of overridable configurations
  * Recommended seek ChaosSearch support before using
  
## Attribute Reference
* `hocon_json` - The hocon configuration for the _Subaccount_ returned as json
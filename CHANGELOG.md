## 1.0.5

### Enhancements:
* docs/resources/*: Documentation has been included for resources
* provider/resources/objectgroup.go: Removed `interval` and `index_retention.for_partition` for simplicity and redundancy

## 1.0.4

### Features:
* provider/resources/view.go: Enable users to provide multiple `preds` in views
* provider/provider.go: Enable the usage of API Key driven Auth

### Enhancements:
* provider/resources/view.go: Simplify `ViewRequestDTO`, rename to `ViewData`
* client/clientrequest.go: Add a note to JWT parsing error indicating failed authentication
* client/clientrequest.go: Break up request construction to allow for API Key driven auth 
  
## 1.0.3

### Enhancements:
* CHANGELOG.md: Adding
* provider/resources/usergroup.go: Convert `chaossearch_user_group` resource to take in json for `permissions` attribute

## 1.0.2

### Bug Fixes:
* provider/resources/index.go: Catch error thrown by `client.ReadIndexModel()`

## 1.0.1

### Enhancements:
* provider/resources/objectgroup.go: Add checking to ensure `index_retention.overall` and `target_active_index` have valid values
* provider/resources/index.go: Enable a `delete_timeout` to be configured for index deletion
* provider/resources/usergroup.gp: Flatten `user_group` schema for easier maintenance and readability
* provider/tests/provider_test.go: Add random UUID to to generated resources names

### Bug Fixes:
* provider/resources/view.go: Pass in `cacheable` into `CreateViewRequest`
* provider/resources/usergroup.go: Enable `user_group` ids to be stored in state, allowing for proper updates
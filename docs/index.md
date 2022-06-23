# ChaosSearch Provider

Official provider for the ChaosSearch product.

### The Data Lake Platform for Analytics at Scale

> Know Better™. Now you can, while realizing the true promise of data lake economics for log analytics at scale, insightful product-led growth and agile BI.

Now **Automate Better** with our provider.

Check out the ChaosSearch documentation here: [Chaos Docs](https://docs.chaossearch.io/docs)

## Example Usage

```hcl
provider "chaossearch" {
  url               = "" 
  access_key_id     = "" 
  secret_access_key = "" 
  region            = "" 
  login {
    user_name       = "" 
    password        = ""
    parent_user_id  = ""
  }
}
```

## Argument Reference
The following all have environment variable default functions.
* `url` - **(Required)** Your ChaosSearch cluster's Url
  * Env Var -> `CS_URL`
  * e.g -> 'cluster.chaossearch.com'
* `access_key_id` - **(Required)** Your ChaosSearch user's Access Key ID
  * This can be found in 'Settings > API Keys'
  * Env Var -> `CS_ACCESS_KEY`
* `secret_access_key`- **(Required)** Your ChaosSearch user's Secret Access Key
  * This can be found in 'Settings > API Keys'
  *  Env Var -> `CS_SECRET_KEY`
* `region` - **(Required)** Your ChaosSearch cluster's deploy region. 
  * This can be found in 'Settings > (AWS/GCP) Credentials' 
  * Env Var -> `CS_REGION`
  * e.g -> `us-east-`
* `login` - **(Required)** Login block for housing credentials
  * `user_name` - **(Required)** Your ChaosSearch username
    * Env Var -> `CS_USERNAME`
  * `password` - **(Required)** Your ChaosSearch password
    * Env Var -> `CS_PASSWORD`
  * `parent_user_id` - **(Optional)** If using a _Subaccount_, provide the parent account's `uuid`
    * Env Var -> `CS_PARENT_USER_ID`

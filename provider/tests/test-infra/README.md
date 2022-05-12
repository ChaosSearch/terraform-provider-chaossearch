# Testing Infrastructure

Here is where we can keep Terraform code surrounding deploying test infrastructure.

This would allow for us to perform acceptance testing in any cloud environment, to ensure compatibility between them. Currently Terraform does not seem to support the ability to pre-provision resources to external providers, but this should also be a simple work around.

## How

Spinning up the test infrastructure should be as easy as 

```
terraform apply
```
**Note** For the CS provider to init in automation, the following environment variables will need to be defined:
    - CS_URL 
    - CS_ACCESS_KEY
    - CS_SECRET_KEY
    - CS_REGION
    - CS_USERNAME
    - CS_PASSWORD
    - CS_PARENT_USER_ID (OPTIONAL) 

### Troubleshooting

A common error to throw during acceptance testing is: 
    ```After applying this test step and performing a `terraform refresh`, the plan was not empty.```

Generally this points to issues in the `Create` that get papered over in `Refresh` prior to `Plan` running.
This is usually reproduceable through the CLI by doing a `terraform apply` and  a `terraform plan -refresh=false`.
We would expect the plan to be empty due to the apply, but this error shows it to not be the case.
# Destination Resource

> ChaosSearch provides the ability to send notifications to destinations such as third-party tools. You can define either a Slack integration or a Webhook integration in the Destinations section of the ChaosSearch Alerting area.

Create a _Destination_ for your alerts using webhooks

Check out the _Destination_ documentation here: [Destination Docs](https://docs.chaossearch.io/docs/alert-destinations)

## Example Usage

### Slack Integration
```hcl
resource "chaossearch_destination" "dest" {
  name = "tf-provider-destination"
  type = "slack"
  slack {
    url = "http://slack.com"
  }
}
```

### Custom Webhook w/ URL
```hcl
resource "chaossearch_destination" "dest_custom" {
  name = "tf-provider-destination-custom"
  type = "custom_webhook"
  custom_webhook {
    url = "http://test.com"
  }
}
```

### Custom Webhook
```hcl
resource "chaossearch_destination" "dest_custom_host" {
  name = "tf-provider-destination-custom-host"
  type = "custom_webhook"
  custom_webhook {
    scheme = "HTTPS"
    host = "test.com"
    path = "/api/test"
    port = "8080"
    method = "POST"
    query_params = jsonencode({
      "test": "value"
    })
    header_params = jsonencode({
      "Content-Type": "application/json"
    })
  }
}
```

**Note** Slack and custom_webhook configs cannot both be declared in the same resource block
**Note** When using a custom_webhook, url is not needed if host, path, and/or port are defined
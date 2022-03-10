package client

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func (c *CSClient) Set(data *schema.ResourceData, key string, value interface{}) {
	err := data.Set(key, value)
	if err != nil {
		return
	}
}

package main

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceObjectGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectGroupsRead,
		Schema: map[string]*schema.Schema{
			"object_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// TODO BucketType (object group, view, native s3 bucket)
						// TODO Predicate
					},
				},
			},
		},
	}
}

func dataSourceObjectGroupsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// todo to be implemented
	client := meta.(*ProviderMeta).Client
	tokenValue := meta.(*ProviderMeta).token

	clientResponse, err := client.ListBuckets(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, len(clientResponse.BucketsCollection.Buckets))
	for i := 0; i < len(clientResponse.BucketsCollection.Buckets); i++ {
		result[i] = map[string]interface{}{
			"id":   clientResponse.BucketsCollection.Buckets[i].Name,
			"name": clientResponse.BucketsCollection.Buckets[i].Name,
		}
	}

	var diags diag.Diagnostics

	objectGroups := result

	if err := data.Set("object_groups", objectGroups); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

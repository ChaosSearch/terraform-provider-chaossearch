package cs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceObjectGroupsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	clientResponse, err := client.ListBuckets(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	result := GetBucketData(clientResponse)
	var diags diag.Diagnostics
	objectGroups := result
	if err := data.Set("object_groups", objectGroups); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

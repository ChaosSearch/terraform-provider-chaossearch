package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

func dataSourceSubAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: readAllSubAccounts,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sub_accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activated": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"full_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hocon": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func readAllSubAccounts(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	usersResponse, err := client.ListUsers(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, len(usersResponse.Users))
	for i := 0; i < len(usersResponse.Users); i++ {
		log.Info(result)
		//result[i] = map[string]interface{}{
		//	"id":   usersResponse.BucketsCollection.Buckets[i].Name,
		//	"name": usersResponse.BucketsCollection.Buckets[i].Name,
		//}
	}

	return nil
}

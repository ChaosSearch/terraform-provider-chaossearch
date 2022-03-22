package cs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceSubAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: readAllSubAccounts,
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
	}
}

func readAllSubAccounts(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	usersResponse, err := client.ListUsers(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, len(usersResponse.Users[0].SubAccounts))
	for i := 0; i < len(usersResponse.Users[0].SubAccounts); i++ {
		subAccount := usersResponse.Users[0].SubAccounts[i]
		result[i] = map[string]interface{}{
			"activated": subAccount.Activated,
			"full_name": subAccount.FullName,
			"hocon":     subAccount.Hocon,
			"uid":       subAccount.UID,
			"username":  subAccount.Username,
			"group_ids": subAccount.GroupIds,
		}
	}
	var diags diag.Diagnostics
	subAccounts := result
	if err := data.Set("sub_accounts", subAccounts); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

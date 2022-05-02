package datasources

import (
	"context"
	"cs-tf-provider/provider/models"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSubAccounts() *schema.Resource {
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
	var diags diag.Diagnostics
	var result []map[string]interface{}

	client := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	usersResponse, err := client.ListUsers(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	users := usersResponse.Users
	if len(users) > 0 {
		subAccounts := users[0].SubAccounts
		result = make([]map[string]interface{}, len(subAccounts))
		for i, subAccount := range subAccounts {
			result[i] = map[string]interface{}{
				"activated": subAccount.Activated,
				"full_name": subAccount.FullName,
				"hocon":     subAccount.Hocon,
				"uid":       subAccount.UID,
				"username":  subAccount.Username,
				"group_ids": subAccount.GroupIds,
			}
		}
	}

	if err := data.Set("sub_accounts", result); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

package datasources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: readUserGroup,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: readUserGroups,
		Schema: map[string]*schema.Schema{
			"user_groups": {
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
						"permissions": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readUserGroups(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(client.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenValue := meta.(*models.ProviderMeta).Token
	usersResponse, err := client.ListUsers(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	users := usersResponse.Users
	if len(users) > 0 {
		userGroups := make([]map[string]interface{}, len(users[0].UserGroups))
		for i, userGroup := range users[0].UserGroups {
			permissions, err := json.Marshal(userGroup.Permissions)
			if err != nil {
				return diag.FromErr(utils.MarshalJsonError(err))
			}

			userGroups[i] = map[string]interface{}{
				"id":          userGroup.ID,
				"name":        userGroup.Name,
				"permissions": string(permissions),
			}
		}

		if err := data.Set("user_groups", userGroups); err != nil {
			return diag.FromErr(err)
		}
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func readUserGroup(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	id := data.Get("id").(string)
	data.SetId(id)
	diags := diag.Diagnostics{}
	tokenValue := meta.(*models.ProviderMeta).Token
	req := &client.BasicRequest{
		AuthToken: tokenValue,
		Id:        id,
	}

	resp, err := c.ReadUserGroup(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		return diag.Errorf("Couldn't find User Group: %s", err)
	}

	userGroupContent, err := resources.CreateUserGroupResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("name", userGroupContent[0]["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("permissions", userGroupContent[0]["permissions"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

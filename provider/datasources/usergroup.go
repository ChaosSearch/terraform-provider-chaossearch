package datasources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{

			"user_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
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
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
									"effect": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
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

func DataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: readAllUserGroups,
		Schema: map[string]*schema.Schema{
			"groups": {
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
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
									"effect": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
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

func readAllUserGroups(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	usersResponse, err := client.ListUsers(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, len(usersResponse.Users[0].UserGroups))
	for i := 0; i < len(usersResponse.Users[0].UserGroups); i++ {
		permissionArr := make([]interface{}, 1)
		permissionMap := make(map[string]interface{})

		group := usersResponse.Users[0].UserGroups[i]

		permissionMap["version"] = group.Permissions[0].Version
		permissionMap["resources"] = group.Permissions[0].Resources
		permissionMap["effect"] = group.Permissions[0].Effect
		permissionMap["actions"] = group.Permissions[0].Actions
		permissionArr[0] = permissionMap

		result[i] = map[string]interface{}{
			"id":          group.ID,
			"name":        group.Name,
			"permissions": permissionArr,
		}
	}

	var diags diag.Diagnostics
	userGroups := result
	if err := data.Set("groups", userGroups); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func dataSourceUserGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userGroupInterface := data.Get("user_groups").(*schema.Set).List()[0].(map[string]interface{})
	data.SetId(userGroupInterface["id"].(string))
	diags := diag.Diagnostics{}
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	req := &client.ReadUserGroupRequest{
		AuthToken: tokenValue,
		ID:        data.Id(),
	}
	resp, err := c.ReadUserGroup(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.Errorf("Couldn't find User Group: %s", err)
	}
	userGroupContent := resources.CreateUserGroupResponse(resp)
	if err := data.Set("user_groups", userGroupContent); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

package main

import (
	"context"
	//"github.com/google/martian/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceUserGroups() *schema.Resource {
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
	client := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	usersResponse, err := client.ListUsers(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, len(usersResponse.Users[0].Groups))
	for i := 0; i < len(usersResponse.Users[0].Groups); i++ {
		permissionArr := make([]interface{}, 1)
		permissionMap := make(map[string]interface{})

		group := usersResponse.Users[0].Groups[i]

		permissionMap["version"] = group.Permissions[0].Version
		permissionMap["resources"] = group.Permissions[0].Resources
		permissionMap["effect"] = group.Permissions[0].Effect
		permissionMap["actions"] = group.Permissions[0].Actions
		permissionArr[0] = permissionMap

		result[i] = map[string]interface{}{
			"id":          group.Id,
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

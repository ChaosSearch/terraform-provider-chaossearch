package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"permissions": {
							Type:     schema.TypeSet,
							Optional: true,
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
										Optional: true,
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
										Optional: true,
									},
									"conditions": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"starts_with": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"equals": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"not_equals": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"like": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
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

func resourceGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	resp, err := c.CreateUserGroup(ctx, GroupObject(data, tokenValue))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.ID)
	return resourceGroupRead(ctx, data, meta)
}

func GroupObject(data *schema.ResourceData, authToken string) *client.CreateUserGroupRequest {
	var clientConditionGroup client.ConditionGroup
	var clientEquals *client.Equals
	var clientStartsWith *client.StartsWith
	var clientNotEquals *client.NotEquals
	var clientLike *client.Like
	var permissionList []client.Permission

	userGroupList := data.Get("user_groups").(*schema.Set).List()
	if len(userGroupList) > 0 {
		userGroupMap := userGroupList[0].(map[string]interface{})
		permissionsList := userGroupMap["permissions"].(*schema.Set).List()
		if len(permissionsList) > 0 {
			for _, permission := range permissionsList {
				permissionMap := permission.(map[string]interface{})
				conditionsList := permissionMap["conditions"].(*schema.Set).List()
				if len(conditionsList) > 0 {
					conditionsMap := conditionsList[0].(map[string]interface{})
					conditionList := conditionsMap["condition"].(*schema.Set).List()
					if len(conditionList) > 0 {
						conditionMap := conditionList[0].(map[string]interface{})
						equalsList := conditionMap["equals"].(*schema.Set).List()
						if len(equalsList) > 0 {
							equalsMap := equalsList[0].(map[string]interface{})
							clientEquals = &client.Equals{
								ChaosDocumentAttributesTitle: equalsMap["chaos_document_attributes_title"].(string),
							}
						}

						startsWithList := conditionMap["starts_with"].(*schema.Set).List()
						if len(startsWithList) > 0 {
							startsWithMap := startsWithList[0].(map[string]interface{})
							clientStartsWith = &client.StartsWith{
								ChaosDocumentAttributesTitle: startsWithMap["chaos_document_attributes_title"].(string),
							}
						}

						notEqualsList := conditionMap["not_equals"].(*schema.Set).List()
						if len(notEqualsList) > 0 {
							notEqualsMap := notEqualsList[0].(map[string]interface{})
							clientNotEquals = &client.NotEquals{
								ChaosDocumentAttributesTitle: notEqualsMap["chaos_document_attributes_title"].(string),
							}
						}

						likeList := conditionMap["like"].(*schema.Set).List()
						if len(likeList) > 0 {
							likeMap := likeList[0].(map[string]interface{})
							clientLike = &client.Like{
								ChaosDocumentAttributesTitle: likeMap["chaos_document_attributes_title"].(string),
							}
						}
					}

					clientConditionGroup = client.ConditionGroup{
						Condition: []client.Condition{
							{
								Equals:     *clientEquals,
								StartsWith: *clientStartsWith,
								NotEquals:  *clientNotEquals,
								Like:       *clientLike,
							},
						},
					}
				}

				permissionList = append(
					permissionList,
					client.Permission{
						Effect:         permissionMap["effect"].(string),
						Actions:        permissionMap["actions"].([]interface{}),
						Resources:      permissionMap["resources"].([]interface{}),
						Version:        permissionMap["version"].(string),
						ConditionGroup: clientConditionGroup,
					},
				)
			}
		}

		var userGroupId string
		if data.Id() != "" {
			userGroupId = data.Id()
		} else if userGroupMap["id"].(string) != "" {
			userGroupId = userGroupMap["id"].(string)
		}

		return &client.CreateUserGroupRequest{
			AuthToken:  authToken,
			ID:         userGroupId,
			Name:       userGroupMap["name"].(string),
			Permission: permissionList,
		}
	}

	return nil
}

func resourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Id())
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

	userGroupContent := CreateUserGroupResponse(resp)
	if err := data.Set("user_groups", userGroupContent); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func CreateUserGroupResponse(resp *client.UserGroup) []map[string]interface{} {
	permissions := make([]interface{}, len(resp.Permissions))
	if resp.Permissions != nil && len(resp.Permissions) > 0 {
		for i, permission := range resp.Permissions {
			permissions[i] = map[string]interface{}{
				"effect":    permission.Effect,
				"actions":   permission.Actions,
				"resources": permission.Resources,
				"version":   permission.Version,
			}
		}
	}

	return []map[string]interface{}{
		{
			"id":          resp.ID,
			"name":        resp.Name,
			"permissions": permissions,
		},
	}
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	resp, err := c.UpdateUserGroup(ctx, GroupObject(data, tokenValue))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.ID)
	return nil
}

func resourceGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	deleteUserGroupRequest := &client.DeleteUserGroupRequest{
		AuthToken: tokenValue,
		ID:        data.Id(),
	}

	if err := c.DeleteUserGroup(ctx, deleteUserGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Id())
	return nil
}

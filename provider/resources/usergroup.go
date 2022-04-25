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

	var id, name string
	var permissionList []client.Permission
	var conditionList []client.Condition

	GroupObjectDTO := GroupObjectDTO{
		ID:             id,
		Name:           name,
		ConditionList:  conditionList,
		PermissionList: permissionList,
	}
	id, name, permissionList = GroupObject(data, GroupObjectDTO)

	createUserGroupRequest := &client.CreateUserGroupRequest{
		AuthToken:  tokenValue,
		ID:         id,
		Name:       name,
		Permission: permissionList,
	}
	resp, err := c.CreateUserGroup(ctx, createUserGroupRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.ID)
	return resourceGroupRead(ctx, data, meta)
}

type GroupObjectDTO struct {
	ID             string
	Name           string
	ConditionList  []client.Condition
	PermissionList []client.Permission
}

func GroupObject(data *schema.ResourceData, dto GroupObjectDTO) (string, string, []client.Permission) {
	if data.Get("user_groups").(*schema.Set).Len() > 0 {
		userGroupInterface := data.Get("user_groups").(*schema.Set).List()[0].(map[string]interface{})
		dto.ID = userGroupInterface["id"].(string)
		dto.Name = userGroupInterface["name"].(string)
		permissions := userGroupInterface["permissions"].(*schema.Set).List()
		if len(permissions) > 0 {
			for _, permissionsElement := range permissions {
				permissionMap := permissionsElement.(map[string]interface{})
				var ConditionGroup client.ConditionGroup
				if len(permissionMap["conditions"].(*schema.Set).List()) > 0 {
					conditionObj := permissionMap["conditions"].(*schema.Set).List()[0]
					conditions := conditionObj.(map[string]interface{})["condition"].(*schema.Set).List()
					conditionMap := conditions[0].(map[string]interface{})
					equalObj := conditionMap["equals"].(*schema.Set).List()[0]
					equal := equalObj.(map[string]interface{})["chaos_document_attributes_title"].(string)
					startsWithObj := conditionMap["starts_with"].(*schema.Set).List()[0]
					startsWith := startsWithObj.(map[string]interface{})["chaos_document_attributes_title"].(string)
					notEqualsObj := conditionMap["not_equals"].(*schema.Set).List()[0]
					notEquals := notEqualsObj.(map[string]interface{})["chaos_document_attributes_title"].(string)
					likeObj := conditionMap["like"].(*schema.Set).List()[0]
					like := likeObj.(map[string]interface{})["chaos_document_attributes_title"].(string)

					equalObject := client.Equals{
						ChaosDocumentAttributesTitle: equal,
					}
					likeObject := client.Like{
						ChaosDocumentAttributesTitle: like,
					}
					notEqualsObject := client.NotEquals{
						ChaosDocumentAttributesTitle: notEquals,
					}
					startsWithObject := client.StartsWith{
						ChaosDocumentAttributesTitle: startsWith,
					}
					dto.ConditionList = append(
						dto.ConditionList,
						client.Condition{
							Equals:     equalObject,
							StartsWith: startsWithObject,
							NotEquals:  notEqualsObject,
							Like:       likeObject,
						})
					ConditionGroup = client.ConditionGroup{
						Condition: dto.ConditionList,
					}
				}

				dto.PermissionList = append(
					dto.PermissionList,
					client.Permission{
						Effect:         permissionMap["effect"].(string),
						Actions:        permissionMap["actions"].([]interface{}),
						Resources:      permissionMap["resources"].([]interface{}),
						Version:        permissionMap["version"].(string),
						ConditionGroup: ConditionGroup,
					})
				dto.ConditionList = nil
			}
		}
	}
	return dto.ID, dto.Name, dto.PermissionList
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
	permissionContent := make(map[string]interface{})
	result := make([]map[string]interface{}, 1)
	userGroupContentMap := make(map[string]interface{})
	if resp.Permissions != nil && len(resp.Permissions) > 0 {
		permissions := make([]interface{}, len(resp.Permissions))
		for i := 0; i < len(resp.Permissions); i++ {
			permissionContent["effect"] = resp.Permissions[i].Effect
			permissionContent["actions"] = resp.Permissions[i].Actions
			permissionContent["resources"] = resp.Permissions[i].Resources
			permissionContent["version"] = resp.Permissions[i].Version
			permissions[i] = permissionContent
		}
		userGroupContentMap["permissions"] = permissions
	}
	userGroupContentMap["id"] = resp.ID
	userGroupContentMap["name"] = resp.Name
	result[0] = userGroupContentMap

	return result
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	var id, name string
	var permissionList []client.Permission
	var conditionList []client.Condition

	GroupObjectDTO := GroupObjectDTO{
		ID:             id,
		Name:           name,
		ConditionList:  conditionList,
		PermissionList: permissionList,
	}
	id, name, permissionList = GroupObject(data, GroupObjectDTO)

	createUserGroupRequest := &client.CreateUserGroupRequest{
		AuthToken:  tokenValue,
		ID:         id,
		Name:       name,
		Permission: permissionList,
	}

	resp, err := c.UpdateUserGroup(ctx, createUserGroupRequest)
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

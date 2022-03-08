package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroup() *schema.Resource {
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
				Optional: false,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: false,
							ForceNew: false,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: false,
							ForceNew: false,
							Optional: true,
						},
						"permissions": {
							Type:     schema.TypeSet,
							Optional: true,
							Required: false,
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
										Required: false,
										ForceNew: false,
										Optional: true,
										//Computed: true,
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
										Required: false,
										ForceNew: false,
										Optional: true,
										//Computed: true,
									},
									"conditions": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     schema.TypeSet,
													Optional: true,
													Required: false,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"starts_with": {
																Type:     schema.TypeSet,
																Optional: true,
																ForceNew: false,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Required: false,
																			ForceNew: false,
																			Optional: true,
																		},
																	},
																},
															},
															"equals": {
																Type:     schema.TypeSet,
																Optional: true,
																ForceNew: false,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Required: false,
																			ForceNew: false,
																			Optional: true,
																		},
																	},
																},
															},
															"not_equals": {
																Type:     schema.TypeSet,
																Optional: true,
																ForceNew: false,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Required: false,
																			ForceNew: false,
																			Optional: true,
																		},
																	},
																},
															},
															"like": {
																Type:     schema.TypeSet,
																Optional: true,
																ForceNew: false,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"chaos_document_attributes_title": {
																			Type:     schema.TypeString,
																			Required: false,
																			ForceNew: false,
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
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token

	var id, name string
	var permissionList []client.Permission
	var conditionList []client.Condition
	var actionsList []interface{}
	var resourcesList []interface{}

	id, name, permissionList = CreateUserGroupObject(data, id, name, conditionList, actionsList, resourcesList,
		permissionList)

	createUserGroupRequest := &client.CreateUserGroupRequest{
		AuthToken:  tokenValue,
		Id:         id,
		Name:       name,
		Permission: permissionList,
	}
	resp, err := c.CreateUserGroup(ctx, createUserGroupRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Id)
	return resourceGroupRead(ctx, data, meta)
}

func CreateUserGroupObject(data *schema.ResourceData, id string, name string, conditionList []client.Condition, actionsList []interface{}, resourcesList []interface{}, permissionList []client.Permission) (string, string, []client.Permission) {
	if data.Get("user_groups").(*schema.Set).Len() > 0 {
		userGroupInterface := data.Get("user_groups").(*schema.Set).List()[0].(map[string]interface{})
		id = userGroupInterface["id"].(string)
		name = userGroupInterface["name"].(string)
		permissions := userGroupInterface["permissions"].(*schema.Set).List()
		if len(permissions) > 0 {
			for _, permissionsElement := range permissions {
				permissionMap := permissionsElement.(map[string]interface{})
				var ConditionGroup client.ConditionGroup
				if len(permissionMap["conditions"].(*schema.Set).List()) > 0 {
					conditions := permissionMap["conditions"].(*schema.Set).List()[0].(map[string]interface{})["condition"].(*schema.Set).List()
					conditionMap := conditions[0].(map[string]interface{})
					equal := conditionMap["equals"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)
					startsWith := conditionMap["starts_with"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)
					notEquals := conditionMap["not_equals"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)
					like := conditionMap["like"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)

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
					conditionList = append(
						conditionList,
						client.Condition{
							Equals:     equalObject,
							StartsWith: startsWithObject,
							NotEquals:  notEqualsObject,
							Like:       likeObject,
						})
					ConditionGroup = client.ConditionGroup{
						Condition: conditionList,
					}
				}

				actionsList = permissionMap["actions"].([]interface{})
				resourcesList = permissionMap["resources"].([]interface{})
				permissionList = append(
					permissionList,
					client.Permission{
						Effect:         permissionMap["effect"].(string),
						Actions:        actionsList,
						Resources:      resourcesList,
						Version:        permissionMap["version"].(string),
						ConditionGroup: ConditionGroup,
					})
				conditionList = nil
			}
		}
	}
	return id, name, permissionList
}

func resourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Id())
	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
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

func CreateUserGroupResponse(resp *client.Group) []map[string]interface{} {
	userGroupContent := make([]map[string]interface{}, 1)
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
	userGroupContentMap["id"] = resp.Id
	userGroupContentMap["name"] = resp.Name
	result[0] = userGroupContentMap
	userGroupContent = result
	return userGroupContent
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	var id, name string
	var permissionList []client.Permission
	var conditionList []client.Condition
	var actionsList []interface{}
	var resourcesList []interface{}

	id, name, permissionList = CreateUserGroupObject(data, id, name, conditionList, actionsList, resourcesList, permissionList)

	createUserGroupRequest := &client.CreateUserGroupRequest{
		AuthToken:  tokenValue,
		Id:         id,
		Name:       name,
		Permission: permissionList,
	}

	resp, err := c.UpdateUserGroup(ctx, createUserGroupRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Id)
	return nil
}

func resourceGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient

	tokenValue := meta.(*ProviderMeta).token
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

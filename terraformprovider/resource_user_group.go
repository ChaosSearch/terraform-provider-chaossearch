package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
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
							Computed: true,
						},

						"permissions": {
							Type:     schema.TypeSet,
							Optional: false,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"permission": {
										Type: schema.TypeSet,
										//Required: true,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"effect": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"resources": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"conditions": {
													Type:     schema.TypeSet,
													Optional: true,
													ForceNew: false,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"condition": {
																Type:     schema.TypeSet,
																Optional: false,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"starts_with": {
																			Type: schema.TypeSet,
																			//Required: true,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"equals": {
																			Type: schema.TypeSet,
																			//Required: true,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"not_equals": {
																			Type: schema.TypeSet,
																			//Required: true,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"like": {
																			Type: schema.TypeSet,
																			//Required: true,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
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
				},
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("creating groups")
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	log.Debug("token value------------>>>>", tokenValue)
	set := data.Get("user_groups")
	log.Info("user_groups===>", set)

	if data.Get("user_groups").(*schema.Set).Len() > 0 {
		policiesInterfaces := data.Get("user_groups").(*schema.Set).List()[0]
		policiesInterface := policiesInterfaces.(map[string]interface{})

		//get permission map as a list
		objectList := policiesInterface["permissions"].(*schema.Set).List()
		permissionObject := objectList[0]
		log.Info("objectList====>", objectList)
		permission := permissionObject.(map[string]interface{})
		i := permission["permission"].(*schema.Set).List()
		//log.Debug("index", index)
		for index1, element1 := range i {
			//get permission map one by one
			permissionMap := element1.(map[string]interface{})
			log.Info("policy1==>", permissionMap)
			log.Debug("action====>", permissionMap["action"].(string))
			log.Debug("resources====>", permissionMap["resources"].(string))
			log.Debug("effect====>", permissionMap["effect"].(string))
			log.Debug("conditions====>", permissionMap["conditions"].(*schema.Set).List()[index1])
			log.Debug(index1, "index")
		}
		//}
	}

	createUserGroupRequest := &client.CreateUserGroupRequest{
		AuthToken: tokenValue,
		Id:        data.Get("id").(string),
		Name:      data.Get("name").(string),
	}

	log.Debug("createUserGroupRequest.id-->", createUserGroupRequest.Id)
	log.Debug("createUserGroupRequest.name-->", createUserGroupRequest.Name)

	if err := c.CreateUserGroup(ctx, createUserGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))

	return nil
}

func resourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("reading groups")
	// c := meta.(*ProviderMeta).Client
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value------------>>>>", tokenValue)
	//TODO to be developed
	return nil
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("updating groups")
	// c := meta.(*ProviderMeta).Client
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value------------>>>>", tokenValue)
	//TODO to be developed
	return nil
}

func resourceGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("deleting groups")
	// c := meta.(*ProviderMeta).Client
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value------------>>>>", tokenValue)
	//TODO to be developed
	return nil
}

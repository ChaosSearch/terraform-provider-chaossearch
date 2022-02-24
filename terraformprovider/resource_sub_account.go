package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubAccountCreate,
		ReadContext:   resourceSubAccountRead,
		UpdateContext: resourceSubAccountUpdate,
		DeleteContext: resourceSubAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_info_block": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"full_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"hocon": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSubAccountCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient
	formatColumnSelectionInterface := data.Get("user_info_block").(*schema.Set).List()[0].(map[string]interface{})

	userInfoBlock := client.UserInfoBlock{
		Username: formatColumnSelectionInterface["username"].(string),
		FullName: formatColumnSelectionInterface["full_name"].(string),
		Email:    formatColumnSelectionInterface["email"].(string),
	}

	createSubAccountRequest := &client.CreateSubAccountRequest{
		AuthToken:     meta.(*ProviderMeta).token,
		UserInfoBlock: userInfoBlock,
		GroupIds:      data.Get("group_ids").([]interface{}),
		Password:      data.Get("password").(string),
		HoCon:         data.Get("hocon").([]interface{}),
	}

	if err := c.CreateSubAccount(ctx, createSubAccountRequest); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(formatColumnSelectionInterface["username"].(string))
	return nil
}

func resourceSubAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSubAccountUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSubAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

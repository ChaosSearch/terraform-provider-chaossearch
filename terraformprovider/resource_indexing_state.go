package main

import (
	"context"
	"cs-tf-provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIndexingState() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIndexingStateCreate,
		ReadContext:   resourceIndexingStateRead,
		UpdateContext: resourceIndexingStateUpdate,
		DeleteContext: resourceIndexingStateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"object_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the live indexing should be running or not",
				Required:    true,
				ForceNew:    false,
			},
		},
	}
}

func resourceIndexingStateCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).Client

	updateIndexingStateRequest := &client.UpdateIndexingStateRequest{
		ObjectGroupName: data.Get("object_group_name").(string),
		Active:          data.Get("active").(bool),
	}

	if err := c.UpdateIndexingState(ctx, updateIndexingStateRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("object_group_name").(string))

	return resourceIndexingStateRead(ctx, data, meta)
}

func resourceIndexingStateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}

	c := meta.(*ProviderMeta).Client

	readIndexingStateRequest := &client.ReadIndexingStateRequest{
		ObjectGroupName: data.Get("object_group_name").(string),
	}

	resp, err := c.ReadIndexingState(ctx, readIndexingStateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.ObjectGroupName)
	data.Set("active", resp.Active)

	return diags
}

func resourceIndexingStateUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).Client

	updateIndexingStateRequest := &client.UpdateIndexingStateRequest{
		ObjectGroupName: data.Get("object_group_name").(string),
		Active:          data.Get("active").(bool),
	}
	if err := c.UpdateIndexingState(ctx, updateIndexingStateRequest); err != nil {
		return diag.FromErr(err)
	}

	return resourceIndexingStateRead(ctx, data, meta)
}

func resourceIndexingStateDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).Client

	stopIndexingRequest := &client.UpdateIndexingStateRequest{
		ObjectGroupName: data.Get("object_group_name").(string),
		Active:          false,
	}
	if err := c.UpdateIndexingState(ctx, stopIndexingRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return nil
}

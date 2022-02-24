package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	log "github.com/sirupsen/logrus"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBucketImportOrHide,
		UpdateContext: resourceBucketImportOrHide,
		ReadContext:   resourceBucketRead,
		DeleteContext: resourceBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hide_bucket": {
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceBucketRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//TODO Import or Hide bucket read to be implemented
	return nil
}
func resourceBucketDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//TODO Import or Hide bucket delete or revert to be implemented
	return nil
}
func resourceBucketImportOrHide(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("Importing / hiding bucket")
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token

	importBucketRequest := &client.ImportBucketRequest{
		AuthToken:  tokenValue,
		Bucket:     data.Get("bucket").(string),
		HideBucket: data.Get("hide_bucket").(bool),
	}

	if err := c.ImportBucket(ctx, importBucketRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))

	return nil
}

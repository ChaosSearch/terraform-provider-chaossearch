package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIndexModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceIndexModel,
		ReadContext:   readResourceIndexModel,
		DeleteContext: deleteResourceIndexModel,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: true,
				Optional: true,
			},
			"model_mode": {
				Type:     schema.TypeInt,
				Required: false,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func createResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	indexModel := &client.IndexModelRequest{
		AuthToken:  meta.(*models.ProviderMeta).Token,
		BucketName: data.Get("bucket_name").(string),
		ModelMode:  data.Get("model_mode").(int),
	}
	resp, err := c.CreateIndexModel(ctx, indexModel)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strings.Join([]string{"BucketName:" + resp.BucketName, " Result: " + strconv.FormatBool(resp.Result)}, ","))
	return nil
}

func readResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func deleteResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

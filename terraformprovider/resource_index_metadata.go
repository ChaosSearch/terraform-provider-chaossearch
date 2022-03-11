package cs

import (
	"context"
	"cs-tf-provider/client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceIndexMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceIndexMetadata,
		ReadContext:   readResourceIndexMetadata,
		DeleteContext: deleteResourceIndexMetadata,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"bucket_names": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createResourceIndexMetadata(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient
	indexModel := &client.IndexMetadataRequest{
		AuthToken:   meta.(*ProviderMeta).token,
		BucketNames: data.Get("bucket_names").(string),
	}
	resp, err := c.ReadIndexMetadata(ctx, indexModel)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strings.Join([]string{"Bucket:" + resp.Bucket,
		" LastIndexTime: " + fmt.Sprint(resp.LastIndexTime), " State: " + resp.State}, ","))
	return nil
}

func readResourceIndexMetadata(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func deleteResourceIndexMetadata(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

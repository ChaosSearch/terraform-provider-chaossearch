package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"fmt"
	"strconv"
	"time"

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
				ForceNew: true,
				Required: true,
			},
			"model_mode": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"delete_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
		},
	}
}

func createResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var bucketName string
	var modelMode int
	c := meta.(*models.ProviderMeta).CSClient
	if data.Get("bucket_name") != nil {
		bucketName = data.Get("bucket_name").(string)
	}

	if data.Get("model_mode") != nil {
		modelMode = data.Get("model_mode").(int)
	}

	indexModel := &client.IndexModelRequest{
		AuthToken:  meta.(*models.ProviderMeta).Token,
		BucketName: bucketName,
		ModelMode:  modelMode,
	}

	resp, err := c.CreateIndexModel(ctx, indexModel)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("BucketName: %s, Result: %s", resp.BucketName, strconv.FormatBool(resp.Result)))
	return nil
}

func readResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func deleteResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deleteEnabled := data.Get("delete_enabled").(bool)
	if deleteEnabled {
		var listBucketResp *client.ListBucketResponse
		c := meta.(*models.ProviderMeta).CSClient
		authToken := meta.(*models.ProviderMeta).Token
		bucketName := data.Get("bucket_name").(string)
		listBucketResp, err := c.ReadIndexModel(
			ctx,
			bucketName,
			authToken,
		)

		if err != nil {
			return diag.FromErr(err)
		}

		if listBucketResp.Contents != nil {
			err = c.DeleteIndexModel(ctx, listBucketResp.Contents.Key, authToken)
			if err != nil {
				return diag.FromErr(err)
			}

			// await return until index confirmed deletion
			tickerCounter := 0
			for tickerCounter <= 300 {
				listBucketResp, err = c.ReadIndexModel(ctx, bucketName, authToken)
				if listBucketResp.Contents == nil {
					break
				} else if listBucketResp.Contents.Key == "" {
					break
				}

				time.Sleep(15 * time.Second)
				tickerCounter += 15
				if tickerCounter >= 300 {
					err = fmt.Errorf(`
						Failure confirming index deletion => Timeout (5 Minutes)
						Note:
							This does not mean there was a failure with index deletion.
							Please confirm the state of the index within ChaosSearch.
					`)
				}
			}

			if err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	} else {
		return diag.Errorf(`
			Failure deleting index model => Deletion is not enabled
			WARNING: 
				Enabling will allow for all index data within an Object Group to be deleted, 
				default is set to false as a safeguard.
				If you're sure, you can set 'delete_enabled' to true.
				This is also persisted in your .tfstate
			Note: 
				Index data existing will block Object Group deletion.
		`)
	}
}

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
				Optional: true,
			},
			"model_mode": {
				Type:     schema.TypeInt,
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

	data.SetId(fmt.Sprintf("BucketName: %s, Result: %s", resp.BucketName, strconv.FormatBool(resp.Result)))
	return nil
}

func readResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func deleteResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		quit := make(chan bool)
		ticker := time.NewTicker(15 * time.Second)
		go func() {
			<-time.After(5 * time.Minute)
			close(quit)
			err = fmt.Errorf("Failure confirming index deletion => Timeout (5 Minutes)")
		}()

		func() {
			for {
				select {
				case <-ticker.C:
					listBucketResp, err = c.ReadIndexModel(ctx, bucketName, authToken)
					if listBucketResp.Contents == nil {
						close(quit)
					} else if listBucketResp.Contents.Key == "" {
						close(quit)
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

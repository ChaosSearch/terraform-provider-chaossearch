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

const (
	BucketName     = "bucket_name"
	ModelMode      = "model_mode"
	Result         = "result"
	Indexed        = "indexed"
	Options        = "options"
	DeleteEnabled  = "delete_enabled"
	DeleteTimeout  = "delete_timeout"
	SkipIndexPause = "skip_index_pause"
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
			BucketName: {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			ModelMode: {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			Result: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			Indexed: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			Options: {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DeleteEnabled: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						DeleteTimeout: {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							ForceNew:    true,
							Description: "The amount of time before a delete request times out, in seconds",
						},
						SkipIndexPause: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable if you don't want to wait for indexing to complete",
						},
					},
				},
			},
		},
	}
}

func createResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var skipIndexPause bool
	options := data.Get(Options).(*schema.Set).List()
	if len(options) > 0 {
		optionsMap := options[0].(map[string]interface{})
		skipIndexPause = optionsMap[SkipIndexPause].(bool)
	}

	var bucketName string
	var modelMode int
	c := meta.(*models.ProviderMeta).CSClient
	if data.Get(BucketName) != nil {
		bucketName = data.Get(BucketName).(string)
	}

	if data.Get(ModelMode) != nil {
		modelMode = data.Get(ModelMode).(int)
	}

	authToken := meta.(*models.ProviderMeta).Token
	indexModel := &client.IndexModelRequest{
		AuthToken:  authToken,
		BucketName: bucketName,
		ModelMode:  modelMode,
	}

	resp, err := c.CreateIndexModel(ctx, indexModel)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("BucketName: %s, Result: %s", resp.BucketName, strconv.FormatBool(resp.Result)))
	err = data.Set(BucketName, resp.BucketName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set(Result, resp.Result)
	if err != nil {
		return diag.FromErr(err)
	}

	// Confirm index status before ending create
	if !skipIndexPause {
		indexed := false
		for !indexed {
			checkResp, err := c.CheckIndexModel(ctx, bucketName, authToken)
			if err != nil {
				return diag.FromErr(err)
			}

			indexed = checkResp.Indexed
			time.Sleep(15 * time.Second)
		}

		err = data.Set(Indexed, indexed)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = data.Set(Indexed, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func readResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func deleteResourceIndexModel(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var deleteEnabled bool
	var deleteTimeout int
	var bucketName string
	c := meta.(*models.ProviderMeta).CSClient

	options := data.Get(Options).(*schema.Set).List()
	if len(options) > 0 {
		optionsMap := options[0].(map[string]interface{})
		deleteEnabled = optionsMap[DeleteEnabled].(bool)
		deleteTimeout = optionsMap[DeleteTimeout].(int)
	}

	authToken := meta.(*models.ProviderMeta).Token
	if data.Get(BucketName) != nil {
		bucketName = data.Get(BucketName).(string)
	}

	indexModel := &client.IndexModelRequest{
		AuthToken:  authToken,
		BucketName: bucketName,
		ModelMode:  -1,
	}

	_, err := c.CreateIndexModel(ctx, indexModel)
	if err != nil {
		return diag.FromErr(err)
	}

	if deleteEnabled {
		var listBucketResp *client.ListBucketResponse
		authToken := meta.(*models.ProviderMeta).Token
		bucketName := data.Get(BucketName).(string)
		listBucketResp, err := c.ReadIndexModel(
			ctx,
			bucketName,
			authToken,
		)

		if err != nil {
			return diag.FromErr(err)
		}

		if listBucketResp.Contents != nil {
			for _, content := range *listBucketResp.Contents {
				err = c.DeleteIndexModel(ctx, content.Key, authToken)
				if err != nil {
					return diag.FromErr(err)
				}
			}

			// await return until index confirmed deletion
			timeout := false
			tickerCounter := 0
			for !timeout {
				listBucketResp, err = c.ReadIndexModel(ctx, bucketName, authToken)
				if err != nil {
					return diag.FromErr(err)
				}

				contentEmpty := true
				for _, content := range *listBucketResp.Contents {
					if content.Key != "" {
						contentEmpty = false
					}
				}

				if contentEmpty {
					break
				}

				time.Sleep(15 * time.Second)
				tickerCounter += 15
				if deleteTimeout != 0 && tickerCounter >= deleteTimeout {
					err = fmt.Errorf(`
						Failure confirming index deletion => Timeout (%v Seconds)
						Note:
							This does not mean there was a failure with index deletion.
							Please confirm the state of the index within ChaosSearch.
					`, deleteTimeout)

					timeout = true
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

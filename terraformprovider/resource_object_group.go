package main

import (
	"context"
	"cs-tf-provider/client"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceObjectGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectGroupCreate,
		ReadContext:   resourceObjectGroupRead,
		UpdateContext: resourceObjectGroupUpdate,
		DeleteContext: resourceObjectGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON", "LOG"}, false),
			},
			"filter_json": {
				Type:         schema.TypeString,
				Default:      `[{"field":"key","regex":".*"}]`,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"compression": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"", "gzip", "zlib", "snappy"}, true),
			},
			"live_events_sqs_arn": {
				Type:         schema.TypeString,
				Default:      "",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^arn:aws:sqs:.*`), `must be an SQS ARN: "arn:aws:sqs:..."`),
			},
			"partition_by": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
				ForceNew: true,
			},
			"pattern": {
				Type:     schema.TypeString,
				Default:  ".*",
				Optional: true,
				ForceNew: true,
			},
			"index_retention": {
				Type:        schema.TypeInt,
				Default:     14,
				Description: "Number of days to keep the data before deleting it",
				Optional:    true,
				ForceNew:    false,
			},
			"array_flatten_depth": {
				Type:         schema.TypeInt,
				Default:      0,
				Description:  "Array flattening level. 0 - disabled, -1 - unlimited, >1 - the respective flattening level",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(-1),
			},
			"column_renames": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Default:     make(map[string]string),
				ForceNew:    true,
				Description: "A map specifying names of columns to rename (keys) and what to rename them to (values)",
			},
			"keep_original": {
				Type:        schema.TypeBool,
				Default:     false,
				Description: "Works in connection with the `column_selection`, dictates whether to keep the fields filtered out by the column_selection in a separate field",
				Optional:    true,
				ForceNew:    true,
			},
			"column_selection": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"includes": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
							ForceNew: true,
						},
					},
				},
				Optional:    true,
				ForceNew:    true,
				Description: "List of fields in logs to include or exclude from parsing. If nothing is specified, all fields will be parsed",
			},

			// Workaround. Otherwise Terraform fails with "All fields are ForceNew or Computed w/out Optional, Update is superfluous"
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceObjectGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).Client

	// "unlimited" flattening represented as "null" in the api, and as -1 in the terraform module
	// because the terraform sdk doesn't support nil values in configs https://github.com/hashicorp/terraform-plugin-sdk/issues/261
	// We represent "null" as an int pointer to nil in the code.
	arrayFlattenTF := data.Get("array_flatten_depth").(int)
	var arrayFlattenCS *int
	if arrayFlattenTF == -1 {
		// -1 in terraform represents "null" in the ChaosSearch API call
		arrayFlattenCS = nil
	} else {
		// any other value is passed as is
		arrayFlattenCS = &arrayFlattenTF
	}
	var columnSelection map[string]interface{}

	if data.Get("column_selection").(*schema.Set).Len() > 0 {
		columnSelectionInterfaces := data.Get("column_selection").(*schema.Set).List()[0]
		columnSelectionInterface := columnSelectionInterfaces.(map[string]interface{})

		includesListOfInterfaces := columnSelectionInterface["includes"].(*schema.Set).List()

		columnSelection = map[string]interface{}{
			"type":     columnSelectionInterface["type"].(string),
			"includes": includesListOfInterfaces,
		}
	}

	createObjectGroupRequest := &client.CreateObjectGroupRequest{
		Name:              data.Get("name").(string),
		SourceBucket:      data.Get("source_bucket").(string),
		FilterJSON:        data.Get("filter_json").(string),
		Format:            data.Get("format").(string),
		Compression:       data.Get("compression").(string),
		LiveEventsSqsArn:  data.Get("live_events_sqs_arn").(string),
		PartitionBy:       data.Get("partition_by").(string),
		Pattern:           data.Get("pattern").(string),
		IndexRetention:    data.Get("index_retention").(int),
		KeepOriginal:      data.Get("keep_original").(bool),
		ArrayFlattenDepth: arrayFlattenCS,
		ColumnRenames:     data.Get("column_renames").(map[string]interface{}),
		ColumnSelection:   columnSelection,
	}

	if err := c.CreateObjectGroup(ctx, createObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("name").(string))

	return resourceObjectGroupRead(ctx, data, meta)
}

func resourceObjectGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).Client

	req := &client.ReadObjectGroupRequest{
		ID: data.Id(),
	}

	resp, err := c.ReadObjectGroup(ctx, req)
	if err != nil {
		return diag.Errorf("Failed to read object group: %s", err)
	}

	data.Set("name", data.Id())
	data.Set("filter_json", resp.FilterJSON)
	data.Set("format", resp.Format)
	data.Set("live_events_sqs_arn", resp.LiveEventsSqsArn)
	data.Set("index_retention", resp.IndexRetention)

	// When the object in an Object Group use no compression, you need to create it with
	// `compression = ""`. However, when querying an Object Group whose object are not
	// compressed, the API returns `compression = "none"`. We coerce the "none" value to
	// an empty string in order not to confuse Terraform.
	compressionOrEmptyString := resp.Compression
	if strings.ToLower(compressionOrEmptyString) == "none" {
		compressionOrEmptyString = ""
	}
	data.Set("compression", compressionOrEmptyString)

	data.Set("partition_by", resp.PartitionBy)
	data.Set("pattern", resp.Pattern)
	data.Set("source_bucket", resp.SourceBucket)

	data.Set("column_selection", resp.ColumnSelection)

	// "unlimited" flattening represented as "null" in the api, and as -1 in the terraform module
	// because the terraform sdk doesn't support nil values in configs https://github.com/hashicorp/terraform-plugin-sdk/issues/261
	// We represent "null" as an int pointer to nil in the code.
	if resp.ArrayFlattenDepth == nil {
		data.Set("array_flatten_depth", -1)
	} else {
		data.Set("array_flatten_depth", resp.ArrayFlattenDepth)
	}

	return diags
}

func resourceObjectGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).Client

	updateObjectGroupRequest := &client.UpdateObjectGroupRequest{
		Name:           data.Get("name").(string),
		IndexRetention: data.Get("index_retention").(int),
	}

	if err := c.UpdateObjectGroup(ctx, updateObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	return resourceObjectGroupRead(ctx, data, meta)
}

func resourceObjectGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).Client

	deleteObjectGroupRequest := &client.DeleteObjectGroupRequest{
		Name: data.Get("name").(string),
	}

	if err := c.DeleteObjectGroup(ctx, deleteObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

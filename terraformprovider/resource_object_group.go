package main

import (
	"context"
	"cs-tf-provider/client"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceObjectGroup() *schema.Resource {
	log.Info("called resourceObjectGroup")
	return &schema.Resource{
		CreateContext: resourceObjectGroupCreate,
		ReadContext:   resourceObjectGroupRead,
		UpdateContext: resourceObjectGroupUpdate,
		DeleteContext: resourceObjectGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"format": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"column_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"header_row": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"row_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"index_retention": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"for_partition": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"overall": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"regex_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"regex": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"interval": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"mode": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_irregular": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"realtime": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
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

	formatColumnSelectionInterface := data.Get("format").(*schema.Set).List()[0].(map[string]interface{})
	intervalColumnSelectionInterface := data.Get("interval").(*schema.Set).List()[0].(map[string]interface{})
	indexRetentionColumnSelectionInterface := data.Get("index_retention").(*schema.Set).List()[0].(map[string]interface{})
	optionsColumnSelectionInterface := data.Get("options").(*schema.Set).List()[0].(map[string]interface{})
	//filterColumnSelectionInterface := data.Get("filter").(*schema.Set).List()[0].(map[string]interface{})

	format := client.Format{
		Type:            formatColumnSelectionInterface["_type"].(string),
		ColumnDelimiter: formatColumnSelectionInterface["column_delimiter"].(string),
		RowDelimiter:    formatColumnSelectionInterface["row_delimiter"].(string),
		HeaderRow:       formatColumnSelectionInterface["header_row"].(bool),
	}

	interval := client.Interval{
		Mode:   intervalColumnSelectionInterface["mode"].(int),
		Column: intervalColumnSelectionInterface["column"].(int),
	}

	indexRetention := client.IndexRetention{
		ForPartition: indexRetentionColumnSelectionInterface["for_partition"].([]interface{}),
		Overall:      indexRetentionColumnSelectionInterface["overall"].(int),
	}

	options := client.Options{
		IgnoreIrregular: optionsColumnSelectionInterface["ignore_irregular"].(bool),
	}

	var fieldOne string
	var prefix string
	var fieldTwo string
	var regeX string

	if data.Get("filter").(*schema.Set).Len() > 0 {
		filterColumnSelectionInterface := data.Get("filter").(*schema.Set).List()[0]
		filterColumnSelection := filterColumnSelectionInterface.(map[string]interface{})

		prefixFilter := filterColumnSelection["prefix_filter"].(*schema.Set).List()[0].(map[string]interface{})
		regexFilter := filterColumnSelection["regex_filter"].(*schema.Set).List()[0].(map[string]interface{})

		fieldOne = prefixFilter["field"].(string)
		prefix = prefixFilter["prefix"].(string)

		fieldTwo = regexFilter["field"].(string)
		regeX = regexFilter["regex"].(string)
	}

	classOne := client.PrefixFilter{
		Field:  fieldOne,
		Prefix: prefix,
	}

	classTwo := client.RegexFilter{
		Field: fieldTwo,
		Regex: regeX,
	}
	filter := &client.Filter{
		&classOne,
		&classTwo,
	}

	createObjectGroupRequest := &client.CreateObjectGroupRequest{
		Bucket:         data.Get("bucket").(string),
		Source:         data.Get("source").(string),
		Format:         &format,
		Interval:       &interval,
		IndexRetention: &indexRetention,
		Filter:         filter,
		Options:        &options,
		Realtime:       data.Get("realtime").(bool),
	}

	if err := c.CreateObjectGroup(ctx, createObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(data.Get("bucket").(string))
	return nil
}

func resourceObjectGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Info("called READ")

	log.Info("11111111111111111111")
	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).Client

	log.Info("22222222222222")

	req := &client.ReadObjectGroupRequest{
		ID: data.Id(),
	}

	log.Info("33333333333333")

	log.Warn("req---->", req)
	resp, err := c.ReadObjectGroup(ctx, req)
	log.Info("4444444444444")

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

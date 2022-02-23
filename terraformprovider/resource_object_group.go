package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
	"strings"
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
			"bucket": {
				Type: schema.TypeString,
				//Required: true,
				Optional: true,
				ForceNew: false,
			},
			"source": {
				Type: schema.TypeString,
				//Required: true,
				Optional: true,
				ForceNew: false,
			},
			"format": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"column_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"header_row": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"row_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"index_retention": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
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
							ForceNew: false,
						},
					},
				},
			},
			"filter": {
				Type: schema.TypeSet,
				//Required: true,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: false,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: false,
									},
								},
							},
						},
						"regex_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: false,
									},
									"regex": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: false,
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
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"mode": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_irregular": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
			"realtime": {
				Type: schema.TypeBool,
				//Required:    true,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"index_parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"index_retention_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"target_active_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
			"live_events_parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
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
	c := meta.(*ProviderMeta).CSClient

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

	var prefixFilterField string
	var prefix string
	var regexFilterField string
	var regeX string

	if data.Get("filter").(*schema.Set).Len() > 0 {
		filterColumnSelectionInterface := data.Get("filter").(*schema.Set).List()[0]
		filterColumnSelection := filterColumnSelectionInterface.(map[string]interface{})

		prefixFilter := filterColumnSelection["prefix_filter"].(*schema.Set).List()[0].(map[string]interface{})
		regexFilter := filterColumnSelection["regex_filter"].(*schema.Set).List()[0].(map[string]interface{})

		prefixFilterField = prefixFilter["field"].(string)
		prefix = prefixFilter["prefix"].(string)

		regexFilterField = regexFilter["field"].(string)
		regeX = regexFilter["regex"].(string)
	}

	prefixFilter := client.PrefixFilter{
		Field:  prefixFilterField,
		Prefix: prefix,
	}

	regexFilter := client.RegexFilter{
		Field: regexFilterField,
		Regex: regeX,
	}
	filter := &client.Filter{
		PrefixFilter: &prefixFilter,
		RegexFilter:  &regexFilter,
	}
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value------------>>>>", tokenValue)
	createObjectGroupRequest := &client.CreateObjectGroupRequest{
		AuthToken:      tokenValue,
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
	return resourceObjectGroupRead(ctx, data, meta)
	//return nil
}

func resourceObjectGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).CSClient

	data.SetId(data.Get("bucket").(string))

	//req := &client.ReadObjectGroupRequest{
	//	ID: bucketId,
	//}
	tokenValue := meta.(*ProviderMeta).token
	req := &client.ReadObjectGroupRequest{
		ID:        data.Id(),
		AuthToken: tokenValue,
	}

	resp, err := c.ReadObjectGroup(ctx, req)

	if resp == nil {
		return diag.Errorf("Couldn't find object group: %s", err)
	}

	data.Set("object_group_id", resp.ID)
	if err != nil {
		return diag.Errorf("Failed to read object group: %s", err)
	}

	//do not change working fine
	data.Set("name", data.Id())
	data.Set("_public", resp.Public)
	data.Set("_type", resp.Type)
	data.Set("content_type", resp.ContentType)
	data.Set("_realtime", resp.Realtime)
	data.Set("bucket", resp.Bucket)

	var prefixFilter interface{} = resp.ObjectFilter.And[0]
	var prefixFilterResponse = map[string]string{}
	for k, v := range prefixFilter.(map[string]interface{}) {
		prefixFilterResponse[k] = v.(string)
	}
	prefixFilterField := prefixFilterResponse["field"]
	prefixFilterPrefix := prefixFilterResponse["prefix"]
	var regexFilter interface{} = resp.ObjectFilter.And[1]
	var regexFilterRes = map[string]string{}
	for k, v := range regexFilter.(map[string]interface{}) {
		regexFilterRes[k] = v.(string)
	}
	regexFilterField := regexFilterRes["field"]
	regexFilterRegex := regexFilterRes["regex"]

	PrefixFilterObjectMap := make(map[string]interface{})
	PrefixFilterObjectMap["field"] = prefixFilterField
	PrefixFilterObjectMap["prefix"] = prefixFilterPrefix

	RegexFilterObjectMap := make(map[string]interface{})
	RegexFilterObjectMap["field"] = regexFilterField
	RegexFilterObjectMap["regex"] = regexFilterRegex

	filter := make([]interface{}, 1)
	PrefixFilter := make([]interface{}, 1)
	PrefixFilter[0] = PrefixFilterObjectMap
	RegexFilter := make([]interface{}, 1)
	RegexFilter[0] = RegexFilterObjectMap
	filterObjectMap := make(map[string]interface{})
	filterObjectMap["prefix_filter"] = PrefixFilter
	filterObjectMap["regex_filter"] = RegexFilter
	filter[0] = filterObjectMap
	data.Set("filter", filter)

	if resp.Format != nil {
		format := make([]interface{}, 1)
		formatObjectMap := make(map[string]interface{})
		formatObjectMap["_type"] = resp.Format.Type
		formatObjectMap["header_row"] = resp.Format.HeaderRow
		formatObjectMap["column_delimiter"] = resp.Format.ColumnDelimiter
		formatObjectMap["row_delimiter"] = resp.Format.RowDelimiter
		format[0] = formatObjectMap
		data.Set("format", format)
		log.Info("format====>", format)
	}

	if resp.Interval != nil {
		interval := make([]interface{}, 1)
		intervalObjectMap := make(map[string]interface{})
		intervalObjectMap["column"] = resp.Interval.Mode
		intervalObjectMap["mode"] = resp.Interval.Column
		interval[0] = intervalObjectMap
		data.Set("interval", interval)
	}

	if resp.Metadata != nil {
		metadata := make([]interface{}, 1)
		metadataObjectMap := make(map[string]interface{})
		metadataObjectMap["creation_date"] = resp.Metadata.CreationDate
		metadata[0] = metadataObjectMap
		data.Set("metadata", metadata)
	}

	if resp.Options != nil {
		options := make([]interface{}, 1)
		optionsObjectMap := make(map[string]interface{})
		optionsObjectMap["ignore_irregular"] = resp.Options.IgnoreIrregular
		options[0] = optionsObjectMap
		data.Set("options", options)
	}

	RegionAvailability := make([]interface{}, 1)
	RegionAvailability[0] = resp.RegionAvailability
	data.Set("region_availability", RegionAvailability[0])

	data.Set("source", resp.Source)
	data.Set("filter_json", resp.FilterJSON)
	data.Set("live_events_sqs_arn", resp.LiveEventsSqsArn)

	//data.Set("index_retention", resp.IndexRetention)
	data.Set("partition_by", resp.PartitionBy)
	data.Set("pattern", resp.Pattern)
	data.Set("source_bucket", resp.SourceBucket)
	data.Set("column_selection", resp.ColumnSelection)

	// When the object in an Object Group use no compression, you need to create it with
	// `compression = ""`. However, when querying an Object Group whose object are not
	// compressed, the API returns `compression = "none"`. We coerce the "none" value to
	// an empty string in order not to confuse Terraform.
	compressionOrEmptyString := resp.Compression
	if strings.ToLower(compressionOrEmptyString) == "none" {
		compressionOrEmptyString = ""
	}
	log.Info("compression", compressionOrEmptyString)
	data.Set("compression", compressionOrEmptyString)
	data.Set("partition_by", resp.PartitionBy)
	data.Set("pattern", resp.Pattern)
	data.Set("source_bucket", resp.SourceBucket)

	data.Set("column_selection", resp.ColumnSelection)

	//"unlimited" flattening represented as "null" in the api, and as -1 in the terraform module
	//because the terraform sdk doesn't support nil values in configs https://github.com/hashicorp/terraform-plugin-sdk/issues/261
	//We represent "null" as an int pointer to nil in the code.
	if resp.ArrayFlattenDepth == nil {
		data.Set("array_flatten_depth", -1)
	} else {
		data.Set("array_flatten_depth", resp.ArrayFlattenDepth)
	}

	return diags
}

func resourceObjectGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	updateObjectGroupRequest := &client.UpdateObjectGroupRequest{
		AuthToken:             tokenValue,
		Bucket:                data.Get("bucket").(string),
		IndexParallelism:      data.Get("index_parallelism").(int),
		IndexRetention:        data.Get("index_retention_value").(int),
		TargetActiveIndex:     data.Get("target_active_index").(int),
		LiveEventsParallelism: data.Get("live_events_parallelism").(int),
	}

	if err := c.UpdateObjectGroup(ctx, updateObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	return resourceObjectGroupRead(ctx, data, meta)
}

func resourceObjectGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	deleteObjectGroupRequest := &client.DeleteObjectGroupRequest{
		AuthToken: tokenValue,
		Name:      data.Get("bucket").(string),
	}

	if err := c.DeleteObjectGroup(ctx, deleteObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))
	return nil
}

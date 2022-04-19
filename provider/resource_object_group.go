package provider

import (
	"context"
	"cs-tf-provider/client"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"source": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeSet,
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
				Type:        schema.TypeBool,
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
	var regex string

	if data.Get("filter").(*schema.Set).Len() > 0 {
		filterColumnSelectionInterface := data.Get("filter").(*schema.Set).List()[0]
		filterColumnSelection := filterColumnSelectionInterface.(map[string]interface{})

		prefixFilter := filterColumnSelection["prefix_filter"].(*schema.Set).List()[0].(map[string]interface{})
		regexFilter := filterColumnSelection["regex_filter"].(*schema.Set).List()[0].(map[string]interface{})

		prefixFilterField = prefixFilter["field"].(string)
		prefix = prefixFilter["prefix"].(string)

		regexFilterField = regexFilter["field"].(string)
		regex = regexFilter["regex"].(string)
	}

	prefixFilter := client.PrefixFilter{
		Field:  prefixFilterField,
		Prefix: prefix,
	}

	regexFilter := client.RegexFilter{
		Field: regexFilterField,
		Regex: regex,
	}
	filter := &client.Filter{
		PrefixFilter: &prefixFilter,
		RegexFilter:  &regexFilter,
	}
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value-->", tokenValue)
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

}

func resourceObjectGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).CSClient

	data.SetId(data.Get("bucket").(string))

	tokenValue := meta.(*ProviderMeta).token
	req := &client.ReadObjGroupReq{
		ID:        data.Id(),
		AuthToken: tokenValue,
	}

	resp, err := c.ReadObjGroup(ctx, req)

	if resp == nil {
		return diag.Errorf("Couldn't find object group: %s", err)
	}

	c.Set(data, "object_group_id", resp.ID)
	if err != nil {
		return diag.Errorf("Failed to read object group: %s", err)
	}

	c.Set(data, "name", data.Id())
	c.Set(data, "_public", resp.Public)
	c.Set(data, "_type", resp.Type)
	c.Set(data, "content_type", resp.ContentType)
	c.Set(data, "_realtime", resp.Realtime)
	c.Set(data, "bucket", resp.Bucket)

	var prefixFilterResponse = map[string]string{}
	for k, v := range resp.ObjectFilter.And[0].(map[string]interface{}) {
		prefixFilterResponse[k] = v.(string)
	}

	var regexFilterRes = map[string]string{}
	for k, v := range resp.ObjectFilter.And[1].(map[string]interface{}) {
		regexFilterRes[k] = v.(string)
	}

	PrefixFilterObjectMap := make(map[string]interface{})
	PrefixFilterObjectMap["field"] = prefixFilterResponse["field"]
	PrefixFilterObjectMap["prefix"] = prefixFilterResponse["prefix"]

	RegexFilterObjectMap := make(map[string]interface{})
	RegexFilterObjectMap["field"] = regexFilterRes["field"]
	RegexFilterObjectMap["regex"] = regexFilterRes["regex"]

	filter := make([]interface{}, 1)
	PrefixFilter := make([]interface{}, 1)
	PrefixFilter[0] = PrefixFilterObjectMap
	RegexFilter := make([]interface{}, 1)
	RegexFilter[0] = RegexFilterObjectMap
	filterObjectMap := make(map[string]interface{})
	filterObjectMap["prefix_filter"] = PrefixFilter
	filterObjectMap["regex_filter"] = RegexFilter
	filter[0] = filterObjectMap
	c.Set(data, "filter", filter)

	if resp.Format != nil {
		format := make([]interface{}, 1)
		formatObjectMap := make(map[string]interface{})
		formatObjectMap["_type"] = resp.Format.Type
		formatObjectMap["header_row"] = resp.Format.HeaderRow
		formatObjectMap["column_delimiter"] = resp.Format.ColumnDelimiter
		formatObjectMap["row_delimiter"] = resp.Format.RowDelimiter
		format[0] = formatObjectMap
		c.Set(data, "format", format)
	}

	if resp.Interval != nil {
		interval := make([]interface{}, 1)
		intervalObjectMap := make(map[string]interface{})
		intervalObjectMap["column"] = resp.Interval.Mode
		intervalObjectMap["mode"] = resp.Interval.Column
		interval[0] = intervalObjectMap
		c.Set(data, "interval", interval)
	}

	if resp.Metadata != nil {
		metadata := make([]interface{}, 1)
		metadataObjectMap := make(map[string]interface{})
		metadataObjectMap["creation_date"] = resp.Metadata.CreationDate
		metadata[0] = metadataObjectMap
		c.Set(data, "metadata", metadata)
	}

	if resp.Options != nil {
		options := make([]interface{}, 1)
		optionsObjectMap := make(map[string]interface{})
		optionsObjectMap["ignore_irregular"] = resp.Options.IgnoreIrregular
		options[0] = optionsObjectMap
		c.Set(data, "options", options)
	}

	RegionAvailability := make([]interface{}, 1)
	RegionAvailability[0] = resp.RegionAvailability
	c.Set(data, "region_availability", RegionAvailability[0])

	c.Set(data, "source", resp.Source)
	c.Set(data, "filter_json", resp.FilterJSON)
	c.Set(data, "live_events_sqs_arn", resp.LiveEventsSqsArn)

	c.Set(data, "partition_by", resp.PartitionBy)
	c.Set(data, "pattern", resp.Pattern)
	c.Set(data, "source_bucket", resp.SourceBucket)
	c.Set(data, "column_selection", resp.ColumnSelection)

	compressionOrEmptyString := resp.Compression
	if strings.ToLower(compressionOrEmptyString) == "none" {
		compressionOrEmptyString = ""
	}
	log.Info("compression", compressionOrEmptyString)
	c.Set(data, "compression", compressionOrEmptyString)
	c.Set(data, "partition_by", resp.PartitionBy)
	c.Set(data, "pattern", resp.Pattern)
	c.Set(data, "source_bucket", resp.SourceBucket)

	c.Set(data, "column_selection", resp.ColumnSelection)

	if resp.ArrayFlattenDepth == nil {
		c.Set(data, "array_flatten_depth", -1)
	} else {
		c.Set(data, "array_flatten_depth", resp.ArrayFlattenDepth)
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

package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/fatih/structs"
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
				ForceNew: true,
			},
			"source": {
				Type: schema.TypeString,
				//Required: true,
				Optional: true,
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
				Type: schema.TypeSet,
				//Required: true,
				Optional: true,
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
				Type: schema.TypeBool,
				//Required:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"index_parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"index_retention_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"target_active_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"live_events_parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
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
	log.Info("called READ")

	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).CSClient

	log.Info("keyyyyy====>", data.Get("objid"))
	//data.Get("objid")
	//bucketId := data.Get("objid").(string)
	log.Info("data.Id()", data.Id())

	var resourceGroupReqId string
	if data.Id() != "" {
		resourceGroupReqId = data.Id()
	} else if data.Get("object_group_id").(string) != "" {
		resourceGroupReqId = data.Get("object_group_id").(string)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unable to find id for object group",
			Detail:   "Unable to find id for object group",
		})
		return diags

	}
	//req := &client.ReadObjectGroupRequest{
	//	ID: bucketId,
	//}

	log.Info("33333333333333")
	req := &client.ReadObjectGroupRequest{
		ID: resourceGroupReqId,
	}
	log.Warn("req---->", req)
	resp, err := c.ReadObjectGroup(ctx, req)
	log.Info("response Id===", resp)
	if resp == nil {
		return diag.Errorf("Couldn't find object group: %s", err)
	}
	data.SetId(resp.ID)
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

	log.Info("ffffffffff====>", resp.Format)

	type Format struct {
		Type            string
		RowDelimiter    string
		HeaderRow       bool
		ColumnDelimiter string
	}
	format := Format{
		Type:            resp.Type,
		RowDelimiter:    resp.Format.RowDelimiter,
		HeaderRow:       resp.Format.HeaderRow,
		ColumnDelimiter: resp.Format.ColumnDelimiter,
	}
	structs.Map(format)
	//format := resp.Format

	//type format struct {
	//	Type        string `json:"_type"`
	//	Bucket      string `json:"bucket"`
	//	Bucket      string `json:"bucket"`
	//	Bucket      string `json:"bucket"`
	//	// contains filtered or unexported fields
	//}

	//TODO need to set format
	//data.Set("format", resp.Format)

	//data.Set("format.column_delimiter", resp.Format.HeaderRow)
	//data.Set("format.header_row", resp.Format.RowDelimiter)

	data.Set("Interval.Mode", resp.Interval.Mode)
	data.Set("RegionAvailability", resp.RegionAvailability)
	data.Set("Source", resp.Source)
	data.Set("Options", resp.Options)
	data.Set("Metadata.CreationDate", resp.Metadata.CreationDate)
	data.Set("Format.ColumnDelimiter", resp.Format.ColumnDelimiter)
	data.Set("Format.HeaderRow", resp.Format.HeaderRow)
	data.Set("Format.RowDelimiter", resp.Format.RowDelimiter)

	data.Set("filter_json", resp.FilterJSON)
	//data.Set("format", resp.Format)
	data.Set("live_events_sqs_arn", resp.LiveEventsSqsArn)

	log.Info("live_events_sqs_arn", resp.LiveEventsSqsArn)
	log.Info("index_retention", resp.IndexRetention)

	log.Info("partition_by", resp.PartitionBy)
	log.Info("pattern", resp.Pattern)
	log.Info("source_bucket", resp.SourceBucket)

	log.Info("column_selection", resp.ColumnSelection)

	//data.Set("index_retention", resp.IndexRetention)
	data.Set("partition_by", resp.PartitionBy)
	data.Set("pattern", resp.Pattern)
	data.Set("source_bucket", resp.SourceBucket)
	data.Set("column_selection", resp.ColumnSelection)

	log.Info("Type", resp.Type)
	log.Info("ContentType", resp.ContentType)
	log.Info("public=========>", resp.Public)
	log.Info("RealTime", resp.Realtime)
	log.Info("Type", resp.Type)
	log.Info("Bucket", resp.Bucket)
	log.Info("Interval.Column", resp.Interval.Column)
	log.Info("Interval.Mode", resp.Interval.Mode)
	log.Info("RegionAvailability", resp.RegionAvailability)
	log.Info("Source", resp.Source)
	log.Info("Options", resp.Options)
	log.Info("Metadata.CreationDate", resp.Metadata.CreationDate)
	log.Info("Format.ColumnDelimiter", resp.Format.ColumnDelimiter)
	log.Info("Format.HeaderRow", resp.Format.HeaderRow)
	log.Info("Format.RowDelimiter", resp.Format.RowDelimiter)
	log.Info("Format.filter", resp.Filter)
	log.Info("filter_json", resp.FilterJSON)
	log.Info("format", resp.Format)

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
	log.Debug("resourceObjectGroupUpdate called >>> ")
	log.Debug("target >>> ", data.Get("target_activeIndex").(int))
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

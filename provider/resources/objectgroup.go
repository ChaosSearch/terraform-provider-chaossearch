package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceObjectGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectGroupCreate,
		ReadContext:   ResourceObjectGroupRead,
		UpdateContext: resourceObjectGroupUpdate,
		DeleteContext: resourceObjectGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"format": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"column_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"header_row": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"row_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"index_retention": {
				Type:     schema.TypeSet,
				Optional: true,
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
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"regex_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
									},
									"regex": {
										Type:     schema.TypeString,
										Optional: true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"mode": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_irregular": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"realtime": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"index_parallelism": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"index_retention_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target_active_index": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"live_events_parallelism": {
				Type:     schema.TypeInt,
				Optional: true,
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
	var prefixFilter *client.PrefixFilter
	var prefixFilterField string
	var prefix string
	var regexFilter *client.RegexFilter
	var regexFilterField string
	var regex string

	c := meta.(*models.ProviderMeta).CSClient
	formatList := data.Get("format").(*schema.Set).List()
	formatMap := make(map[string]interface{}, 1)
	if len(formatList) > 0 {
		formatMap = formatList[0].(map[string]interface{})
	}

	format := client.Format{
		Type:            formatMap["type"].(string),
		ColumnDelimiter: formatMap["column_delimiter"].(string),
		RowDelimiter:    formatMap["row_delimiter"].(string),
		HeaderRow:       formatMap["header_row"].(bool),
	}

	intervalList := data.Get("interval").(*schema.Set).List()
	intervalMap := make(map[string]interface{})
	if len(intervalList) > 0 {
		intervalMap = intervalList[0].(map[string]interface{})
	}

	interval := client.Interval{
		Mode:   intervalMap["mode"].(int),
		Column: intervalMap["column"].(int),
	}

	indexList := data.Get("index_retention").(*schema.Set).List()
	indexMap := make(map[string]interface{})
	if len(indexList) > 0 {
		indexMap = indexList[0].(map[string]interface{})
	}

	indexRetention := client.IndexRetention{
		ForPartition: indexMap["for_partition"].([]interface{}),
		Overall:      indexMap["overall"].(int),
	}

	optionsList := data.Get("options").(*schema.Set).List()
	optionsMap := make(map[string]interface{})
	if len(optionsList) > 0 {
		optionsMap = optionsList[0].(map[string]interface{})
	}

	options := client.Options{
		IgnoreIrregular: optionsMap["ignore_irregular"].(bool),
	}

	filterSet := data.Get("filter").(*schema.Set)
	if filterSet.Len() > 0 {
		var prefixMap map[string]interface{}
		filterList := data.Get("filter").(*schema.Set).List()[0]
		filterMap := filterList.(map[string]interface{})

		prefixList := filterMap["prefix_filter"].(*schema.Set).List()
		if len(prefixList) > 0 {
			prefixMap = prefixList[0].(map[string]interface{})
			prefixFilterField = prefixMap["field"].(string)
			prefix = prefixMap["prefix"].(string)

			prefixFilter = &client.PrefixFilter{
				Field:  prefixFilterField,
				Prefix: prefix,
			}
		} else {
			prefixFilter = nil
		}

		regexList := filterMap["regex_filter"].(*schema.Set).List()
		if len(regexList) > 0 {
			regexMap := regexList[0].(map[string]interface{})
			regexFilterField = regexMap["field"].(string)
			regex = regexMap["regex"].(string)

			regexFilter = &client.RegexFilter{
				Field: regexFilterField,
				Regex: regex,
			}
		} else {
			regexFilter = nil
		}
	}

	filter := &client.Filter{
		PrefixFilter: prefixFilter,
		RegexFilter:  regexFilter,
	}

	tokenValue := meta.(*models.ProviderMeta).Token
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
	return ResourceObjectGroupRead(ctx, data, meta)
}

func ResourceObjectGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var prefixFilterMap = map[string]string{}
	var regexFilterMap = map[string]string{}

	diags := diag.Diagnostics{}
	c := meta.(*models.ProviderMeta).CSClient

	data.SetId(data.Get("bucket").(string))

	tokenValue := meta.(*models.ProviderMeta).Token
	req := &client.ReadObjGroupReq{
		ID:        data.Id(),
		AuthToken: tokenValue,
	}

	resp, err := c.ReadObjGroup(ctx, req)
	if err != nil {
		return diag.Errorf("Failed to read object group: %s", err)
	}

	if resp == nil {
		return diag.Errorf("Couldn't find object group: %s", err)
	}

	objectFilter := resp.ObjectFilter
	if len(objectFilter.And) > 0 {
		for index, filters := range objectFilter.And {
			for key, val := range filters.(map[string]interface{}) {
				if index == 0 {
					prefixFilterMap[key] = val.(string)
				} else if index == 1 {
					regexFilterMap[key] = val.(string)
				}
			}
		}
	}

	c.Set(data, "filter", []interface{}{
		map[string]interface{}{
			"prefix_filter": []interface{}{
				prefixFilterMap,
			},
			"regex_filter": []interface{}{
				regexFilterMap,
			},
		},
	})

	if resp.Format != nil {
		c.Set(data, "format", []interface{}{
			map[string]interface{}{
				"type":             resp.Format.Type,
				"header_row":       resp.Format.HeaderRow,
				"column_delimiter": resp.Format.ColumnDelimiter,
				"row_delimiter":    resp.Format.RowDelimiter,
			},
		})
	}

	if resp.Interval != nil {
		c.Set(data, "interval", []interface{}{
			map[string]interface{}{
				"column": resp.Interval.Column,
				"mode":   resp.Interval.Mode,
			},
		})
	}

	if resp.Metadata != nil {
		c.Set(data, "metadata", []interface{}{
			map[string]interface{}{
				"creation_date": resp.Metadata.CreationDate,
			},
		})
	}

	if resp.Options != nil {
		c.Set(data, "options", []interface{}{
			map[string]interface{}{
				"ignore_irregular": resp.Options.IgnoreIrregular,
			},
		})
	}

	if strings.ToLower(resp.Compression) == "none" {
		c.Set(data, "compression", "")
	} else {
		c.Set(data, "compression", resp.Compression)
	}

	if resp.ArrayFlattenDepth == nil {
		c.Set(data, "array_flatten_depth", -1)
	} else {
		c.Set(data, "array_flatten_depth", resp.ArrayFlattenDepth)
	}

	c.Set(data, "region_availability", []interface{}{
		resp.RegionAvailability,
	})

	c.Set(data, "id", resp.ID)
	c.Set(data, "name", data.Id())
	c.Set(data, "public", resp.Public)
	c.Set(data, "type", resp.Type)
	c.Set(data, "content_type", resp.ContentType)
	c.Set(data, "realtime", resp.Realtime)
	c.Set(data, "bucket", resp.Bucket)
	c.Set(data, "source", resp.Source)
	c.Set(data, "filter_json", resp.FilterJSON)
	c.Set(data, "live_events_sqs_arn", resp.LiveEventsSqsArn)
	c.Set(data, "partition_by", resp.PartitionBy)
	c.Set(data, "pattern", resp.Pattern)
	c.Set(data, "source_bucket", resp.SourceBucket)
	c.Set(data, "column_selection", resp.ColumnSelection)

	return diags
}

func resourceObjectGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
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

	return ResourceObjectGroupRead(ctx, data, meta)
}

func resourceObjectGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
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

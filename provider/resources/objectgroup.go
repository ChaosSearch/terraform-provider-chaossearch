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
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_type": {
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
							Default:  -1,
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
							Optional: true,
							Default:  0,
						},
						"mode": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_date": {
							Type:     schema.TypeInt,
							Computed: true,
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
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"region_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"compression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"array_flatten_depth": {
				Type:     schema.TypeInt,
				Computed: true,
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
	var format *client.Format
	var interval *client.Interval
	var indexRetention *client.IndexRetention
	var options *client.Options

	c := meta.(*models.ProviderMeta).CSClient
	formatList := data.Get("format").(*schema.Set).List()
	if len(formatList) > 0 {
		var typeStr string
		var columnDelimit string
		var rowDelimit string
		var headerRow bool

		formatMap := formatList[0].(map[string]interface{})
		if formatMap["type"] != nil {
			typeStr = formatMap["type"].(string)
		}

		if formatMap["column_delimiter"] != nil {
			columnDelimit = formatMap["column_delimiter"].(string)
		}

		if formatMap["row_delimiter"] != nil {
			rowDelimit = formatMap["row_delimiter"].(string)
		}

		if formatMap["header_row"] != nil {
			headerRow = formatMap["header_row"].(bool)
		}
		format = &client.Format{
			Type:            typeStr,
			ColumnDelimiter: columnDelimit,
			RowDelimiter:    rowDelimit,
			HeaderRow:       headerRow,
		}
	}

	intervalList := data.Get("interval").(*schema.Set).List()
	if len(intervalList) > 0 {
		var mode int
		var column int
		intervalMap := intervalList[0].(map[string]interface{})

		if intervalMap["mode"] != nil {
			mode = intervalMap["mode"].(int)
		}

		if intervalMap["column"] != nil {
			column = intervalMap["column"].(int)
		}

		interval = &client.Interval{
			Mode:   mode,
			Column: column,
		}
	}

	indexList := data.Get("index_retention").(*schema.Set).List()
	if len(indexList) > 0 {
		var forPartition []interface{}
		var overall int

		indexMap := indexList[0].(map[string]interface{})
		if indexMap["for_partition"] != nil {
			forPartition = indexMap["for_partition"].([]interface{})
		}

		if indexMap["overall"] != nil {
			overall = indexMap["overall"].(int)
		}
		indexRetention = &client.IndexRetention{
			ForPartition: forPartition,
			Overall:      overall,
		}
	}

	optionsList := data.Get("options").(*schema.Set).List()
	if len(optionsList) > 0 {
		var ignore bool
		optionsMap := optionsList[0].(map[string]interface{})
		if optionsMap["ignore_irregular"] != nil {
			ignore = optionsMap["ignore_irregular"].(bool)
		}

		options = &client.Options{
			IgnoreIrregular: ignore,
		}
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
		Format:         format,
		Interval:       interval,
		IndexRetention: indexRetention,
		Filter:         filter,
		Options:        options,
		Realtime:       false,
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
		for _, filter := range objectFilter.And {
			filterMap := filter.(map[string]interface{})
			for key, val := range filterMap {
				if _, ok := filterMap["prefix"]; ok {
					prefixFilterMap[key] = val.(string)
				} else if _, ok := filterMap["regex"]; ok {
					regexFilterMap[key] = val.(string)
				}
			}
		}
	}

	filterArr := []interface{}{}
	if len(prefixFilterMap) > 0 {
		filterArr = append(filterArr, map[string]interface{}{
			"prefix_filter": []interface{}{
				prefixFilterMap,
			},
		})
	}

	if len(regexFilterMap) > 0 {
		filterArr = append(filterArr, map[string]interface{}{
			"regex_filter": []interface{}{
				regexFilterMap,
			},
		})
	}

	err = data.Set("filter", filterArr)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Format != nil {
		err = data.Set("format", []interface{}{
			map[string]interface{}{
				"type":             resp.Format.Type,
				"header_row":       resp.Format.HeaderRow,
				"column_delimiter": resp.Format.ColumnDelimiter,
				"row_delimiter":    resp.Format.RowDelimiter,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if resp.Interval != nil {
		err = data.Set("interval", []interface{}{
			map[string]interface{}{
				"column": resp.Interval.Column,
				"mode":   resp.Interval.Mode,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if resp.Metadata != nil {
		err = data.Set("metadata", []interface{}{
			map[string]interface{}{
				"creation_date": resp.Metadata.CreationDate,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if resp.Options != nil {
		err = data.Set("options", []interface{}{
			map[string]interface{}{
				"ignore_irregular": resp.Options.IgnoreIrregular,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if strings.ToLower(resp.Compression) == "none" {
		err = data.Set("compression", "")
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = data.Set("compression", resp.Compression)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if resp.ArrayFlattenDepth == nil {
		err = data.Set("array_flatten_depth", -1)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = data.Set("array_flatten_depth", resp.ArrayFlattenDepth)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = data.Set("region_availability", resp.RegionAvailability)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("id", resp.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("public", resp.Public)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("type", resp.Type)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("content_type", resp.ContentType)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("bucket", resp.Bucket)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("source", resp.Source)
	if err != nil {
		return diag.FromErr(err)
	}

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

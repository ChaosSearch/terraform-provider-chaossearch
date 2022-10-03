package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/models"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"format": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"column_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"header_row": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"row_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"array_flatten_depth": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"strip_prefix": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"horizontal": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"live_events": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"index_retention": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"overall": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"range": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"max": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"equals": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"regex": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"interval": {
				Type:     schema.TypeSet,
				Computed: true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ignore_irregular": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"col_types": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
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
			"target_active_index": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"live_events_parallelism": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}

}

func resourceObjectGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var format *client.Format
	var indexRetention *client.IndexRetention
	FILTER_FIELDS := []string{
		"lastModified",
		"size",
		"storageClass",
		"key",
	}

	rangeFields := []string{FILTER_FIELDS[0], FILTER_FIELDS[1]}
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

	indexList := data.Get("index_retention").(*schema.Set).List()
	if len(indexList) > 0 {
		var overall int

		indexMap := indexList[0].(map[string]interface{})
		if indexMap["overall"] != nil {
			overall = indexMap["overall"].(int)
		}
		indexRetention = &client.IndexRetention{
			ForPartition: []interface{}{},
			Overall:      overall,
		}
	}

	filters := []client.Filter{}
	filterList := data.Get("filter").([]interface{})
	if len(filterList) > 0 {
		for _, filterSet := range filterList {
			filterMap := filterSet.(map[string]interface{})
			field := filterMap["field"].(string)

			if !utils.ContainsString(FILTER_FIELDS, field) {
				return utils.CreateObjectGroupError(
					fmt.Sprintf(`Invalid field supplied: %s Acceptable values are: %v`, field, FILTER_FIELDS),
				)
			}

			rangeList := filterMap["range"].(*schema.Set).List()
			equals := filterMap["equals"].(string)
			regex := filterMap["regex"].(string)
			prefix := filterMap["prefix"].(string)

			if len(rangeList) > 0 && utils.ContainsString(rangeFields, field) {
				return utils.CreateObjectGroupError(`Range is currently not supported`)
				/*
					rangeMap := rangeList[0].(map[string]interface{})
					filters = append(filters, client.Filter{
						Field: field,
						Range: client.Range{
							Min: rangeMap["min"].(string),
						},
					})

					filters = append(filters, client.Filter{
						Field: field,
						Range: client.Range{
							Max: rangeMap["max"].(string),
						},
					})
				*/
			} else if len(rangeList) > 0 && !utils.ContainsString(rangeFields, field) {
				return utils.CreateObjectGroupError(
					fmt.Sprintf(`Range used with incompatible field. Range can only be used with %v`, rangeFields),
				)
			}

			if field == FILTER_FIELDS[2] && equals != "" {
				filters = append(filters, client.Filter{
					Field:  field,
					Equals: equals,
				})
			} else if field == FILTER_FIELDS[2] && equals == "" {
				return utils.CreateObjectGroupError("'storageClass' field must be used with equals param")
			}

			if field == FILTER_FIELDS[3] && regex == "" && prefix == "" {
				return utils.CreateObjectGroupError("'key' field requires either regex or prefix be defined")
			}

			if field == FILTER_FIELDS[3] && regex != "" && prefix != "" {
				return utils.CreateObjectGroupError(
					"'key' field requires only one of regex or prefix to be defined \n" +
						"Note: Break these out into separate filter blocks. One containing regex and the other with prefix",
				)
			}

			if field == FILTER_FIELDS[3] && regex != "" {
				filters = append(filters, client.Filter{
					Field: field,
					Regex: regex,
				})
			}

			if field == FILTER_FIELDS[3] && prefix != "" {
				filters = append(filters, client.Filter{
					Field:  field,
					Prefix: prefix,
				})
			}
		}
	}

	options := &client.Options{
		IgnoreIrregular: true,
	}

	optionsList := data.Get("options").(*schema.Set).List()
	if len(optionsList) > 0 {
		optionsMap := optionsList[0].(map[string]interface{})
		options.Compression = optionsMap["compression"].(string)
		colTypesString := optionsMap["col_types"].(string)
		err := json.Unmarshal([]byte(colTypesString), &options.ColTypes)
		if err != nil {
			return diag.FromErr(utils.UnmarshalJsonError(err))
		}
	}

	tokenValue := meta.(*models.ProviderMeta).Token
	createObjectGroupRequest := &client.CreateObjectGroupRequest{
		AuthToken:      tokenValue,
		Bucket:         data.Get("bucket").(string),
		Source:         data.Get("source").(string),
		Format:         format,
		IndexRetention: indexRetention,
		Filter:         filters,
		LiveEvents:     data.Get("live_events").(string),
		Options:        options,
		Interval: &client.Interval{
			Mode:   0,
			Column: 0,
		},
		Realtime: false,
	}

	if err := c.CreateObjectGroup(ctx, createObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	return ResourceObjectGroupRead(ctx, data, meta)
}

func ResourceObjectGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var filters = []map[string]interface{}{}

	diags := diag.Diagnostics{}
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	req := &client.BasicRequest{
		Id:        data.Get("bucket").(string),
		AuthToken: tokenValue,
	}

	resp, err := c.ReadObjGroup(ctx, req)
	if err != nil {
		return diag.Errorf("Failed to read object group: %s", err)
	}

	if resp == nil {
		return diag.Errorf("Couldn't find object group: %s", err)
	}

	data.SetId(resp.ID)
	objectFilter := resp.ObjectFilter
	if len(objectFilter.And) > 0 {
		for _, filter := range objectFilter.And {
			filterMap := filter.(map[string]interface{})
			for key, val := range filterMap {
				if _, ok := filterMap["prefix"]; ok {
					filters = append(filters, map[string]interface{}{
						"field":  key,
						"prefix": val.(string),
					})
				} else if _, ok := filterMap["regex"]; ok {
					filters = append(filters, map[string]interface{}{
						"field": key,
						"regex": val.(string),
					})
				} else if _, ok := filterMap["equals"]; ok {
					filters = append(filters, map[string]interface{}{
						"field":  key,
						"equals": val.(string),
					})
				}
			}
		}
	}

	err = data.Set("filter", filters)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Format != nil {
		err = data.Set("format", []interface{}{
			map[string]interface{}{
				"type":                resp.Format.Type,
				"header_row":          resp.Format.HeaderRow,
				"column_delimiter":    resp.Format.ColumnDelimiter,
				"row_delimiter":       resp.Format.RowDelimiter,
				"array_flatten_depth": resp.ArrayFlattenDepth,
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
				"compression":      resp.Compression,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = data.Set("region_availability", resp.RegionAvailability)
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

	err = data.Set("source_id", resp.Source)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceObjectGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var indexRetention int
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token

	indexList := data.Get("index_retention").(*schema.Set).List()
	if len(indexList) > 0 {
		indexMap := indexList[0].(map[string]interface{})
		if indexMap["overall"] != nil {
			indexRetention = indexMap["overall"].(int)
		}
	}

	if indexRetention == 0 || indexRetention < -1 {
		return diag.Errorf(`Failure Updating Object Group => Invalid Index Retention
			Note:
				index_retention.overall cannot == 0 or < -1 during update
		`)
	}

	activeIndex := data.Get("target_active_index").(int)
	if activeIndex == 0 || activeIndex < -1 {
		return diag.Errorf(`Failure Updating Object Group => Invalid Active Index
			Note:
				target_active_index cannot == 0 or < -1 during update
				This value is optional on create, but required on update.
		`)
	}

	updateObjectGroupRequest := &client.UpdateObjectGroupRequest{
		AuthToken:             tokenValue,
		Bucket:                data.Get("bucket").(string),
		IndexParallelism:      data.Get("index_parallelism").(int),
		IndexRetention:        indexRetention,
		TargetActiveIndex:     activeIndex,
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
	deleteObjectGroupRequest := &client.BasicRequest{
		AuthToken: tokenValue,
		Id:        data.Get("bucket").(string),
	}

	if err := c.DeleteObjectGroup(ctx, deleteObjectGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))
	return nil
}

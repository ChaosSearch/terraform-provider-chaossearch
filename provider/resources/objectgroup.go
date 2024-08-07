package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/models"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
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
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"array_flatten_depth": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"horizontal": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"array_selection": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSelectionPolicy,
						},
						"field_selection": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSelectionPolicy,
						},
						"vertical_selection": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSelectionPolicy,
						},
					},
				},
			},
			"live_events": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"live_events_aws", "live_events_gcp"},
				Deprecated:    "Use `live_events_aws` for this value",
			},
			"live_events_aws": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"live_events", "live_events_gcp"},
			},
			"live_events_gcp": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"live_events_aws", "live_events"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subscription_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"index_retention": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"overall": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIndexVal,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"range": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"max": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"equals": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"regex": {
							Type:     schema.TypeString,
							Optional: true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"col_types": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"col_renames": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"col_selection": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSelectionPolicy,
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIndexVal,
			},
			"live_events_parallelism": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"partition_by": {
				Type:     schema.TypeString,
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
		var arrayFlattenDepth int
		var stripPrefix bool
		var horizontal bool
		var arraySelection []map[string]interface{}
		var fieldSelection []map[string]interface{}
		var verticalSelection []map[string]interface{}

		formatMap := formatList[0].(map[string]interface{})
		if formatMap["type"] != nil {
			typeStr = formatMap["type"].(string)
		}

		if typeStr != "JSON" {
			stripPrefix = false
		} else {
			stripPrefix = true
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

		if formatMap["array_flatten_depth"] != nil {
			arrayFlattenDepth = formatMap["array_flatten_depth"].(int)
		}

		if formatMap["horizontal"] != nil {
			horizontal = formatMap["horizontal"].(bool)
		}

		arraySelectString := formatMap["array_selection"].(string)
		if arraySelectString != "" {
			arraySelectJson, err := structure.NormalizeJsonString(arraySelectString)
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}

			err = json.Unmarshal([]byte(arraySelectJson), &arraySelection)
			if err != nil {
				return diag.FromErr(
					fmt.Errorf("Object Group Resource Failure => %s \n %s", utils.UnmarshalJsonError(err), arraySelectJson),
				)
			}
		}

		fieldSelectString := formatMap["field_selection"].(string)
		if fieldSelectString != "" {
			fieldSelectJson, err := structure.NormalizeJsonString(fieldSelectString)
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}

			err = json.Unmarshal([]byte(fieldSelectJson), &fieldSelection)
			if err != nil {
				return diag.FromErr(
					fmt.Errorf("Object Group Resource Failure => %s \n %s", utils.UnmarshalJsonError(err), fieldSelectJson),
				)
			}
		}

		verticalSelectStr := formatMap["vertical_selection"].(string)
		if verticalSelectStr != "" {
			verticalSelectJson, err := structure.NormalizeJsonString(verticalSelectStr)
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}

			err = json.Unmarshal([]byte(verticalSelectJson), &verticalSelection)
			if err != nil {
				return diag.FromErr(
					fmt.Errorf("Object Group Resource Failure => %s \n %s", utils.UnmarshalJsonError(err), verticalSelectJson),
				)
			}
		}

		format = &client.Format{
			Type:              typeStr,
			ColumnDelimiter:   columnDelimit,
			RowDelimiter:      rowDelimit,
			HeaderRow:         headerRow,
			ArrayFlattenDepth: arrayFlattenDepth,
			StripPrefix:       stripPrefix,
			Horizontal:        horizontal,
			ArraySelection:    arraySelection,
			FieldSelection:    fieldSelection,
			VerticalSelection: verticalSelection,
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
		if colTypesString != "" {
			err := json.Unmarshal([]byte(colTypesString), &options.ColTypes)
			if err != nil {
				return diag.FromErr(utils.UnmarshalJsonError(err))
			}
		}

		colRenamesString := optionsMap["col_renames"].(string)
		if colRenamesString != "" {
			err := json.Unmarshal([]byte(colRenamesString), &options.ColRenames)
			if err != nil {
				return diag.FromErr(utils.UnmarshalJsonError(err))
			}
		}

		colSelectionString := optionsMap["col_selection"].(string)
		if colSelectionString != "" {
			err := json.Unmarshal([]byte(colSelectionString), &options.ColSelection)
			if err != nil {
				return diag.FromErr(utils.UnmarshalJsonError(err))
			}
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
		PartitionBy:    data.Get("partition_by").(string),
		Options:        options,
		Interval: &client.Interval{
			Mode:   0,
			Column: 0,
		},
		Realtime: false,
	}

	if targetActiveIndex, ok := data.GetOk("target_active_index"); ok {
		val := targetActiveIndex.(int)
		createObjectGroupRequest.TargetActiveIndex = &val
	}

	liveEvents := data.Get("live_events").(string)
	liveEventsAws := data.Get("live_events_aws").(string)
	liveEventsGcp := data.Get("live_events_gcp").(*schema.Set).List()

	if liveEvents != "" && liveEventsAws != "" {
		return diag.Errorf("Both live_events and live_events_aws found defined, please only use one")
	} else if liveEvents != "" {
		liveEventsAws = liveEvents
	}

	if liveEventsAws != "" && len(liveEventsGcp) > 0 {
		err := "Live Events found defined for both AWS and GCP, please ensure you configure only one for your cluster type"
		return diag.Errorf(err)
	} else if liveEventsAws != "" {
		createObjectGroupRequest.LiveEventsAws = liveEventsAws
	} else if len(liveEventsGcp) > 0 {
		liveEventsMap := liveEventsGcp[0].(map[string]interface{})
		createObjectGroupRequest.LiveEventsGcp = &client.LiveEventsGcp{
			ProjectId:      liveEventsMap["project_id"].(string),
			SubscriptionId: liveEventsMap["subscription_id"].(string),
		}
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
			field := filterMap["field"].(string)
			if _, ok := filterMap["prefix"]; ok {
				filters = append(filters, map[string]interface{}{
					"field":  field,
					"prefix": filterMap["prefix"].(string),
				})
			} else if _, ok := filterMap["regex"]; ok {
				var regex string
				regex_map, ok := filterMap["regex"].(map[string]interface{})
				if !ok {
					regex = filterMap["regex"].(string)
				} else {
					regex = regex_map["pattern"].(string)
				}

				filters = append(filters, map[string]interface{}{
					"field": field,
					"regex": regex,
				})
			} else if _, ok := filterMap["equals"]; ok {
				filters = append(filters, map[string]interface{}{
					"field":  field,
					"equals": filterMap["equals"].(string),
				})
			}
		}
	}

	err = data.Set("filter", filters)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Format != nil {
		var arraySelectString string
		var fieldSelectString string
		var vertSelectString string

		formatRespMap := map[string]interface{}{
			"type":             resp.Format.Type,
			"header_row":       resp.Format.HeaderRow,
			"column_delimiter": resp.Format.ColumnDelimiter,
			"row_delimiter":    resp.Format.RowDelimiter,
			"horizontal":       resp.Format.Horizontal,
		}

		arraySelect := resp.Format.ArraySelection
		if len(arraySelect) > 0 {

			selectJson, err := json.Marshal(arraySelect)
			if err != nil {
				return diag.FromErr(utils.MarshalJsonError(err))
			}

			arraySelectString, err = structure.NormalizeJsonString(string(selectJson))
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}
		}

		fieldSelect := resp.Format.FieldSelection
		if len(fieldSelect) > 0 {
			selectJson, err := json.Marshal(fieldSelect)
			if err != nil {
				return diag.FromErr(utils.MarshalJsonError(err))
			}

			fieldSelectString, err = structure.NormalizeJsonString(string(selectJson))
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}
		}

		vertSelect := resp.Format.VerticalSelection
		if len(vertSelect) > 0 {
			selectJson, err := json.Marshal(vertSelect)
			if err != nil {
				return diag.FromErr(utils.MarshalJsonError(err))
			}

			vertSelectString, err = structure.NormalizeJsonString(string(selectJson))
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}
		}

		formatList := data.Get("format").(*schema.Set).List()
		if len(formatList) > 0 {
			var arrayMap []map[string]interface{}
			var arrayStateMap []map[string]interface{}
			var fieldMap []map[string]interface{}
			var fieldStateMap []map[string]interface{}
			var vertMap []map[string]interface{}
			var vertStateMap []map[string]interface{}
			formatMap := formatList[0].(map[string]interface{})

			arraySelectStateString := formatMap["array_selection"].(string)
			_ = json.Unmarshal([]byte(arraySelectString), &arrayMap)
			_ = json.Unmarshal([]byte(arraySelectStateString), &arrayStateMap)
			if !reflect.DeepEqual(arrayMap, arrayStateMap) {
				formatRespMap["array_selection"] = arraySelectString
			} else {
				formatRespMap["array_selection"] = arraySelectStateString
			}

			fieldSelectStateString := formatMap["field_selection"].(string)
			_ = json.Unmarshal([]byte(fieldSelectString), &fieldMap)
			_ = json.Unmarshal([]byte(fieldSelectStateString), &fieldStateMap)
			if !reflect.DeepEqual(fieldMap, fieldStateMap) {
				formatRespMap["field_selection"] = fieldSelectString
			} else {
				formatRespMap["field_selection"] = fieldSelectStateString
			}

			vertSelectStateString := formatMap["vertical_selection"].(string)
			_ = json.Unmarshal([]byte(vertSelectString), &vertMap)
			_ = json.Unmarshal([]byte(vertSelectStateString), &vertStateMap)
			if !reflect.DeepEqual(vertMap, vertStateMap) {
				formatRespMap["vertical_selection"] = vertSelectString
			} else {
				formatRespMap["vertical_selection"] = vertSelectStateString
			}

			formatRespMap["array_flatten_depth"] = formatMap["array_flatten_depth"]
		} else {
			formatRespMap["array_selection"] = arraySelectString
			formatRespMap["field_selection"] = fieldSelectString
			formatRespMap["vertical_selection"] = vertSelectString
		}

		err = data.Set("format", []interface{}{formatRespMap})

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

	optionsList := data.Get("options").(*schema.Set).List()
	if resp.Options != nil && len(optionsList) > 0 {
		var options map[string]interface{}

		optionsMap := optionsList[0].(map[string]interface{})
		compression := optionsMap["compression"].(string)

		if strings.EqualFold(resp.Options.Compression, compression) {
			options = map[string]interface{}{"compression": compression}
		} else {
			options = map[string]interface{}{"compression": resp.Options.Compression}
		}

		colTypesString := optionsMap["col_types"].(string)
		if colTypesString != "" {
			options["col_types"] = colTypesString
		} else if len(resp.Options.ColTypes) > 0 {
			colTypes, _ := json.Marshal(resp.Options.ColTypes)
			colTypesJson, _ := structure.NormalizeJsonString(string(colTypes))
			options["col_types"] = colTypesJson
		}

		colSelectionString := optionsMap["col_selection"].(string)
		if colSelectionString != "" {
			options["col_selection"] = colSelectionString
		} else if len(resp.Options.ColSelection) > 0 {
			colSelect, _ := json.Marshal(resp.Options.ColSelection)
			colSelectJson, _ := structure.NormalizeJsonString(string(colSelect))
			options["col_selection"] = colSelectJson
		}

		colRenamesString := optionsMap["col_renames"].(string)
		if colRenamesString != "" {
			options["col_renames"] = colRenamesString
		} else if len(resp.Options.ColRenames) > 0 {
			colRenames, _ := json.Marshal(resp.Options.ColRenames)
			colRenamesJson, _ := structure.NormalizeJsonString(string(colRenames))
			options["col_renames"] = colRenamesJson
		}

		err = data.Set("options", []interface{}{options})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if resp.PartitionBy != nil {
		var pattern string
		pattern, ok := resp.PartitionBy.(string)
		if !ok {
			patternMap := resp.PartitionBy.(map[string]interface{})
			byList := patternMap["by"].([]interface{})
			pattern = byList[0].(map[string]interface{})["pattern"].(string)
		}

		err = data.Set("partition_by", pattern)
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
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	hasChanges := data.HasChanges(
		"target_active_index",
		"index_parallelism",
		"live_events_parallelism",
	)

	if data.HasChange("index_retention") {
		oldVal, newVal := data.GetChange("index_retention")
		if oldVal != nil && newVal != nil {
			var oldOverall int
			var newOverall int
			if oldValList := oldVal.(*schema.Set).List(); len(oldValList) > 0 {
				oldOverall = oldValList[0].(map[string]interface{})["overall"].(int)
			}

			if newValList := newVal.(*schema.Set).List(); len(newValList) > 0 {
				newOverall = newValList[0].(map[string]interface{})["overall"].(int)
			}

			if oldOverall != newOverall {
				hasChanges = true
			}
		}
	}

	if hasChanges {
		req := &client.UpdateObjectGroupRequest{
			AuthToken: tokenValue,
			Bucket:    data.Get("bucket").(string),
		}

		indexList := data.Get("index_retention").(*schema.Set).List()
		if len(indexList) > 0 {
			indexMap := indexList[0].(map[string]interface{})
			if indexMap["overall"] != nil {
				val := indexMap["overall"].(int)
				req.IndexRetention = &val
			}
		}

		if activeIndex, ok := data.GetOk("target_active_index"); ok {
			val := activeIndex.(int)
			req.TargetActiveIndex = &val
		}

		if indexParallelism, ok := data.GetOk("index_parallelism"); ok {
			val := indexParallelism.(int)
			req.IndexParallelism = &val
		}

		if liveParallelism, ok := data.GetOk("live_events_parallelism"); ok {
			val := liveParallelism.(int)
			req.IndexParallelism = &val
		}

		if err := c.UpdateObjectGroup(ctx, req); err != nil {
			return diag.FromErr(err)
		}
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

func validateSelectionPolicy(i interface{}, k string) (warnings []string, errors []error) {
	var policyArr []map[string]interface{}
	types := []string{
		"blacklist",
		"whitelist",
		"regex",
	}

	warnings, errors = validation.StringIsJSON(i, k)
	if len(warnings) > 0 || len(errors) > 0 {
		return warnings, errors
	}

	policyJson, _ := structure.NormalizeJsonString(i.(string))
	err := json.Unmarshal([]byte(policyJson), &policyArr)
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	for _, policyMap := range policyArr {
		selectType, ok := policyMap["type"].(string)
		if !ok {
			errors = append(errors, fmt.Errorf("json expected to have 'type' field for select policy"))
			return warnings, errors
		}

		if !utils.ContainsString(types, selectType) {
			errors = append(errors, fmt.Errorf("invalid type found for selection policy. type: %s", selectType))
			return warnings, errors
		}

		if selectType == types[0] {
			excludes := policyMap["excludes"]
			if excludes == nil {
				errors = append(errors, fmt.Errorf("selection json expected an 'excludes' array in the case of type %s", types[0]))
			}
		}

		if selectType == types[1] {
			includes := policyMap["includes"]
			if includes == nil {
				errors = append(errors, fmt.Errorf("selection json expected an 'includes' array in the case of type %s", types[1]))
			}
		}

		if selectType == types[2] {
			patterns := policyMap["patterns"]
			if patterns == nil {
				errors = append(errors, fmt.Errorf("selection json expected a 'patterns' array in the case of type %s", types[2]))
			}
		}
	}

	return warnings, errors
}

func validateIndexVal(i interface{}, k string) (warnings []string, errors []error) {
	if i != nil && (i.(int) == 0 || i.(int) < -1) {
		return nil, []error{fmt.Errorf("Value cannot == 0 or be < -1")}
	}

	return nil, nil
}

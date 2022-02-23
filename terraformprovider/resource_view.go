package main

import (
	"context"
	"cs-tf-provider/client"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	log "github.com/sirupsen/logrus"
)

func resourceView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceViewCreate,
		ReadContext:   resourceViewRead,
		UpdateContext: resourceViewCreate,
		DeleteContext: resourceViewDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"index_pattern": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
			},

			"overwrite": {
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
				Default:  false,
				Optional: true,
			},
			"case_insensitive": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"index_retention": {
				Type:        schema.TypeInt,
				Default:     14,
				Description: "",
				Optional:    true,
				ForceNew:    false,
			},
			"time_field_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"transforms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predicate": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pred": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: false,
												},
												"_type": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: false,
												},
												"query": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: false,
												},
												"state": {
													Type:     schema.TypeSet,
													Required: true,
													ForceNew: false,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"_type": {
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
									"_type": {
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
		},
	}
}

func resourceViewCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	filterColumnSelectionInterface := data.Get("filter").(*schema.Set).List()[0].(map[string]interface{})
	predicateColumnSelectionInterface := filterColumnSelectionInterface["predicate"].(*schema.Set).List()[0].(map[string]interface{})
	predColumnSelectionInterface := predicateColumnSelectionInterface["pred"].(*schema.Set).List()[0].(map[string]interface{})
	stateColumnSelectionInterface := predColumnSelectionInterface["state"].(*schema.Set).List()[0].(map[string]interface{})

	var predicateType string

	var fieldValue string
	var queryValue string
	var typeValue string

	var stateType string

	predicateType = predicateColumnSelectionInterface["_type"].(string)
	fieldValue = predColumnSelectionInterface["field"].(string)
	queryValue = predColumnSelectionInterface["query"].(string)
	typeValue = predColumnSelectionInterface["_type"].(string)
	stateType = stateColumnSelectionInterface["_type"].(string)

	state := client.State{
		Type_: stateType,
	}

	pred := client.Pred{
		Field: fieldValue,
		Query: queryValue,
		State: state,
		Type_: typeValue,
	}

	Predicate := client.Predicate{
		Type_: predicateType,
		Pred:  pred,
	}

	filter := &client.FilterPredicate{
		Predicate: &Predicate,
	}

	//filter := client.Filter{
	//	Predicate: filterColumnSelectionInterface["predicate"].(*schema.Set),
	//}

	// predicate := client.Predicate{
	// 	Field: predicateColumnSelectionInterface["field"].(string),
	// 	Query: predicateColumnSelectionInterface["query"].(string),
	// 	State: predicateColumnSelectionInterface["state"].(*client.State),
	// 	Type_: predicateColumnSelectionInterface["_type"].(string),
	// }

	// state := client.State{
	// 	Type_: predicateColumnSelectionInterface["_type"].(string),
	// }

	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value------------>>>>", tokenValue)

	// arrayFlattenTF := data.Get("array_flatten_depth").(int)
	// log.Info("arrayFlattenTF-->",arrayFlattenTF)
	// var arrayFlattenCS *int

	// if arrayFlattenTF == -1 {
	// -1 in terraform represents "null" in the ChaosSearch API call
	// arrayFlattenCS = nil
	// } else {
	// any other value is passed as is
	// arrayFlattenCS = &arrayFlattenTF
	// }

	// var indexRetention map[string]interface{}

	// if data.Get("index_retention").(*schema.Set).Len() > 0 {
	// 	columnSelectionInterfaces := data.Get("index_retention").(*schema.Set).List()[0]
	// 	columnSelectionInterface := columnSelectionInterfaces.(map[string]interface{})

	// 	indexRetention = map[string]interface{}{
	// 		"value": columnSelectionInterface["value"].(string),
	// 	}
	// }
	// log.Debug("indexretention", indexRetention)
	sources_, ok := data.GetOk("sources")
	if !ok {
		log.Error(" sources not available")
	}
	log.Debug("sources_-->", sources_)
	var sourcesStrings []interface{}

	if sources_ != nil {
		sourcesStrings = sources_.([]interface{})
	}

	transforms_, ok := data.GetOk("transforms")
	if !ok {
		log.Error(" transforms not available")
	}
	var transforms []interface{}

	if transforms_ != nil {
		transforms = transforms_.([]interface{})
	}

	//patterns_, ok := data.GetOk("pattern")
	//if !ok {
	//	log.Error(" sources not available")
	//}
	//patterns := patterns_.([]interface{})
	createViewRequest := &client.CreateViewRequest{
		AuthToken: tokenValue,

		Bucket:          data.Get("bucket").(string),
		Sources:         sourcesStrings,
		IndexPattern:    data.Get("index_pattern").(string),
		Overwrite:       data.Get("overwrite").(bool),
		CaseInsensitive: data.Get("case_insensitive").(bool),
		IndexRetention:  data.Get("index_retention").(int),
		TimeFieldName:   data.Get("time_field_name").(string),
		Transforms:      transforms,
		FilterPredicate: filter,

		// Cacheable:         data.Get("cachable").(bool),
		//ArrayFlattenDepth: arrayFlattenCS,
	}

	if err := c.CreateView(ctx, createViewRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))

	return resourceViewRead(ctx, data, meta)

}

func resourceViewRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//when call view_by_id view_id get from here
	data.SetId(data.Get("bucket").(string))
	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).CSClient

	tokenValue := meta.(*ProviderMeta).token
	req := &client.ReadViewRequest{
		AuthToken: tokenValue,
		ID:        data.Id(),
	}

	resp, err := c.ReadView(ctx, req)
	if resp == nil {
		return diag.Errorf("Couldn't find View: %s", err)
	}

	if err != nil {
		return diag.Errorf("Failed to read View: %s", err)
	}

	data.Set("name", data.Id())
	data.Set("view_id", resp.ID)
	data.Set("_cacheable", resp.Cacheable)
	data.Set("_case_insensitive", resp.CaseInsensitive)
	data.Set("_type", resp.Type)
	data.Set("bucket", resp.Bucket)
	data.Set("index_pattern", resp.IndexPattern)
	data.Set("time_field_name", resp.TimeFieldName)

	RegionAvailability := make([]interface{}, 1)
	RegionAvailability[0] = resp.RegionAvailability
	data.Set("region_availability", RegionAvailability[0])

	if resp.MetaData != nil {
		metadata := make([]interface{}, 1)
		metadataObjectMap := make(map[string]interface{})
		metadataObjectMap["creation_date"] = resp.MetaData.CreationDate
		metadata[0] = metadataObjectMap
		data.Set("metadata", metadata)
	}

	//data.Set("resp.Public", resp.Public)
	//data.Set("resp.ContentType", resp.ContentType)
	//data.Set("resp.Type", resp.Type)

	//data.Set("Type", resp.Type)
	//data.Set("ContentType", resp.ContentType)
	//data.Set("Public", resp.Public)
	//data.Set("RealTime", resp.Realtime)
	//data.Set("Type", resp.Type)
	//data.Set("Bucket", resp.Bucket)
	//data.Set("Interval.Column", resp.Interval.Column)
	//data.Set("Interval.Mode", resp.Interval.Mode)
	//data.Set("RegionAvailability", resp.RegionAvailability)
	//data.Set("Source", resp.Source)
	//data.Set("Options", resp.Options)
	//data.Set("Metadata.CreationDate", resp.Metadata.CreationDate)
	//data.Set("Format.ColumnDelimiter", resp.Format.ColumnDelimiter)
	//data.Set("Format.HeaderRow", resp.Format.HeaderRow)
	//data.Set("Format.RowDelimiter", resp.Format.RowDelimiter)
	//data.Set("filter_json", resp.FilterJSON)
	//data.Set("format", resp.Format)
	//data.Set("live_events_sqs_arn", resp.LiveEventsSqsArn)
	//data.Set("index_retention", resp.IndexRetention)

	// When the object in an Object Group use no compression, you need to create it with
	// `compression = ""`. However, when querying an Object Group whose object are not
	// compressed, the API returns `compression = "none"`. We coerce the "none" value to
	// an empty string in order not to confuse Terraform.
	compressionOrEmptyString := resp.Compression
	if strings.ToLower(compressionOrEmptyString) == "none" {
		compressionOrEmptyString = ""
	}
	//data.Set("compression", compressionOrEmptyString)
	//
	//data.Set("partition_by", resp.PartitionBy)
	//data.Set("pattern", resp.Pattern)
	//data.Set("source_bucket", resp.SourceBucket)
	//
	//data.Set("column_selection", resp.ColumnSelection)
	//
	//// "unlimited" flattening represented as "null" in the api, and as -1 in the terraform module
	//// because the terraform sdk doesn't support nil values in configs https://github.com/hashicorp/terraform-plugin-sdk/issues/261
	//// We represent "null" as an int pointer to nil in the code.
	if resp.ArrayFlattenDepth == nil {
		data.Set("array_flatten_depth", -1)
	} else {
		data.Set("array_flatten_depth", resp.ArrayFlattenDepth)
	}
	//
	return diags
}

func resourceViewUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO to be developed
	return nil
}

func resourceViewDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderMeta).CSClient

	tokenValue := meta.(*ProviderMeta).token
	deleteViewRequest := &client.DeleteViewRequest{
		AuthToken: tokenValue,
		Name:      data.Get("bucket").(string),
	}

	if err := c.DeleteView(ctx, deleteViewRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))
	return nil
}

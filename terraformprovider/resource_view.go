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

	state := client.State{
		Type_: stateColumnSelectionInterface["_type"].(string),
	}

	pred := client.Pred{
		Field: predColumnSelectionInterface["field"].(string),
		Query: predColumnSelectionInterface["query"].(string),
		State: state,
		Type_: predColumnSelectionInterface["_type"].(string),
	}

	Predicate := client.Predicate{
		Type_: predicateColumnSelectionInterface["_type"].(string),
		Pred:  pred,
	}

	filter := &client.FilterPredicate{
		Predicate: &Predicate,
	}

	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token

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

	if resp.FilterPredicate != nil {

		filter := make([]interface{}, 1)
		predicate := make([]interface{}, 1)
		pred := make([]interface{}, 1)
		state := make([]interface{}, 1)

		predObjectMap := make(map[string]interface{})
		if resp.FilterPredicate.Predicate != nil {
			predObjectMap["field"] = resp.FilterPredicate.Predicate.Pred.Field
			predObjectMap["_type"] = resp.FilterPredicate.Predicate.Pred.Type_
			predObjectMap["query"] = resp.FilterPredicate.Predicate.Pred.Query

			stateObjectMap := make(map[string]interface{})
			stateObjectMap["_type"] = resp.FilterPredicate.Predicate.Pred.State.Type_
			state[0] = stateObjectMap
			predObjectMap["state"] = state
			pred[0] = predObjectMap

			predicatePredObjectMap := make(map[string]interface{})
			predicatePredObjectMap["pred"] = pred
			predicatePredObjectMap["_type"] = resp.FilterPredicate.Predicate.Type_
			predicate[0] = predicatePredObjectMap

			filterPredicateObjectMap := make(map[string]interface{})
			filterPredicateObjectMap["predicate"] = predicate
			filter[0] = filterPredicateObjectMap
			data.Set("filter", filter)
		}
	}

	compressionOrEmptyString := resp.Compression
	if strings.ToLower(compressionOrEmptyString) == "none" {
		compressionOrEmptyString = ""
	}

	if resp.ArrayFlattenDepth == nil {
		data.Set("array_flatten_depth", -1)
	} else {
		data.Set("array_flatten_depth", resp.ArrayFlattenDepth)
	}
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

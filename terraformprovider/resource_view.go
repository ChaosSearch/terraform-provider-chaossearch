package cs

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceViewCreate,
		ReadContext:   resourceViewRead,
		UpdateContext: resourceViewUpdate,
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

	ViewRequestDTO := setViewRequest(data, meta)

	createViewRequest := &client.CreateViewRequest{
		AuthToken: ViewRequestDTO.Token,

		Bucket:          data.Get("bucket").(string),
		Sources:         ViewRequestDTO.Source,
		IndexPattern:    data.Get("index_pattern").(string),
		Overwrite:       data.Get("overwrite").(bool),
		CaseInsensitive: data.Get("case_insensitive").(bool),
		IndexRetention:  data.Get("index_retention").(int),
		TimeFieldName:   data.Get("time_field_name").(string),
		Transforms:      ViewRequestDTO.Transforms,
		FilterPredicate: ViewRequestDTO.FilterPredicate,
	}

	if err := ViewRequestDTO.CSClient.CreateView(ctx, createViewRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))

	return resourceViewRead(ctx, data, meta)

}

type ViewRequestDTO struct {
	FilterPredicate *client.FilterPredicate
	CSClient        *client.CSClient
	Token           string
	Source          []interface{}
	Transforms      []interface{}
}

func setViewRequest(data *schema.ResourceData, meta interface{}) *ViewRequestDTO {

	filterInterface := data.Get("filter").(*schema.Set).List()[0].(map[string]interface{})
	predicateInterface := filterInterface["predicate"].(*schema.Set).List()[0].(map[string]interface{})
	predInterface := predicateInterface["pred"].(*schema.Set).List()[0].(map[string]interface{})
	stateInterface := predInterface["state"].(*schema.Set).List()[0].(map[string]interface{})

	state := client.State{
		Type: stateInterface["_type"].(string),
	}

	pred := client.Pred{
		Field: predInterface["field"].(string),
		Query: predInterface["query"].(string),
		State: state,
		Type:  predInterface["_type"].(string),
	}

	Predicate := client.Predicate{
		Type: predicateInterface["_type"].(string),
		Pred: pred,
	}

	filter := &client.FilterPredicate{
		Predicate: &Predicate,
	}

	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token

	sources, ok := data.GetOk("sources")
	if !ok {
		log.Error(" sources not available")
	}
	var sourcesStrings []interface{}

	if sources != nil {
		sourcesStrings = sources.([]interface{})
	}

	transformElem, ok := data.GetOk("transforms")
	if !ok {
		log.Error(" transforms not available")
	}
	var transforms []interface{}

	if transformElem != nil {
		transforms = transformElem.([]interface{})
	}
	ViewRequestDTO := ViewRequestDTO{
		FilterPredicate: filter,
		CSClient:        c,
		Token:           tokenValue,
		Source:          sourcesStrings,
		Transforms:      transforms,
	}
	return &ViewRequestDTO
}

func resourceViewRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	c.Set(data, "name", data.Id())
	c.Set(data, "view_id", resp.ID)
	c.Set(data, "_cacheable", resp.Cacheable)
	c.Set(data, "_case_insensitive", resp.CaseInsensitive)
	c.Set(data, "_type", resp.Type)
	c.Set(data, "bucket", resp.Bucket)
	c.Set(data, "index_pattern", resp.IndexPattern)
	c.Set(data, "time_field_name", resp.TimeFieldName)

	RegionAvailability := make([]interface{}, 1)
	RegionAvailability[0] = resp.RegionAvailability
	c.Set(data, "region_availability", RegionAvailability[0])

	if resp.MetaData != nil {
		metadata := make([]interface{}, 1)
		metadataObjectMap := make(map[string]interface{})
		metadataObjectMap["creation_date"] = resp.MetaData.CreationDate
		metadata[0] = metadataObjectMap
		c.Set(data, "metadata", metadata)
	}
	if resp.FilterPredicate != nil {
		filter := make([]interface{}, 1)
		predicate := make([]interface{}, 1)
		pred := make([]interface{}, 1)
		state := make([]interface{}, 1)

		predObjectMap := make(map[string]interface{})
		if resp.FilterPredicate.Predicate != nil {
			predObjectMap["field"] = resp.FilterPredicate.Predicate.Pred.Field
			predObjectMap["_type"] = resp.FilterPredicate.Predicate.Pred.Type
			predObjectMap["query"] = resp.FilterPredicate.Predicate.Pred.Query

			stateObjectMap := make(map[string]interface{})
			stateObjectMap["_type"] = resp.FilterPredicate.Predicate.Pred.State.Type
			state[0] = stateObjectMap
			predObjectMap["state"] = state
			pred[0] = predObjectMap

			predicatePredObjectMap := make(map[string]interface{})
			predicatePredObjectMap["pred"] = pred
			predicatePredObjectMap["_type"] = resp.FilterPredicate.Predicate.Type
			predicate[0] = predicatePredObjectMap

			filterPredicateObjectMap := make(map[string]interface{})
			filterPredicateObjectMap["predicate"] = predicate
			filter[0] = filterPredicateObjectMap
			c.Set(data, "filter", filter)
		}
	}

	if resp.ArrayFlattenDepth == nil {
		c.Set(data, "array_flatten_depth", -1)
	} else {
		c.Set(data, "array_flatten_depth", resp.ArrayFlattenDepth)
	}
	return diags
}

func resourceViewUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ViewRequestDTO := setViewRequest(data, meta)

	createViewRequest := &client.CreateViewRequest{
		AuthToken: ViewRequestDTO.Token,

		Bucket:          data.Get("bucket").(string),
		Sources:         ViewRequestDTO.Source,
		IndexPattern:    data.Get("index_pattern").(string),
		Overwrite:       true,
		CaseInsensitive: data.Get("case_insensitive").(bool),
		IndexRetention:  data.Get("index_retention").(int),
		TimeFieldName:   data.Get("time_field_name").(string),
		Transforms:      ViewRequestDTO.Transforms,
		FilterPredicate: ViewRequestDTO.FilterPredicate,
	}

	if err := ViewRequestDTO.CSClient.CreateView(ctx, createViewRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))

	return resourceViewRead(ctx, data, meta)
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

package main

import (
	"context"
	"cs-tf-provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predicate": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pred": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"_type": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"query": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"state": {
													Type:     schema.TypeSet,
													Required: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"_type": {
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
									"_type": {
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

	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value------------>>>>", tokenValue)

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

	// TODO to be developed
	// return resourceObjectGroupRead(ctx, data, meta)
	return nil
}

func resourceViewRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO to be developed
	return nil
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

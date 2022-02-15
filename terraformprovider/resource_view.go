package main

import (
	"context"
	"cs-tf-provider/client"

	// "fmt"

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
			// "pattern": {
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// 	ForceNew: false,
			// 	Optional: true,
			// },
			// "bucket": {
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// 	ForceNew: false,
			// 	Optional: true,
			// },
			// "sources": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// 	Optional: true,
			// },
			// "index_pattern": {
			// 	Type:     schema.TypeString,
			// 	Required: false,
			// 	ForceNew: false,
			// 	Optional: true,
			// },
			// "case_insensitive": {
			// 	Type:     schema.TypeBool,
			// 	Required: true,
			// },
			// "index_retention": {Type: schema.TypeInt,
			// 	Default:     14,
			// 	Description: "Number of days to keep the data before deleting it",
			// 	Optional:    true,
			// 	ForceNew:    false,
			// 	// Type:     schema.TypeSet,
			// 	// Required: false,
			// 	// Elem: &schema.Resource{
			// 	// 	Schema: map[string]*schema.Schema{
			// 	// 		"value": {
			// 	// 			Type:     schema.TypeString,
			// 	// 			Required: false,
			// 	// 			Optional: true,
			// 	// 		},
			// 	// 	},
			// 	// },
			// },
			// "filter_json": {
			// 	Type:         schema.TypeString,
			// 	Default:      `[{"field":"key","regex":".*"}]`,
			// 	Optional:     true,
			// 	ForceNew:     true,
			// 	ValidateFunc: validation.StringIsJSON,
			// },
			// // "time_field_name": {
			// // 	Type:     schema.TypeString,
			// // 	Required: false,
			// // 	ForceNew: false,
			// // },
			// "cachable": {
			// 	Type:     schema.TypeBool,
			// 	Required: false,
			// 	ForceNew: false,
			// 	Optional: true,
			// },
			// "overwrite": {
			// 	Type:     schema.TypeBool,
			// 	Required: false,
			// 	ForceNew: false,
			// 	Default:  false,
			// 	Optional: true,
			// },
			// "transforms": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// 	Optional: true,
			// },
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

	log.Info("predicateColumnSelectionInterface===", predicateColumnSelectionInterface)
	log.Info("stateColumnSelectionInterface===", stateColumnSelectionInterface)

	log.Info("predicateColumnSelectionInterface _type ]===", predicateColumnSelectionInterface["_type"].(string))
	log.Info("predicateColumnSelectionInterface pred ]===", predicateColumnSelectionInterface["pred"].(*schema.Set))

	log.Info("predColumnSelectionInterface field ]===", predColumnSelectionInterface["field"].(string))
	log.Info("predColumnSelectionInterface query ]===", predColumnSelectionInterface["query"].(string))
	log.Info("predColumnSelectionInterface state ]===", predColumnSelectionInterface["state"].(*schema.Set))
	log.Info("predColumnSelectionInterface _type ]===", predColumnSelectionInterface["_type"].(string))

	log.Info("stateColumnSelectionInterface _type ]===", stateColumnSelectionInterface["_type"].(string))

	var predicateType string
	// var pred_value *schema.Set

	var field_value string
	var query_value string
	// var state_value *schema.Set
	var _type_value string

	var stateType string

	predicateType = predicateColumnSelectionInterface["_type"].(string)
	// pred_value = predicateColumnSelectionInterface["pred"].(*schema.Set)

	field_value = predColumnSelectionInterface["field"].(string)
	query_value = predColumnSelectionInterface["query"].(string)
	// state_value = predColumnSelectionInterface["state"].(*schema.Set)
	_type_value = predColumnSelectionInterface["_type"].(string)

	stateType=stateColumnSelectionInterface["_type"].(string)

	state := client.State{
		Type_: stateType,
	}

	pred := client.Pred{
		Field: field_value,
		Query: query_value,
		State: state,
		Type_: _type_value,
	}

	
	Predicate := client.Predicate{
		Type_: predicateType,
		Pred:  pred,
	}
	
	filter := &client.Filter{
		&Predicate,
	}

	// filter := client.Filter{
	// 	Predicate: filterColumnSelectionInterface["predicate"].(*schema.Set),
	// }

	// predicate := client.Predicate{
	// 	Field: predicateColumnSelectionInterface["field"].(string),
	// 	Query: predicateColumnSelectionInterface["query"].(string),
	// 	State: predicateColumnSelectionInterface["state"].(*client.State),
	// 	Type_: predicateColumnSelectionInterface["_type"].(string),
	// }

	// state := client.State{
	// 	Type_: predicateColumnSelectionInterface["_type"].(string),
	// }

	c := meta.(*ProviderMeta).Client
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
		log.Debug("sourcesStrings-->", sourcesStrings)
	}

	log.Debug("sourcesStrings-->", sourcesStrings)

	transforms_, ok := data.GetOk("transforms")
	if !ok {
		log.Error(" transforms not available")
	}
	var transforms []interface{}

	if transforms_ != nil {
		transforms = transforms_.([]interface{})
	}

	// patterns_, ok := data.GetOk("pattern")
	// if !ok {
	// 	log.Error(" sources not available")
	// }
	// patterns := patterns_.([]interface{})
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
		 Filter:         filter,

		// Cacheable:         data.Get("cachable").(bool),
		// ArrayFlattenDepth: arrayFlattenCS,

	}

	log.Info("createViewRequest.Bucket--->", createViewRequest.Bucket)
	log.Info("createViewRequest.TimeFieldName--->", createViewRequest.TimeFieldName)
	// log.Info("createViewRequest.Pattern--->", createViewRequest.Pattern)
	log.Info("createViewRequest.IndexRetention--->", createViewRequest.IndexRetention)
	log.Info("createViewRequest.Cacheable--->", createViewRequest.Cacheable)
	// log.Info("createViewRequest.Overwrite--->", createViewRequest.Overwrite)

	if err := c.CreateView(ctx, createViewRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("bucket").(string))

	// return resourceObjectGroupRead(ctx, data, meta)
	return nil
}

func resourceViewRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceViewUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceViewDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

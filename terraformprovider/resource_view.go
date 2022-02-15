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
									"field": {
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
	stateColumnSelectionInterface := predicateColumnSelectionInterface["state"].(*schema.Set).List()[0].(map[string]interface{})

	log.Info("predicateColumnSelectionInterface===", predicateColumnSelectionInterface)
	log.Info("stateColumnSelectionInterface===", stateColumnSelectionInterface)

	filter := client.Filter{
		Predicate: filterColumnSelectionInterface["predicate"].(*client.Predicate),
	}


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

		Bucket:  data.Get("bucket").(string),
		Sources: sourcesStrings,
		IndexPattern:   data.Get("index_pattern").(string),
		Overwrite:       data.Get("overwrite").(bool),
		CaseInsensitive: data.Get("case_insensitive").(bool),
		IndexRetention:  data.Get("index_retention").(int),
		TimeFieldName:   data.Get("time_field_name").(string),
		Transforms: transforms,
		Filter: &filter,

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

package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ViewRequestDTO struct {
	FilterPredicate *client.FilterPredicate
	CSClient        *client.CSClient
	Token           string
	Source          []interface{}
	Transforms      []interface{}
}

func ResourceView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceViewCreate,
		ReadContext:   ResourceViewRead,
		UpdateContext: resourceViewUpdate,
		DeleteContext: resourceViewDelete,
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
				Optional: true,
			},
			"cacheable": {
				Type:     schema.TypeBool,
				ForceNew: false,
				Optional: true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"case_insensitive": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"region_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"index_retention": {
				Type:        schema.TypeInt,
				Default:     14,
				Description: "",
				Optional:    true,
			},
			"time_field_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transforms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
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
			"array_flatten_depth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predicate": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pred": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"query": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"state": {
													Type:     schema.TypeSet,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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
		AuthToken:       ViewRequestDTO.Token,
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

	return ResourceViewRead(ctx, data, meta)

}

func setViewRequest(data *schema.ResourceData, meta interface{}) *ViewRequestDTO {
	var sourcesStrings []interface{}
	var transforms []interface{}
	var state *client.State
	var pred *client.Pred
	var predicate *client.Predicate
	var filter *client.FilterPredicate

	filterList := data.Get("filter").(*schema.Set).List()
	if len(filterList) > 0 {
		filterMap := filterList[0].(map[string]interface{})
		predicateList := filterMap["predicate"].(*schema.Set).List()
		if len(predicateList) > 0 {
			predicateMap := predicateList[0].(map[string]interface{})
			predList := predicateMap["pred"].(*schema.Set).List()
			if len(predList) > 0 {
				predMap := predList[0].(map[string]interface{})
				stateList := predMap["state"].(*schema.Set).List()
				if len(stateList) > 0 {
					stateMap := stateList[0].(map[string]interface{})
					state = &client.State{
						Type: stateMap["type"].(string),
					}
				}

				pred = &client.Pred{
					Field: predMap["field"].(string),
					Query: predMap["query"].(string),
					State: *state,
					Type:  predMap["type"].(string),
				}
			}

			predicate = &client.Predicate{
				Type: predicateMap["type"].(string),
				Pred: *pred,
			}
		}

		filter = &client.FilterPredicate{
			Predicate: predicate,
		}
	}

	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	sources, _ := data.GetOk("sources")
	if sources != nil {
		sourcesStrings = sources.([]interface{})
	}

	transformElem, _ := data.GetOk("transforms")
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

func ResourceViewRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("bucket").(string))
	diags := diag.Diagnostics{}
	c := meta.(*models.ProviderMeta).CSClient

	tokenValue := meta.(*models.ProviderMeta).Token
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

	if resp.MetaData != nil {
		err = data.Set("metadata", []interface{}{
			map[string]interface{}{
				"creation_date": resp.MetaData.CreationDate,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if resp.FilterPredicate != nil {
		if resp.FilterPredicate.Predicate != nil {
			err = data.Set("filter", []interface{}{
				map[string]interface{}{
					"predicate": []interface{}{
						map[string]interface{}{
							"type": resp.FilterPredicate.Predicate.Type,
							"pred": []interface{}{
								map[string]interface{}{
									"field": resp.FilterPredicate.Predicate.Pred.Field,
									"type":  resp.FilterPredicate.Predicate.Pred.Type,
									"query": resp.FilterPredicate.Predicate.Pred.Query,
									"state": []interface{}{
										map[string]interface{}{
											"type": resp.FilterPredicate.Predicate.Pred.State.Type,
										},
									},
								},
							},
						},
					},
				},
			})

			if err != nil {
				return diag.FromErr(err)
			}
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

	err = data.Set("id", resp.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("cacheable", resp.Cacheable)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("case_insensitive", resp.CaseInsensitive)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("type", resp.Type)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("bucket", resp.Bucket)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("index_pattern", resp.IndexPattern)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("time_field_name", resp.TimeFieldName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("region_availability", resp.RegionAvailability)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceViewUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ViewRequestDTO := setViewRequest(data, meta)

	createViewRequest := &client.CreateViewRequest{
		AuthToken:       ViewRequestDTO.Token,
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

	return ResourceViewRead(ctx, data, meta)
}

func resourceViewDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient

	tokenValue := meta.(*models.ProviderMeta).Token
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

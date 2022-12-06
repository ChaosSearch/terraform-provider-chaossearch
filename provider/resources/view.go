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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ViewData struct {
	FilterPredicate *client.FilterPredicate
	Source          []interface{}
	Transforms      []client.Transform
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
				Default:  false,
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
				Computed: true,
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
				Type:     schema.TypeInt,
				Default:  -1,
				Optional: true,
			},
			"time_field_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transforms": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predicate": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"preds": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsJSON,
										},
									},
									"pred": {
										Type:     schema.TypeSet,
										Optional: true,
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
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token

	viewData, err := setViewRequest(data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createViewRequest := &client.CreateViewRequest{
		AuthToken:       tokenValue,
		Bucket:          data.Get("bucket").(string),
		Sources:         viewData.Source,
		IndexPattern:    data.Get("index_pattern").(string),
		Overwrite:       data.Get("overwrite").(bool),
		CaseInsensitive: data.Get("case_insensitive").(bool),
		IndexRetention:  data.Get("index_retention").(int),
		TimeFieldName:   data.Get("time_field_name").(string),
		Cacheable:       data.Get("cacheable").(bool),
		Transforms:      viewData.Transforms,
		FilterPredicate: viewData.FilterPredicate,
	}

	if err := c.CreateView(ctx, createViewRequest); err != nil {
		return diag.FromErr(err)
	}

	return ResourceViewRead(ctx, data, meta)

}

func setViewRequest(data *schema.ResourceData, meta interface{}) (*ViewData, error) {
	var sourcesStrings []interface{}
	var transforms []client.Transform
	var state *client.State
	var pred *client.Pred
	var preds []client.Pred
	var predicate *client.Predicate
	var filter *client.FilterPredicate

	filterList := data.Get("filter").(*schema.Set).List()
	if len(filterList) > 0 {
		filterMap := filterList[0].(map[string]interface{})
		predicateList := filterMap["predicate"].(*schema.Set).List()
		if len(predicateList) > 0 {
			predicateMap := predicateList[0].(map[string]interface{})
			predList := predicateMap["pred"].(*schema.Set).List()
			predsArr := predicateMap["preds"].([]interface{})

			if (len(predsArr) == 0 && len(predList) == 0) || (len(predsArr) > 0 && len(predList) > 0) {
				return nil, fmt.Errorf("Either 'preds' or 'pred' need to be defined (Not both or none).")
			}

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
					State: state,
					Type:  predMap["type"].(string),
				}
			}

			if len(predsArr) > 0 {
				preds = make([]client.Pred, len(predsArr))
				for index, predItem := range predsArr {
					var predHolder client.Pred

					predJson, err := structure.NormalizeJsonString(predItem)
					if err != nil {
						return nil, utils.NormalizingJsonError(err)
					}

					err = json.Unmarshal([]byte(predJson), &predHolder)
					if err != nil {
						return nil, utils.UnmarshalJsonError(err)
					}

					preds[index] = predHolder
				}
			}

			predicate = &client.Predicate{
				Type:  predicateMap["type"].(string),
				Pred:  pred,
				Preds: preds,
			}
		}

		filter = &client.FilterPredicate{
			Predicate: predicate,
		}
	}

	sources, _ := data.GetOk("sources")
	if sources != nil {
		sourcesStrings = sources.([]interface{})
	}

	transformsArr := data.Get("transforms").([]interface{})
	if len(transformsArr) > 0 {
		transforms = make([]client.Transform, len(transformsArr))
		for index, transform := range transformsArr {
			var transformHolder client.Transform
			transformJson, err := structure.NormalizeJsonString(transform)
			if err != nil {
				return nil, utils.NormalizingJsonError(err)
			}

			err = json.Unmarshal([]byte(transformJson), &transformHolder)
			if err != nil {
				return nil, utils.UnmarshalJsonError(err)
			}

			transforms[index] = transformHolder
		}
	}

	return &ViewData{
		FilterPredicate: filter,
		Source:          sourcesStrings,
		Transforms:      transforms,
	}, nil
}

func ResourceViewRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	c := meta.(*models.ProviderMeta).CSClient

	tokenValue := meta.(*models.ProviderMeta).Token
	req := &client.BasicRequest{
		AuthToken: tokenValue,
		Id:        data.Get("bucket").(string),
	}

	resp, err := c.ReadView(ctx, req)
	if resp == nil {
		return diag.Errorf("Couldn't find View: %s", err)
	}

	if err != nil {
		return diag.Errorf("Failed to read View: %s", err)
	}

	data.SetId(resp.ID)
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
		predicate := resp.FilterPredicate.Predicate
		if predicate != nil {
			var preds []string

			predicateMap := map[string]interface{}{
				"type": predicate.Type,
			}

			if predicate.Pred != nil {
				predicateMap["pred"] = []interface{}{
					map[string]interface{}{
						"field": predicate.Pred.Field,
						"type":  predicate.Pred.Type,
						"query": predicate.Pred.Query,
						"state": []interface{}{
							map[string]interface{}{
								"type": predicate.Pred.State.Type,
							},
						},
					},
				}
			}

			if predicate.Preds != nil {
				preds = make([]string, len(predicate.Preds))
				for index, predItem := range predicate.Preds {
					predMap := map[string]interface{}{
						"field": predItem.Field,
						"_type": predItem.Type,
						"query": predItem.Query,
						"state": map[string]interface{}{
							"_type": predItem.State.Type,
						},
					}

					predBytes, err := json.Marshal(predMap)
					if err != nil {
						return diag.FromErr(utils.MarshalJsonError(err))
					}

					predJson, err := structure.NormalizeJsonString(string(predBytes))
					if err != nil {
						return diag.FromErr(utils.NormalizingJsonError(err))
					}

					preds[index] = predJson
				}

				predicateMap["preds"] = preds
			}

			err = data.Set("filter", []interface{}{
				map[string]interface{}{
					"predicate": []interface{}{
						predicateMap,
					},
				},
			})

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if len(resp.Transforms) > 0 {
		transforms := make([]string, len(resp.Transforms))
		for index, transformItem := range resp.Transforms {
			transformMap := map[string]interface{}{
				"_type":      transformItem.Type,
				"inputField": transformItem.InputField,
				"keyPart":    transformItem.KeyPart,
				"pattern":    transformItem.Pattern,
				"paths":      transformItem.Paths,
				"vertical":   transformItem.Vertical,
				"format":     transformItem.Format,
			}

			transformBytes, err := json.Marshal(transformMap)
			if err != nil {
				return diag.FromErr(utils.MarshalJsonError(err))
			}

			transformJson, err := structure.NormalizeJsonString(string(transformBytes))
			if err != nil {
				return diag.FromErr(utils.NormalizingJsonError(err))
			}

			transforms[index] = transformJson
		}

		err = data.Set("transforms", transforms)
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
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token

	viewData, err := setViewRequest(data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createViewRequest := &client.CreateViewRequest{
		AuthToken:       tokenValue,
		Bucket:          data.Get("bucket").(string),
		Sources:         viewData.Source,
		IndexPattern:    data.Get("index_pattern").(string),
		Overwrite:       true,
		CaseInsensitive: data.Get("case_insensitive").(bool),
		IndexRetention:  data.Get("index_retention").(int),
		TimeFieldName:   data.Get("time_field_name").(string),
		Transforms:      viewData.Transforms,
		FilterPredicate: viewData.FilterPredicate,
	}

	if err := c.CreateView(ctx, createViewRequest); err != nil {
		return diag.FromErr(err)
	}

	return ResourceViewRead(ctx, data, meta)
}

func resourceViewDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	deleteViewRequest := &client.BasicRequest{
		AuthToken: tokenValue,
		Id:        data.Get("bucket").(string),
	}

	if err := c.DeleteView(ctx, deleteViewRequest); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

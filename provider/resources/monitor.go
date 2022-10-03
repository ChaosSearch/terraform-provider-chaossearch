package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/models"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorCreate,
		ReadContext:   resourceMonitorRead,
		UpdateContext: resourceMonitorUpdate,
		DeleteContext: resourceMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "monitor",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"schedule": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"unit": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"inputs": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"indices": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"query": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsJSON,
									},
								},
							},
						},
					},
				},
			},
			"triggers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"condition": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"script": &scriptSchema,
								},
							},
						},
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"destination_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"message_template": &scriptSchema,
									"subject_template": &scriptSchema,
									"throttle_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"throttle": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"unit": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
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

var scriptSchema schema.Schema = schema.Schema{
	Type:     schema.TypeSet,
	Required: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lang": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	},
}

func resourceMonitorCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	req, diagErr := constructCreateMonitorRequest(data, meta)
	if diagErr != nil {
		return diagErr
	}

	resp, err := c.CreateMonitor(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Resp.Id)
	return nil
}

func resourceMonitorRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceMonitorUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	req, diagErr := constructCreateMonitorRequest(data, meta)
	if diagErr != nil {
		return diagErr
	}

	req.Id = data.Id()
	err := c.UpdateMonitor(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMonitorDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := c.DeleteMonitor(ctx, &client.BasicRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		Id:        data.Id(),
	})

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func extractScriptValues(scriptMap map[string]interface{}) client.Script {
	return client.Script{
		Lang:   scriptMap["lang"].(string),
		Source: scriptMap["source"].(string),
	}
}

func constructCreateMonitorRequest(
	data *schema.ResourceData,
	meta interface{},
) (*client.CreateMonitorRequest, diag.Diagnostics) {
	var triggers []client.Trigger
	var actions []client.Action
	var indices []string

	scheduleMap := utils.ConvertSetToMap(data.Get("schedule"))
	periodMap := utils.ConvertSetToMap(scheduleMap["period"])
	inputMap := data.Get("inputs").([]interface{})[0].(map[string]interface{})
	searchMap := utils.ConvertSetToMap(inputMap["search"])
	for _, index := range searchMap["indices"].([]interface{}) {
		indices = append(indices, index.(string))
	}
	search := client.Search{
		Indices: indices,
	}

	err := json.Unmarshal([]byte(searchMap["query"].(string)), &search.Query)
	if err != nil {
		return nil, diag.FromErr(utils.UnmarshalJsonError(err))
	}

	triggersList := data.Get("triggers").([]interface{})
	for _, triggerItem := range triggersList {
		triggerMap := triggerItem.(map[string]interface{})
		actionsList := triggerMap["actions"].([]interface{})
		for _, actionItem := range actionsList {
			actionMap := actionItem.(map[string]interface{})
			msgMap := utils.ConvertSetToMap(actionMap["message_template"])
			subjectMap := utils.ConvertSetToMap(actionMap["subject_template"])
			action := client.Action{
				DestinationId:   actionMap["destination_id"].(string),
				Name:            actionMap["name"].(string),
				SubjectTemplate: extractScriptValues(msgMap),
				MessageTemplate: extractScriptValues(subjectMap),
			}

			if actionMap["throttle_enabled"].(bool) {
				throttleMap := utils.ConvertSetToMap(actionMap["throttle"])
				action.ThrottleEnabled = true
				action.Throttle = client.Throttle{
					Value: throttleMap["value"].(int),
					Unit:  throttleMap["unit"].(string),
				}
			}

			actions = append(actions, action)
		}
		conditionMap := utils.ConvertSetToMap(triggerMap["condition"])
		scriptMap := utils.ConvertSetToMap(conditionMap["script"])
		triggers = append(triggers, client.Trigger{
			Name:     triggerMap["name"].(string),
			Severity: triggerMap["severity"].(string),
			Condition: client.MonitorCondition{
				Script: extractScriptValues(scriptMap),
			},
			Actions: actions,
		})
	}

	return &client.CreateMonitorRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		Name:      data.Get("name").(string),
		Type:      data.Get("type").(string),
		Enabled:   data.Get("enabled").(bool),
		Schedule: client.Schedule{
			Period: client.Period{
				Interval: periodMap["interval"].(int),
				Unit:     periodMap["unit"].(string),
			},
		},
		Inputs: []client.Input{{
			Search: search,
		}},
		Triggers: triggers,
	}, nil
}

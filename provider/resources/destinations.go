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

const (
	Name          = "name"
	Type          = "type"
	CustomWebhook = "custom_webhook"
	Scheme        = "scheme"
	Host          = "host"
	Url           = "url"
	Path          = "path"
	Port          = "port"
	QueryParams   = "query_params"
	HeaderParams  = "header_params"
	Method        = "method"
	Slack         = "slack"
)

func ResourceDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDestinationCreate,
		ReadContext:   resourceDestinationRead,
		UpdateContext: resourceDestinationUpdate,
		DeleteContext: resourceDestinationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			Name: {
				Type:     schema.TypeString,
				Required: true,
			},
			Type: {
				Type:     schema.TypeString,
				Required: true,
			},
			CustomWebhook: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Scheme: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "HTTPS",
						},
						Host: {
							Type:     schema.TypeString,
							Optional: true,
						},
						Url: {
							Type:     schema.TypeString,
							Optional: true,
						},
						Path: {
							Type:     schema.TypeString,
							Optional: true,
						},
						Port: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						QueryParams: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Default:      `{"":""}`,
						},
						HeaderParams: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Default:      `{"Content-Type": "application/json"}`,
						},
						Method: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "POST",
						},
					},
				},
			},
			Slack: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Url: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDestinationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	req, diagErr := constructCreateDestinationRequest(data, meta)
	if diagErr != nil {
		return diagErr
	}

	destinationResp, err := c.CreateDestination(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(destinationResp.Id)
	return nil
}

func resourceDestinationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	resp, err := c.ReadDestination(ctx, &client.BasicRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		Id:        data.Id(),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	for _, dest := range resp.Destinations {
		if dest.Id == data.Id() {
			err := data.Set(Name, dest.Name)
			if err != nil {
				return diag.FromErr(err)
			}

			err = data.Set(Type, dest.Type)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceDestinationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	req, diagErr := constructCreateDestinationRequest(data, meta)
	if diagErr != nil {
		return diagErr
	}

	req.Id = data.Id()
	err := c.UpdateDestination(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDestinationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := c.DeleteDestination(ctx, client.BasicRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		Id:        data.Id(),
	})

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func constructCreateDestinationRequest(
	data *schema.ResourceData,
	meta interface{},
) (*client.CreateDestinationRequest, diag.Diagnostics) {
	req := client.CreateDestinationRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		Name:      data.Get(Name).(string),
		Type:      data.Get(Type).(string),
	}

	slackList := data.Get(Slack).(*schema.Set).List()
	if len(slackList) > 0 {
		slackMap := slackList[0].(map[string]interface{})
		req.Slack = &client.Slack{
			Url: slackMap[Url].(string),
		}
	}

	customList := data.Get(CustomWebhook).(*schema.Set).List()
	if len(customList) > 0 {
		customMap := customList[0].(map[string]interface{})
		if customMap[Url].(string) != "" {
			req.CustomWebhook = &client.CustomWebhook{
				Url:    customMap[Url].(string),
				Scheme: customMap[Scheme].(string),
				Method: customMap[Method].(string),
				Port:   customMap[Port].(int),
			}
		} else {
			req.CustomWebhook = &client.CustomWebhook{
				Scheme: customMap[Scheme].(string),
				Method: customMap[Method].(string),
				Host:   customMap[Host].(string),
				Port:   customMap[Port].(int),
				Path:   customMap[Path].(string),
			}
		}

		queryParamString := customMap[QueryParams].(string)
		if queryParamString != "" {
			err := json.Unmarshal([]byte(queryParamString), &req.CustomWebhook.QueryParams)
			if err != nil {
				return nil, diag.FromErr(utils.UnmarshalJsonError(err))
			}
		}

		headerParamString := customMap[HeaderParams].(string)
		if headerParamString != "" {
			err := json.Unmarshal([]byte(headerParamString), &req.CustomWebhook.HeaderParams)
			if err != nil {
				return nil, diag.FromErr(utils.UnmarshalJsonError(err))
			}
		}
	}

	return &req, nil
}

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_webhook": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "HTTPS",
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"query_params": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Default:      `{"":""}`,
						},
						"header_params": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Default:      `{"Content-Type": "application/json"}`,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "POST",
						},
					},
				},
			},
			"slack": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
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
			err := data.Set("name", dest.Name)
			if err != nil {
				return diag.FromErr(err)
			}

			err = data.Set("type", dest.Type)
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
		Name:      data.Get("name").(string),
		Type:      data.Get("type").(string),
	}

	slackList := data.Get("slack").(*schema.Set).List()
	if len(slackList) > 0 {
		slackMap := slackList[0].(map[string]interface{})
		req.Slack = &client.Slack{
			Url: slackMap["url"].(string),
		}
	}

	customList := data.Get("custom_webhook").(*schema.Set).List()
	if len(customList) > 0 {
		customMap := customList[0].(map[string]interface{})
		if customMap["url"].(string) != "" {
			req.CustomWebhook = &client.CustomWebhook{
				Url:    customMap["url"].(string),
				Scheme: customMap["scheme"].(string),
				Method: customMap["method"].(string),
				Port:   customMap["port"].(int),
			}
		} else {
			req.CustomWebhook = &client.CustomWebhook{
				Scheme: customMap["scheme"].(string),
				Method: customMap["method"].(string),
				Host:   customMap["host"].(string),
				Port:   customMap["port"].(int),
				Path:   customMap["path"].(string),
			}
		}

		queryParamString := customMap["query_params"].(string)
		if queryParamString != "" {
			err := json.Unmarshal([]byte(queryParamString), &req.CustomWebhook.QueryParams)
			if err != nil {
				return nil, diag.FromErr(utils.UnmarshalJsonError(err))
			}
		}

		headerParamString := customMap["header_params"].(string)
		if headerParamString != "" {
			err := json.Unmarshal([]byte(headerParamString), &req.CustomWebhook.HeaderParams)
			if err != nil {
				return nil, diag.FromErr(utils.UnmarshalJsonError(err))
			}
		}
	}

	return &req, nil
}

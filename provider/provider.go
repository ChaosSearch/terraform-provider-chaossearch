package provider

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/datasources"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_access_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						}, "parent_user_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"chaossearch_object_group": resources.ResourceObjectGroup(),
			"chaossearch_view":         resources.ResourceView(),
			"chaossearch_sub_account":  resources.ResourceSubAccount(),
			"chaossearch_user_group":   resources.ResourceUserGroup(),
			"chaossearch_index_model":  resources.ResourceIndexModel(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"chaossearch_retrieve_object_groups": datasources.DataSourceObjectGroups(),
			"chaossearch_retrieve_object_group":  datasources.DataSourceObjectGroup(),
			"chaossearch_retrieve_views":         datasources.DataSourceViews(),
			"chaossearch_retrieve_view":          datasources.DataSourceView(),
			"chaossearch_retrieve_sub_accounts":  datasources.DataSourceSubAccounts(),
			"chaossearch_retrieve_groups":        datasources.DataSourceUserGroups(),
			"chaossearch_retrieve_user_group":    datasources.DataSourceUserGroup(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	accessKeyID := d.Get("access_key_id").(string)
	secretAccessKey := d.Get("secret_access_key").(string)
	region := d.Get("region").(string)

	var username string
	var password string
	var parentUserID string

	if d.Get("login").(*schema.Set).Len() > 0 {
		columnSelectionInterfaces := d.Get("login").(*schema.Set).List()[0]
		columnSelectionInterface := columnSelectionInterfaces.(map[string]interface{})
		username = columnSelectionInterface["user_name"].(string)
		password = columnSelectionInterface["password"].(string)
		if columnSelectionInterface["parent_user_id"] != nil {
			parentUserID = columnSelectionInterface["parent_user_id"].(string)
		}
	}

	login := client.Login{
		Username:     username,
		Password:     password,
		ParentUserID: parentUserID,
	}

	if url == "" {
		return nil, diag.Errorf("Expected 'url' to be defined in provider configuration, but it was not")
	}
	if accessKeyID == "" {
		return nil, diag.Errorf("Expected 'access_key_id' to be defined in provider configuration, but it was not")
	}
	if secretAccessKey == "" {
		return nil, diag.Errorf("Expected 'secret_access_key' to be defined in provider configuration, but it was not")
	}
	if region == "" {
		return nil, diag.Errorf("Expected 'region' to be defined in provider configuration, but it was not")
	}

	config := client.NewConfiguration()
	config.URL = url
	config.AccessKeyID = accessKeyID
	config.SecretAccessKey = secretAccessKey
	config.Region = region

	csClient := client.NewClient(config, &login)

	authResponseString, err := csClient.Auth(ctx)

	if err != nil {
		return nil, diag.Errorf("Token generation fail..")

	}
	tokenData := models.AuthResponse{}

	if err := json.Unmarshal([]byte(authResponseString), &tokenData); err != nil {
		return fmt.Errorf("Provider Config Failure => %s", utils.UnmarshalJsonError(err)), nil
	}

	providerMeta := &models.ProviderMeta{
		CSClient: csClient,
		Token:    tokenData.Token,
	}
	return providerMeta, nil

}

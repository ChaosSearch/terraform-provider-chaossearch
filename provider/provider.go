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
						},
						"parent_user_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Required: true,
				ForceNew: true,
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
	var username string
	var password string
	var parentUserID string
	var authResp models.AuthResponse

	url, diagErr := getConfig(ctx, d, "url")
	if diagErr != nil {
		return nil, diagErr
	}

	accessKeyID, diagErr := getConfig(ctx, d, "access_key_id")
	if diagErr != nil {
		return nil, diagErr
	}

	secretAccessKey, diagErr := getConfig(ctx, d, "secret_access_key")
	if diagErr != nil {
		return nil, diagErr
	}

	region, diagErr := getConfig(ctx, d, "region")
	if diagErr != nil {
		return nil, diagErr
	}

	loginList := d.Get("login").(*schema.Set).List()
	if len(loginList) > 0 {
		loginMap := loginList[0].(map[string]interface{})
		username = loginMap["user_name"].(string)
		if username == "" {
			return nil, diag.Errorf("Failed to configure provider => Expected 'user_name' to be defined")
		}

		password = loginMap["password"].(string)
		if password == "" {
			return nil, diag.Errorf("Failed to configure provider => Expected 'password' to be defined")
		}

		if loginMap["parent_user_id"] != nil {
			parentUserID = loginMap["parent_user_id"].(string)
		}
	}

	login := client.Login{
		Username:     username,
		Password:     password,
		ParentUserID: parentUserID,
	}

	config := client.NewConfiguration(url, accessKeyID, secretAccessKey, region)
	csClient := client.NewClient(config, &login)
	authResponseString, err := csClient.Auth(ctx)
	if err != nil {
		return nil, diag.Errorf("Failed to configure provider => %s", err)

	}

	if err := json.Unmarshal([]byte(authResponseString), &authResp); err != nil {
		return fmt.Errorf("Failed to configure provider => %s", utils.UnmarshalJsonError(err)), nil
	}

	providerMeta := &models.ProviderMeta{
		CSClient: csClient,
		Token:    authResp.Token,
	}

	return providerMeta, nil

}

func getConfig(ctx context.Context, data *schema.ResourceData, value string) (string, diag.Diagnostics) {
	cred := data.Get(value).(string)
	if cred == "" {
		return "", diag.Errorf("Failed to configure provider => Expected '%s' to be defined", value)
	}

	return cred, nil
}

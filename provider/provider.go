package provider

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/datasources"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_URL", ""),
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_ACCESS_KEY", ""),
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_SECRET_KEY", ""),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_REGION", ""),
			},
			"parent_user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_PARENT_USER_ID", ""),
			},
			"login": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							DefaultFunc: schema.EnvDefaultFunc("CS_USERNAME", ""),
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							DefaultFunc: schema.EnvDefaultFunc("CS_PASSWORD", ""),
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"chaossearch_object_group": resources.ResourceObjectGroup(),
			"chaossearch_view":         resources.ResourceView(),
			"chaossearch_sub_account":  resources.ResourceSubAccount(),
			"chaossearch_user_group":   resources.ResourceUserGroup(),
			"chaossearch_index_model":  resources.ResourceIndexModel(),
			"chaossearch_destination":  resources.ResourceDestination(),
			"chaossearch_monitor":      resources.ResourceMonitor(),
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
	var authResp models.AuthResponse
	keyAuthEnabled := true

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

	parentUserID := d.Get("parent_user_id").(string)
	loginList := d.Get("login").(*schema.Set).List()
	if len(loginList) > 0 {
		loginMap := loginList[0].(map[string]interface{})
		username = loginMap["user_name"].(string)
		if username == "" {
			return nil, utils.ConfigurationError("user_name")
		}

		password = loginMap["password"].(string)
		if password == "" {
			return nil, utils.ConfigurationError("password")
		}

		keyAuthEnabled = false
	}

	if keyAuthEnabled && parentUserID == "" {
		return nil, diag.Errorf(`
		Failed to configure provider => 'parent_user_id' must be defined for API Key Auth
		Note: 
			This can be populated with the account's UID, regardless of being a subaccount
		`)
	}

	config := &client.Configuration{
		URL:             url,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		AWSServiceName:  "s3",
		Region:          region,
		UserID:          parentUserID,
		KeyAuthEnabled:  keyAuthEnabled,
	}

	login := client.Login{
		Username:     username,
		Password:     password,
		ParentUserID: parentUserID,
	}

	csClient := client.NewClient(config, &login)
	authResponseString, err := csClient.Auth(ctx)
	if err != nil {
		return nil, diag.Errorf("Failed to configure provider => %s", err)

	}

	if !keyAuthEnabled {
		if err := json.Unmarshal([]byte(authResponseString), &authResp); err != nil {
			return nil, diag.Errorf("Failed to configure provider => %s", utils.UnmarshalJsonError(err))
		}
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
		return "", utils.ConfigurationError(value)
	}

	return cred, nil
}

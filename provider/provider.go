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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	Url             = "url"
	AccessKeyID     = "access_key_id"
	SecretAccessKey = "secret_access_key"
	ParentUserID    = "parent_user_id"
	Region          = "region"
	Login           = "login"
	Username        = "user_name"
	Password        = "password"
	Options         = "options"
	RetryCount      = "retry_count"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			Url: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("CS_URL", ""),
				ValidateFunc: validateUrl,
			},
			AccessKeyID: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_ACCESS_KEY", ""),
			},
			SecretAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_SECRET_KEY", ""),
			},
			Region: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_REGION", ""),
			},
			ParentUserID: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CS_PARENT_USER_ID", ""),
			},
			Login: {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Username: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							DefaultFunc: schema.EnvDefaultFunc("CS_USERNAME", ""),
						},
						Password: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							DefaultFunc: schema.EnvDefaultFunc("CS_PASSWORD", ""),
						},
					},
				},
			},
			Options: {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RetryCount: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  10,
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
	var username, password string
	var authResp models.AuthResponse
	keyAuthMap := map[string]string{}
	keyAuthEnabled := true

	url, diagErr := getConfig(ctx, d, Url)
	if diagErr != nil {
		return nil, diagErr
	}

	region, diagErr := getConfig(ctx, d, Region)
	if diagErr != nil {
		return nil, diagErr
	}

	loginList := d.Get(Login).(*schema.Set).List()
	if len(loginList) > 0 {
		loginMap := loginList[0].(map[string]interface{})
		username = loginMap[Username].(string)
		if username == "" {
			return nil, diag.FromErr(utils.UndefinedError(Username))
		}

		password = loginMap[Password].(string)
		if password == "" {
			return nil, diag.FromErr(utils.UndefinedError(Password))
		}

		keyAuthEnabled = false
	}

	keyAuthMap[ParentUserID] = d.Get(ParentUserID).(string)
	if keyAuthEnabled {
		keyAuthMap[AccessKeyID] = d.Get(AccessKeyID).(string)
		keyAuthMap[SecretAccessKey] = d.Get(SecretAccessKey).(string)

		for key, val := range keyAuthMap {
			if key == ParentUserID && val == "" {
				return nil, utils.ProviderConfigurationError(fmt.Errorf(
					"'parent_user_id' must be defined for API Key Auth => \n" +
						"Note: \n" +
						"This can be populated with the account's UID, regardless of being a subaccount",
				))
			} else if val == "" {
				return nil, utils.ProviderConfigurationError(utils.UndefinedError(val))
			}
		}
	}

	config := &client.Configuration{
		URL:             url,
		AccessKeyID:     keyAuthMap[AccessKeyID],
		SecretAccessKey: keyAuthMap[SecretAccessKey],
		AWSServiceName:  "s3",
		Region:          region,
		UserID:          keyAuthMap[ParentUserID],
		KeyAuthEnabled:  keyAuthEnabled,
	}

	optionsList := d.Get(Options).(*schema.Set).List()
	if len(optionsList) > 0 {
		optionsMap := optionsList[0].(map[string]interface{})
		config.RetryCount = optionsMap[RetryCount].(int)
	} else {
		config.RetryCount = 10
	}

	login := client.Login{
		Username:     username,
		Password:     password,
		ParentUserID: keyAuthMap[ParentUserID],
	}

	csClient := client.NewClient(config, &login)
	if !keyAuthEnabled {
		authResponseString, err := csClient.Auth(ctx)
		if err != nil {
			return nil, utils.ProviderConfigurationError(err)

		}

		if err := json.Unmarshal([]byte(authResponseString), &authResp); err != nil {
			return nil, utils.ProviderConfigurationError(utils.UnmarshalJsonError(err))
		}

		if authResp.Token == nil {
			return nil, utils.ProviderConfigurationError(fmt.Errorf(
				"Login auth failed => \n"+
					"Code: %s \n"+
					"Message: %s",
				*authResp.Code, *authResp.Message,
			))
		}
	}

	providerMeta := &models.ProviderMeta{
		CSClient: csClient,
	}

	if authResp.Token != nil {
		providerMeta.Token = *authResp.Token
	}

	return providerMeta, nil

}

func getConfig(ctx context.Context, data *schema.ResourceData, value string) (string, diag.Diagnostics) {
	cred := data.Get(value).(string)
	if cred == "" {
		return "", utils.ProviderConfigurationError(utils.UndefinedError(value))
	}

	return cred, nil
}

func validateUrl(i interface{}, k string) (warnings []string, errors []error) {
	url := i.(string)
	if !strings.HasPrefix(url, "https://") {
		return nil, []error{fmt.Errorf("Invalid URL, must start with 'https://'")}
	}

	return nil, nil
}

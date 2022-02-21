package main

import (
	"context"
	"cs-tf-provider/client"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

type ProviderMeta struct {
	CSClient *client.CSClient
	token    string
}

type AuthResponse struct {
	Token string
}

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
							Required: true,
							ForceNew: true,
						},
					},
				},
				Required:    true,
				ForceNew:    true,
				Description: "List of fields in logs to include or exclude from parsing. If nothing is specified, all fields will be parsed",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"chaossearch_object_group": resourceObjectGroup(),
			// "chaossearch_indexing_state": resourceIndexingState(),
			"chaossearch_view": resourceView(),
			//"chaossearch_user_group":resourceUserGroup(),
		},
		// DataSourcesMap: map[string]*schema.Resource{
		// 	"chaossearch_object_groups": dataSourceObjectGroups(),
		// },
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	accessKeyID := d.Get("access_key_id").(string)
	secretAccessKey := d.Get("secret_access_key").(string)
	region := d.Get("region").(string)

	var username_ string
	var password_ string
	var parentUserId_ string

	if d.Get("login").(*schema.Set).Len() > 0 {
		columnSelectionInterfaces := d.Get("login").(*schema.Set).List()[0]
		columnSelectionInterface := columnSelectionInterfaces.(map[string]interface{})
		username_ = columnSelectionInterface["user_name"].(string)
		password_ = columnSelectionInterface["password"].(string)
		parentUserId_ = columnSelectionInterface["parent_user_id"].(string)

	}

	login_ := client.Login{
		Username:     username_,
		Password:     password_,
		ParentUserId: parentUserId_,
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

	csClient := client.NewClient(config, &login_)

	authResponseString, err := csClient.Auth(ctx)

	log.Debug("authResponseString-->", authResponseString)

	if err != nil {
		return nil, diag.Errorf("Token generation fail..")

	} else {
		tokenData := AuthResponse{}

		if err := json.Unmarshal([]byte(authResponseString), &tokenData); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err), nil
		}

		providerMeta := &ProviderMeta{
			CSClient: csClient,
			token:    tokenData.Token,
		}
		return providerMeta, nil

	}

}

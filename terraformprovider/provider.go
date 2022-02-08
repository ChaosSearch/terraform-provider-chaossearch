package main

import (
	"context"
	"cs-tf-provider/client"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

type ProviderMeta struct {
	Client *client.Client
	token  string
}

type AuthResponse struct{
	Token string
}
// Provider -
func Provider() *schema.Provider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHAOSSEARCH_URL", ""),
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHAOSSEARCH_ACCESS_KEY_ID", ""),
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHAOSSEARCH_SECRET_ACCESS_KEY", ""),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHAOSSEARCH_REGION", "ap-south-1"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// "chaossearch_object_group":   resourceObjectGroup(),
			// "chaossearch_indexing_state": resourceIndexingState(),
			"chaossearch_view": resourceView(),
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

	csClient := client.NewClient(config)

	// logFile, err := os.Create("terraform-provider-chaossearch.log")
	// if err != nil {
	// 	return nil, diag.FromErr(err)
	// }

	// log.SetOutput(logFile)
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// // Close logFile when context is closed
	// go func() {
	// 	<-context.Background().Done()
	// 	log.Println("END")
	// 	logFile.Sync()
	// 	logFile.Close()
	// }()

	authResponseString, err := Auth()

	log.Debug("authResponseString-->", authResponseString)

	if err != nil {
		return nil, diag.Errorf("Token generation fail..")

	} else {
		tokenData:= AuthResponse{}
		json.Unmarshal([]byte(authResponseString),&tokenData)

		providerMeta := &ProviderMeta{
			Client: csClient,
			token:  tokenData.Token,
		}
		return providerMeta, nil

	}
	//

	//return providerMeta, nil
}

package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/models"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSubAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubAccountCreate,
		ReadContext:   resourceSubAccountRead,
		UpdateContext: resourceSubAccountCreate,
		DeleteContext: resourceSubAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hocon": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hocon_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSubAccountCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	createSubAccountRequest := &client.CreateSubAccountRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		UserInfoBlock: client.UserInfoBlock{
			Username: data.Get("username").(string),
			FullName: data.Get("full_name").(string),
			Email:    data.Get("username").(string),
		},
		GroupIds: data.Get("group_ids").([]interface{}),
		Password: data.Get("password").(string),
		HoCon:    data.Get("hocon").([]interface{}),
	}

	if err := c.CreateSubAccount(ctx, createSubAccountRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("username").(string))

	return nil
}

func resourceSubAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var username string
	client := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(client.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	userResp, err := client.ListUsers(ctx, meta.(*models.ProviderMeta).Token)
	if err != nil {
		return diag.FromErr(err)
	}

	if data.Get("username") != nil {
		username = data.Get("username").(string)
	}

	subaccount, err := SortAndGetSubAccount(userResp.Users[0].SubAccounts, username)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("username", subaccount.Username)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("full_name", subaccount.FullName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("hocon_json", subaccount.Hocon)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("group_ids", subaccount.GroupIds)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSubAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var username string
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	if data.Get("username") != nil {
		username = data.Get("username").(string)
	}

	deleteSubAccountRequest := &client.BasicRequest{
		AuthToken: meta.(*models.ProviderMeta).Token,
		Id:        username,
	}

	if err := c.DeleteSubAccount(ctx, deleteSubAccountRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(username)

	return nil
}

func SortAndGetSubAccount(subAccounts []client.SubAccount, target string) (*client.SubAccount, error) {
	if len(subAccounts) == 1 {
		return &subAccounts[0], nil
	}

	if target == "" {
		return nil, fmt.Errorf("Target Username not supplied")
	}

	sort.Slice(subAccounts[:], func(i, j int) bool {
		return subAccounts[i].Username < subAccounts[j].Username
	})

	low := 0
	high := len(subAccounts)
	for low <= high {
		median := (low + high) / 2
		if subAccounts[median].Username < target {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	if low == len(subAccounts) || subAccounts[low].Username != target {
		return nil, fmt.Errorf("SubAccount not found")
	}

	return &subAccounts[low], nil
}

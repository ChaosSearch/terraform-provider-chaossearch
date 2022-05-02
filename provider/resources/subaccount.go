package resources

import (
	"context"
	"cs-tf-provider/client"
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
			"user_info_block": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hocon": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSubAccountCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	userInfoList := data.Get("user_info_block").(*schema.Set).List()
	if len(userInfoList) > 0 {
		userInfoMap := userInfoList[0].(map[string]interface{})
		userInfoBlock := client.UserInfoBlock{
			Username: userInfoMap["username"].(string),
			FullName: userInfoMap["full_name"].(string),
			Email:    userInfoMap["email"].(string),
		}

		createSubAccountRequest := &client.CreateSubAccountRequest{
			AuthToken:     meta.(*models.ProviderMeta).Token,
			UserInfoBlock: userInfoBlock,
			GroupIds:      data.Get("group_ids").([]interface{}),
			Password:      data.Get("password").(string),
			HoCon:         data.Get("hocon").([]interface{}),
		}

		if err := c.CreateSubAccount(ctx, createSubAccountRequest); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(userInfoMap["username"].(string))
	}

	return nil
}

func resourceSubAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*models.ProviderMeta).CSClient
	userInfoList := data.Get("user_info_block").(*schema.Set).List()
	if len(userInfoList) > 0 {
		userInfoMap := userInfoList[0].(map[string]interface{})
		userResp, err := client.ListUsers(ctx, meta.(*models.ProviderMeta).Token)
		if err != nil {
			return diag.FromErr(err)
		}

		subaccount, err := sortAndGetSubAccount(userResp.Users[0].SubAccounts, userInfoMap["email"].(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = data.Set("user_info_block", []interface{}{
			map[string]interface{}{
				"username":  subaccount.Username,
				"full_name": subaccount.FullName,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}

		err = data.Set("hocon", []string{
			subaccount.Hocon,
		})

		if err != nil {
			return diag.FromErr(err)
		}

		err = data.Set("group_ids", subaccount.GroupIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceSubAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	userInfoList := data.Get("user_info_block").(*schema.Set).List()
	if len(userInfoList) > 0 {
		userInfoMap := userInfoList[0].(map[string]interface{})
		deleteSubAccountRequest := &client.DeleteSubAccountRequest{
			AuthToken: meta.(*models.ProviderMeta).Token,
			Username:  userInfoMap["username"].(string),
		}

		if err := c.DeleteSubAccount(ctx, deleteSubAccountRequest); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(userInfoMap["username"].(string))
	}

	return nil
}

func sortAndGetSubAccount(subAccounts []client.SubAccount, target string) (*client.SubAccount, error) {
	if len(subAccounts) == 1 {
		return &subAccounts[0], nil
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

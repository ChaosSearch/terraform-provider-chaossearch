package resources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/client/utils"
	"cs-tf-provider/provider/models"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenValue := meta.(*models.ProviderMeta).Token
	userGroup, err := GroupObject(data, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.CreateUserGroup(ctx, userGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.ID)
	return resourceGroupRead(ctx, data, meta)
}

func GroupObject(data *schema.ResourceData, authToken string) (*client.CreateUserGroupRequest, error) {
	var policy []client.Permission

	permissions, err := structure.NormalizeJsonString(data.Get("permissions").(string))
	if err != nil {
		return nil, utils.NormalizingJsonError(err)
	}

	err = json.Unmarshal([]byte(permissions), &policy)
	if err != nil {
		return nil, utils.UnmarshalJsonError(err)
	}

	return &client.CreateUserGroupRequest{
		AuthToken:   authToken,
		ID:          data.Id(),
		Name:        data.Get("name").(string),
		Permissions: policy,
	}, nil
}

func resourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenValue := meta.(*models.ProviderMeta).Token
	req := &client.ReadUserGroupRequest{
		AuthToken: tokenValue,
		ID:        data.Id(),
	}

	resp, err := c.ReadUserGroup(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		return diag.Errorf("Couldn't find User Group: %s", err)
	}

	userGroupContent, err := CreateUserGroupResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("name", userGroupContent[0]["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("permissions", userGroupContent[0]["permissions"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func CreateUserGroupResponse(resp *client.UserGroup) ([]map[string]interface{}, error) {
	permission_json, err := json.Marshal(resp.Permissions)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	permissions, err := structure.NormalizeJsonString(string(permission_json))
	if err != nil {
		return nil, utils.NormalizingJsonError(err)
	}

	return []map[string]interface{}{
		{
			"id":          resp.ID,
			"name":        resp.Name,
			"permissions": permissions,
		},
	}, nil
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenValue := meta.(*models.ProviderMeta).Token
	userGroup, err := GroupObject(data, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.UpdateUserGroup(ctx, userGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.ID)
	return nil
}

func resourceGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*models.ProviderMeta).CSClient
	err := utils.ValidateAuthType(c.Config.KeyAuthEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenValue := meta.(*models.ProviderMeta).Token
	deleteUserGroupRequest := &client.DeleteUserGroupRequest{
		AuthToken: tokenValue,
		ID:        data.Id(),
	}

	if err := c.DeleteUserGroup(ctx, deleteUserGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

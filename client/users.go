package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *CSClient) ListUsers(ctx context.Context, authToken string) (*ListUsersResponse, error) {
	var resp ListUsersResponse
	url := fmt.Sprintf("%s/user/manifest", c.config.URL)
	request := ClientRequest{
		RequestType: POST,
		Url:         url,
		AuthToken:   authToken,
	}

	httpResp, err := c.createAndSendReq(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("List Users Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

func (c *CSClient) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*UserGroup, error) {
	var readUserGroupResp []UserGroup
	url := fmt.Sprintf("%s/user/groups", c.config.URL)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}

	request := ClientRequest{
		RequestType: POST,
		Url:         url,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	}

	httpResp, err := c.createAndSendReq(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Create User Group Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &readUserGroupResp[0], nil
}

func (c *CSClient) ReadUserGroup(ctx context.Context, req *ReadUserGroupRequest) (*UserGroup, error) {
	var readUserGroupResp UserGroup
	url := fmt.Sprintf("%s/user/group/%s", c.config.URL, req.ID)
	request := ClientRequest{
		RequestType: GET,
		Url:         url,
		AuthToken:   req.AuthToken,
	}

	httpResp, err := c.createAndSendReq(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Read User Group Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &readUserGroupResp, nil
}

func (c *CSClient) UpdateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*UserGroup, error) {
	var readUserGroupResp []UserGroup
	url := fmt.Sprintf("%s/user/groups", c.config.URL)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}

	request := ClientRequest{
		RequestType: POST,
		Url:         url,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	}

	httpResp, err := c.createAndSendReq(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Update User Group Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &readUserGroupResp[0], nil
}

func (c *CSClient) DeleteUserGroup(ctx context.Context, req *DeleteUserGroupRequest) error {
	deleteUserGroupID := url.PathEscape(req.ID)
	url := fmt.Sprintf("%s/user/group/%s", c.config.URL, deleteUserGroupID)
	request := ClientRequest{
		RequestType: DELETE,
		Url:         url,
		AuthToken:   req.AuthToken,
	}

	httpResp, err := c.createAndSendReq(ctx, request)
	if err != nil {
		return fmt.Errorf("Delete User Group Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func marshalCreateUserGroupRequest(req *CreateUserGroupRequest) ([]byte, error) {
	body := []interface{}{
		map[string]interface{}{
			"id":          req.ID,
			"name":        req.Name,
			"permissions": req.Permission,
		},
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

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
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/manifest", c.config.URL),
		RequestType: POST,
		AuthToken:   authToken,
	})

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
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/groups", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

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
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/group/%s", c.config.URL, req.ID),
		RequestType: GET,
		AuthToken:   req.AuthToken,
	})

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
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/groups", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

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
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/group/%s", c.config.URL, deleteUserGroupID),
		RequestType: DELETE,
		AuthToken:   req.AuthToken,
	})

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

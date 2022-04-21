package client

import (
	"bytes"
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *CSClient) ListUsers(ctx context.Context, authToken string) (*ListUsersResponse, error) {
	url := fmt.Sprintf("%s/user/manifest", c.config.URL)
	httpReq, err := http.NewRequestWithContext(ctx, POST, url, nil)
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(authToken, httpReq, nil)
	if err != nil {
		return nil, utils.SubmitRequestError(POST, url, err)
	}
	defer httpResp.Body.Close()

	var resp ListUsersResponse
	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *CSClient) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*Group, error) {
	url := fmt.Sprintf("%s/user/groups", c.config.URL)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, utils.SubmitRequestError(POST, url, err)
	}
	defer httpResp.Body.Close()

	var readUserGroupResp []Group
	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, err
	}
	return &readUserGroupResp[0], err
}

func (c *CSClient) ReadUserGroup(ctx context.Context, req *ReadUserGroupRequest) (*Group, error) {
	var resp Group
	if err := c.ReadUserGroupByID(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *CSClient) ReadUserGroupByID(ctx context.Context, req *ReadUserGroupRequest, resp *Group) error {
	url := fmt.Sprintf("%s/user/group/%s", c.config.URL, req.ID)
	httpReq, err := http.NewRequestWithContext(ctx, GET, url, nil)
	if err != nil {
		return utils.CreateRequestError(err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return utils.SubmitRequestError(GET, url, err)
	}
	defer httpResp.Body.Close()

	var readUserGroupResp Group
	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return err
	}

	resp.ID = readUserGroupResp.ID
	resp.Name = readUserGroupResp.Name
	resp.Permissions = readUserGroupResp.Permissions
	return nil
}

func (c *CSClient) UpdateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*Group, error) {
	url := fmt.Sprintf("%s/user/groups", c.config.URL)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, PUT, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, utils.SubmitRequestError(POST, url, err)
	}
	defer httpResp.Body.Close()

	var readUserGroupResp []Group
	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, err
	}
	return &readUserGroupResp[0], nil
}

func (c *CSClient) DeleteUserGroup(ctx context.Context, req *DeleteUserGroupRequest) error {
	deleteUserGroupID := url.PathEscape(req.ID)
	deleteUserGroupURL := fmt.Sprintf("%s/user/group/%s", c.config.URL, deleteUserGroupID)
	httpReq, err := http.NewRequestWithContext(ctx, DELETE, deleteUserGroupURL, nil)
	if err != nil {
		return utils.CreateRequestError(err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return utils.SubmitRequestError(DELETE, deleteUserGroupURL, err)
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

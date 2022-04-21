package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *CSClient) ListUsers(ctx context.Context, authToken string) (*ListUsersResponse, error) {
	url := fmt.Sprintf("%s/user/manifest", c.config.URL)

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.signV2AndDo(authToken, httpReq, nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var resp ListUsersResponse
	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response body: %s", err)
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
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", POST, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var readUserGroupResp []Group
	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response body sdjhskdhskdskdskdksdkskjd: %s", err)
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
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", GET, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var readUserGroupResp Group
	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response body : %s", err)
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
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", POST, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var readUserGroupResp []Group
	if err := c.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response body sdjhskdhskdskdskdksdkskjd: %s", err)
	}
	return &readUserGroupResp[0], err
}

func (c *CSClient) DeleteUserGroup(ctx context.Context, req *DeleteUserGroupRequest) error {

	deleteUserGroupID := url.PathEscape(req.ID)
	deleteUserGroupURL := fmt.Sprintf("%s/user/group/%s", c.config.URL, deleteUserGroupID)

	httpReq, err := http.NewRequestWithContext(ctx, DELETE, deleteUserGroupURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, deleteUserGroupURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body %s", err)
		}
	}(httpResp.Body)

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
		return nil, err
	}
	return bodyAsBytes, nil
}

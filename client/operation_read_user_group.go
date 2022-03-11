package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

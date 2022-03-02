package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (csClient *CSClient) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*Group, error) {
	method := "POST"
	url := fmt.Sprintf("%s/user/groups", csClient.config.URL)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	log.Debug(" adding headers...")
	log.Warn("httpReq-->", httpReq)
	httpResp, err := csClient.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var readUserGroupResp []Group
	if err := csClient.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response body sdjhskdhskdskdskdksdkskjd: %s", err)
	}
	return &readUserGroupResp[0], err
}

func marshalCreateUserGroupRequest(req *CreateUserGroupRequest) ([]byte, error) {
	body := []interface{}{
		map[string]interface{}{
			"id":          req.Id,
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

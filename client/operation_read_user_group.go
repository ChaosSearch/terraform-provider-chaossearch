package client

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (csClient *CSClient) ReadUserGroup(ctx context.Context, req *ReadUserGroupRequest) (*Group, error) {
	var resp Group
	if err := csClient.ReadUserGroupById(ctx, req, &resp); err != nil {
		return nil, err
	}
	log.Printf("ReadObjectGroupResponse: %+v", resp)
	return &resp, nil
}

func (csClient *CSClient) ReadUserGroupById(ctx context.Context, req *ReadUserGroupRequest, resp *Group) error {
	method := "GET"
	url := fmt.Sprintf("%s/user/group/%s", csClient.config.URL, req.ID)
	log.Debug("ReadUserGroupById--->")
	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	var sessionToken = req.AuthToken
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var readUserGroupResp Group
	if err := csClient.unmarshalJSONBody(httpResp.Body, &readUserGroupResp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response body : %s", err)
	}
	resp.Id = readUserGroupResp.Id
	resp.Name = readUserGroupResp.Name
	resp.Permissions = readUserGroupResp.Permissions
	return nil
}

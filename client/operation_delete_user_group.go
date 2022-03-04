package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (csClient *CSClient) DeleteUserGroup(ctx context.Context, req *DeleteUserGroupRequest) error {
	method := "DELETE"
	deleteUserGroupId := url.PathEscape(req.ID)
	deleteUserGroupUrl := fmt.Sprintf("%s/user/group/%s", csClient.config.URL, deleteUserGroupId)

	httpReq, err := http.NewRequestWithContext(ctx, method, deleteUserGroupUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	var sessionToken = req.AuthToken
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, deleteUserGroupUrl, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body %s", err)
		}
	}(httpResp.Body)

	return nil
}

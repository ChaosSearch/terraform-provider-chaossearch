package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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

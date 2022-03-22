package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *CSClient) DeleteObjectGroup(ctx context.Context, req *DeleteObjectGroupRequest) error {

	safeObjectGroupName := url.PathEscape(req.Name)
	deleteURL := fmt.Sprintf("%s/V1/%s", c.config.URL, safeObjectGroupName)

	httpReq, err := http.NewRequestWithContext(ctx, DELETE, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)

	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, deleteURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	return nil
}

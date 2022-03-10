package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *CSClient) DeleteView(ctx context.Context, req *DeleteViewRequest) error {

	safeViewName := url.PathEscape(req.Name)
	deleteViewURL := fmt.Sprintf("%s/V1/%s", c.config.URL, safeViewName)

	httpReq, err := http.NewRequestWithContext(ctx, DELETE, deleteViewURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, deleteViewURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body %s", err)
		}
	}(httpResp.Body)

	return nil
}

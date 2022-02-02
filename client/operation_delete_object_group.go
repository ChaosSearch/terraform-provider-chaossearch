package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (client *Client) DeleteObjectGroup(ctx context.Context, req *DeleteObjectGroupRequest) error {
	method := "DELETE"
	safeObjectGroupName := url.PathEscape(req.Name)
	url := fmt.Sprintf("%s/V1/%s", client.config.URL, safeObjectGroupName)

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("Failed to create request: %s", err)
	}

	httpResp, err := client.signAndDo(httpReq, nil)
	if err != nil {
		return fmt.Errorf("Failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()

	return nil
}

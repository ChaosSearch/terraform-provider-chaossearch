package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (client *Client) DeleteView(ctx context.Context, req *DeleteViewRequest) error {
	method := "DELETE"
	safeViewName := url.PathEscape(req.Name)
	url := fmt.Sprintf("%s/V1/%s", client.config.URL, safeViewName)

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := client.signAndDo(httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()

	return nil
}

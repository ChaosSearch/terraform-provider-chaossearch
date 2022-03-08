package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (csClient *CSClient) DeleteObjectGroup(ctx context.Context, req *DeleteObjectGroupRequest) error {
	method := "DELETE"
	safeObjectGroupName := url.PathEscape(req.Name)
	deleteUrl := fmt.Sprintf("%s/V1/%s", csClient.config.URL, safeObjectGroupName)

	httpReq, err := http.NewRequestWithContext(ctx, method, deleteUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, nil)

	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, deleteUrl, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	return nil
}

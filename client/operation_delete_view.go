package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (csClient *CSClient) DeleteView(ctx context.Context, req *DeleteViewRequest) error {
	method := "DELETE"
	safeViewName := url.PathEscape(req.Name)
	deleteViewUrl := fmt.Sprintf("%s/V1/%s", csClient.config.URL, safeViewName)

	httpReq, err := http.NewRequestWithContext(ctx, method, deleteViewUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	var sessionToken = req.AuthToken
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, deleteViewUrl, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body %s", err)
		}
	}(httpResp.Body)

	return nil
}

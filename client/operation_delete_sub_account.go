package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *CSClient) DeleteSubAccount(ctx context.Context, req *DeleteSubAccountRequest) error {

	deleteURL := fmt.Sprintf("%s/user/deleteSubAccount", c.config.URL)

	bodyAsBytes, _ := marshalDeleteSubAccountRequest(req)
	httpReq, err := http.NewRequestWithContext(ctx, POST, deleteURL, bytes.NewReader(bodyAsBytes))

	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
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

func marshalDeleteSubAccountRequest(req *DeleteSubAccountRequest) ([]byte, error) {
	body := map[string]interface{}{
		"Username": req.Username,
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

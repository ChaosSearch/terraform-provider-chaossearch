package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (csClient *CSClient) DeleteSubAccount(ctx context.Context, req *DeleteSubAccountRequest) error {
	method := "POST"
	deleteUrl := fmt.Sprintf("%s/user/deleteSubAccount", csClient.config.URL)

	bodyAsBytes, err := marshalDeleteSubAccountRequest(req)
	httpReq, err := http.NewRequestWithContext(ctx, method, deleteUrl, bytes.NewReader(bodyAsBytes))

	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := csClient.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
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

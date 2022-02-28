package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

func (csClient *CSClient) CreateSubAccount(ctx context.Context, req *CreateSubAccountRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/user/createSubAccount", csClient.config.URL)

	bodyAsBytes, err := marshalCreateSubAccountRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := csClient.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	return nil
}

func marshalCreateSubAccountRequest(req *CreateSubAccountRequest) ([]byte, error) {
	body := map[string]interface{}{
		"UserInfoBlock": map[string]interface{}{
			"Username": req.UserInfoBlock.Username,
			"FullName": req.UserInfoBlock.FullName,
			"Email":    req.UserInfoBlock.Email,
		},
		"GroupIds": req.GroupIds,
		"Password": req.Password,
		"Hocon":    req.HoCon,
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (csClient *CSClient) UpdateObjectGroup(ctx context.Context, req *UpdateObjectGroupRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/updateObjectGroup", csClient.config.URL)

	bodyAsBytes, err := marshalUpdateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	var sessionToken = req.AuthToken
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, bodyAsBytes)
	//httpResp, err := client.signV4AndDo(httpReq, bodyAsBytes)
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

func marshalUpdateObjectGroupRequest(req *UpdateObjectGroupRequest) ([]byte, error) {
	body := map[string]interface{}{
		"bucket":         req.Name,
		"indexRetention": req.IndexRetention,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}

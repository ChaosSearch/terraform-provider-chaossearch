package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (csClient *CSClient) UpdateIndexingState(ctx context.Context, req *UpdateIndexingStateRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/model", csClient.config.URL)

	bodyAsBytes, err := marshalUpdateIndexingStateRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := csClient.signV4AndDo(httpReq, bodyAsBytes)
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

func marshalUpdateIndexingStateRequest(req *UpdateIndexingStateRequest) ([]byte, error) {
	var modelMode int

	if req.Active {
		modelMode = 0
	} else {
		modelMode = -1
	}

	body := map[string]interface{}{
		"BucketName": req.ObjectGroupName,
		"ModelMode":  modelMode,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}

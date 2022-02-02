package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// For documentation see: https://docs.chaossearch.io/reference#bucketmodel
func (client *Client) UpdateIndexingState(ctx context.Context, req *UpdateIndexingStateRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/model", client.config.URL)

	bodyAsBytes, err := marshalUpdateIndexingStateRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("Failed to create request: %s", err)
	}

	httpResp, err := client.signAndDo(httpReq, bodyAsBytes)
	if err != nil {
		return fmt.Errorf("Failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()

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

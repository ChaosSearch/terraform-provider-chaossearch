package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (csClient *CSClient) CreateIndexModel(ctx context.Context, req *IndexModelRequest) (*IndexModelResponse, error) {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/model", csClient.config.URL)
	bodyAsBytes, err := marshalIndexModelRequest(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	httpResp, err := csClient.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	var indexModelResponse IndexModelResponse
	if err := csClient.unmarshalJSONBody(httpResp.Body, &indexModelResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response body: %s", err)
	}
	return &indexModelResponse, err
}

func marshalIndexModelRequest(req *IndexModelRequest) ([]byte, error) {
	body := map[string]interface{}{
		"BucketName": req.BucketName,
		"ModelMode":  req.ModelMode,
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

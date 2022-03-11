package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *CSClient) CreateIndexModel(x context.Context, e *IndexModelRequest) (*IndexModelResponse, error) {
	//method := "POST"
	url := fmt.Sprintf("%s/Bucket/model", c.config.URL)
	bodyAsBytes, err := marshalIndexModelRequest(e)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(x, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	httpResp, err := c.signV2AndDo(e.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", POST, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	var indexModelResponse IndexModelResponse
	if err := c.unmarshalJSONBody(httpResp.Body, &indexModelResponse); err != nil {
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

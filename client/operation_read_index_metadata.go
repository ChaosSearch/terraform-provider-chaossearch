package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (c *CSClient) ReadIndexMetadata(ctx context.Context, req *IndexMetadataRequest) (*IndexMetadataResponse, error) {
	url := fmt.Sprintf("%s/Bucket/metadata", c.config.URL)
	bodyAsBytes, err := marshalIndexMetadataRequest(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", POST, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	result := c.processResponse(req.BucketName, httpResp)
	return &result, err
}

func (c *CSClient) processResponse(requestBucketName string, httpResp *http.Response) IndexMetadataResponse {
	respBodyAsBytes, _ := ioutil.ReadAll(httpResp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBodyAsBytes, &result); err != nil {
		_ = fmt.Errorf("failed to unmarshal JSON: %s", err)
	}
	metadata := result["Metadata"].(map[string]interface{})
	bucket := metadata[requestBucketName].(map[string]interface{})

	bucketName := bucket["Bucket"].(string)
	lastIndexTime := bucket["LastIndexTime"].(float64)
	state := bucket["State"].(string)

	response := IndexMetadataResponse{bucketName, lastIndexTime, state}
	return response
}

func marshalIndexMetadataRequest(req *IndexMetadataRequest) ([]byte, error) {
	body := map[string]interface{}{
		"BucketNames": []interface{}{req.BucketName},
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// There are other nested structures that we don't marshal here because we don't need to.
//  Feel free to add them in the future when needed
type readBucketMetadataResponse struct {
	Bucket string `json:"Bucket"`
	State  string `json:"State"`
}

func (csClient *CSClient) ReadIndexingState(ctx context.Context, req *ReadIndexingStateRequest) (*IndexingState, error) {
	method := "POST"
	body := &readBucketMetadataRequest{
		BucketName: req.ObjectGroupName,
		Stats:      false,
	}

	response, err := makeGetBucketMetadataRequest(method, csClient, ctx, body)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}

	bucketMetadata := &IndexingState{
		ObjectGroupName: response.Bucket,
		Active:          false,
	}

	if response.State == "Active" || response.State == "Idle" {
		bucketMetadata.Active = true
	}

	return bucketMetadata, nil
}

func makeGetBucketMetadataRequest(method string, c *CSClient, ctx context.Context, body *readBucketMetadataRequest) (*readBucketMetadataResponse, error) {
	url := fmt.Sprintf("%s/Bucket/metadata", c.config.URL)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request body: %s", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV4AndDo(httpReq, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var response readBucketMetadataResponse
	err = c.unmarshalJSONBody(httpResp.Body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %s", err)
	}

	return &response, nil
}

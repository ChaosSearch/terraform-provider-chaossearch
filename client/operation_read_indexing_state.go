package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// There are other nested structures that we don't marshal here because we don't need to.
//  Feel free to add them in the future when needed
type readBucketMetadataResponse struct {
	Bucket string `json:"Bucket"`
	State  string `json:"State"`
}

// For documentation see: https://docs.chaossearch.io/reference#bucketmodel
func (client *Client) ReadIndexingState(ctx context.Context, req *ReadIndexingStateRequest) (*IndexingState, error) {
	method := "POST"
	body := &readBucketMetadataRequest{
		BucketName: req.ObjectGroupName,
		Stats:      false,
	}

	response, err := makeGetBucketMetadataRequest(method, client, ctx, body)
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

func makeGetBucketMetadataRequest(method string, client *Client, ctx context.Context, body *readBucketMetadataRequest) (*readBucketMetadataResponse, error) {
	url := fmt.Sprintf("%s/Bucket/metadata", client.config.URL)

	jsonedBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request body: %s", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(jsonedBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := client.signAndDo(httpReq, jsonedBody)
	if err != nil {
		return nil, fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()

	var response readBucketMetadataResponse
	err = client.unmarshalJSONBody(httpResp.Body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %s", err)
	}

	return &response, nil
}

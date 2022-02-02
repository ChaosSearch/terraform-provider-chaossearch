package client

import (
	"context"
	"fmt"
	"net/http"
)

func (client *Client) ListBuckets(ctx context.Context) (*ListBucketsResponse, error) {
	url := fmt.Sprintf("%s/V1/", client.config.URL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := client.signAndDo(httpReq, nil)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	var resp ListBucketsResponse
	if err := client.unmarshalXMLBody(httpResp.Body, &resp); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal XML response body: %s", err)
	}

	return &resp, nil
}

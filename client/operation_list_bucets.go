package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *CSClient) ListBuckets(ctx context.Context, authToken string) (*ListBucketsResponse, error) {
	url := fmt.Sprintf("%s/V1/", c.config.URL)

	httpReq, err := http.NewRequestWithContext(ctx, GET, url, nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.signV2AndDo(authToken, httpReq, nil)

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var resp ListBucketsResponse
	if err := c.unmarshalXMLBody(httpResp.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML response body: %s", err)
	}

	return &resp, nil
}

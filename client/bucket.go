package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"fmt"
)

func (c *CSClient) ListBuckets(ctx context.Context, authToken string) (*ListBucketsResponse, error) {
	var resp ListBucketsResponse
	url := fmt.Sprintf("%s/V1/", c.config.URL)
	request := ClientRequest{
		RequestType: GET,
		Url:         url,
		AuthToken:   authToken,
	}

	httpResp, err := c.createAndSendReq(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("List Bucket Failure => %s", err)
	}

	if err := c.unmarshalXMLBody(httpResp.Body, &resp); err != nil {
		return nil, utils.UnmarshalXmlError(err)
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

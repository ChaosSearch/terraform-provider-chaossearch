package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"fmt"
)

func (c *CSClient) ListBuckets(ctx context.Context, authToken string) (*ListBucketsResponse, error) {
	var resp ListBucketsResponse
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/V1/", c.config.URL),
		RequestType: GET,
		AuthToken:   authToken,
	})

	if err != nil {
		return nil, fmt.Errorf("List Bucket Failure => %s", err)
	}

	if err := c.unmarshalXMLBody(httpResp.Body, &resp); err != nil {
		return nil, utils.UnmarshalXmlError(err)
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

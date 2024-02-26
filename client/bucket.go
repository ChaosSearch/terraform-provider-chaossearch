package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"fmt"
)

func (c *CSClient) ListBuckets(ctx context.Context, authToken string) (*ListBucketsResponse, error) {
	var resp ListBucketsResponse
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/V1/", c.Config.URL),
		RequestType: GET,
		AuthToken:   authToken,
		Headers: map[string]string{
			"x-amz-chaossumo-bucket-tagging":   "true",
			"x-amz-chaossumo-bucket-transform": "exclude-indexes",
		},
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

func (c *CSClient) ReadBucketDataset(
	ctx context.Context,
	authToken,
	bucketName string,
) (*ReadBucketDatasetResp, error) {
	var datasetResp ReadBucketDatasetResp

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/dataset/name/%s?state=true", c.Config.URL, bucketName),
		RequestType: GET,
		AuthToken:   authToken,
	})

	if err != nil {
		return nil, fmt.Errorf("Read Bucket Metadata Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &datasetResp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &datasetResp, nil
}

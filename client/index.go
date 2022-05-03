package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
)

func (c *CSClient) CreateIndexModel(ctx context.Context, req *IndexModelRequest) (*IndexModelResponse, error) {
	var indexModelResponse IndexModelResponse
	bodyAsBytes, err := marshalIndexModelRequest(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/model", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return nil, fmt.Errorf("Create Index Model Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &indexModelResponse); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &indexModelResponse, nil
}

func (c *CSClient) DeleteIndexModel(ctx context.Context, indexName, authToken string) error {
	if indexName != "" {
		httpResp, err := c.createAndSendReq(ctx, ClientRequest{
			Url:         fmt.Sprintf("%s/V1/%s", c.config.URL, indexName),
			RequestType: DELETE,
			AuthToken:   authToken,
		})

		if err != nil {
			return fmt.Errorf("Delete Index Model Failure => %s", err)
		}

		defer httpResp.Body.Close()
	}

	return nil
}

func (c *CSClient) ReadIndexModel(ctx context.Context, bucketName, authToken string) (*ListBucketResponse, error) {
	var listBucketResponse ListBucketResponse
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf(`%s/V1/%s?list-type=2&delimiter=/&max-keys=100`, c.config.URL, bucketName),
		RequestType: GET,
		AuthToken:   authToken,
		Headers: map[string]string{
			"x-amz-chaossumo-bucket-transform": "indexed",
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Read Index Model Failure => %s", err)
	}

	if err := c.unmarshalXMLBody(httpResp.Body, &listBucketResponse); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &listBucketResponse, nil
}

func marshalIndexModelRequest(req *IndexModelRequest) ([]byte, error) {
	body := map[string]interface{}{
		"BucketName": req.BucketName,
		"ModelMode":  req.ModelMode,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

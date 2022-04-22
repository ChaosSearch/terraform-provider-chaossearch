package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
)

func (c *CSClient) CreateIndexModel(ctx context.Context, req *IndexModelRequest) (*IndexModelResponse, error) {
	var indexModelResponse IndexModelResponse
	url := fmt.Sprintf("%s/Bucket/model", c.config.URL)
	bodyAsBytes, err := marshalIndexModelRequest(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.createAndSendReq(ctx, req.AuthToken, url, POST, bodyAsBytes)
	if err != nil {
		return nil, fmt.Errorf("Create Index Model Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &indexModelResponse); err != nil {
		return nil, err
	}

	return &indexModelResponse, nil
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

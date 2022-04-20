package client

import (
	"bytes"
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *CSClient) CreateIndexModel(x context.Context, e *IndexModelRequest) (*IndexModelResponse, error) {
	url := fmt.Sprintf("%s/Bucket/model", c.config.URL)
	bodyAsBytes, err := marshalIndexModelRequest(e)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(x, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(e.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return nil, utils.SubmitRequestError(POST, url, err)
	}
	defer httpResp.Body.Close()

	var indexModelResponse IndexModelResponse
	if err := c.unmarshalJSONBody(httpResp.Body, &indexModelResponse); err != nil {
		return nil, err
	}

	return &indexModelResponse, err
}

func (c *CSClient) ReadIndexMetadata(ctx context.Context, req *IndexMetadataRequest) (*IndexMetadataResponse, error) {
	url := fmt.Sprintf("%s/Bucket/metadata", c.config.URL)
	bodyAsBytes, err := marshalIndexMetadataRequest(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

	if err != nil {
		return nil, utils.SubmitRequestError(POST, url, err)
	}
	defer httpResp.Body.Close()

	result := c.processResponse(req.BucketName, httpResp)
	return &result, err
}

func (c *CSClient) processResponse(requestBucketName string, httpResp *http.Response) IndexMetadataResponse {
	respBodyAsBytes, _ := ioutil.ReadAll(httpResp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBodyAsBytes, &result); err != nil {
		_ = utils.UnmarshalJsonError(err)
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
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
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

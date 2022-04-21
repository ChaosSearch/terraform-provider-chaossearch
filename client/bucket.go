package client

import (
	"bytes"
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *CSClient) ImportBucket(ctx context.Context, req *ImportBucketRequest) error {
	url := fmt.Sprintf("%s/Bucket/importBucket", c.config.URL)
	bodyAsBytes, err := marshalImportBuketRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)
	if err != nil {
		return utils.SubmitRequestError(POST, url, err)
	}
	defer httpResp.Body.Close()

	return nil
}

func (c *CSClient) ListBuckets(ctx context.Context, authToken string) (*ListBucketsResponse, error) {
	url := fmt.Sprintf("%s/V1/", c.config.URL)
	httpReq, err := http.NewRequestWithContext(ctx, GET, url, nil)
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(authToken, httpReq, nil)
	if err != nil {
		return nil, utils.SubmitRequestError(GET, url, err)
	}
	defer httpResp.Body.Close()

	var resp ListBucketsResponse
	if err := c.unmarshalXMLBody(httpResp.Body, &resp); err != nil {
		return nil, utils.UnmarshalXmlError(err)
	}

	return &resp, nil
}

func marshalImportBuketRequest(req *ImportBucketRequest) ([]byte, error) {
	body := map[string]interface{}{
		"bucket":     req.Bucket,
		"hideBucket": req.HideBucket,
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}
	return bodyAsBytes, nil
}

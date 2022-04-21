package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *CSClient) ImportBucket(ctx context.Context, req *ImportBucketRequest) error {
	url := fmt.Sprintf("%s/Bucket/importBucket", c.config.URL)
	log.Debug("Url-->", url)

	bodyAsBytes, err := marshalImportBuketRequest(req)
	if err != nil {
		return err
	}
	log.Debug("method-->", POST)
	log.Debug("bodyAsBytes-->", bodyAsBytes)

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	log.Warn("httpReq-->", httpReq)
	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

	log.Warn("httpResp-->", httpResp)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	return nil
}

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

func marshalImportBuketRequest(req *ImportBucketRequest) ([]byte, error) {
	body := map[string]interface{}{
		"bucket":     req.Bucket,
		"hideBucket": req.HideBucket,
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func (c *CSClient) CreateView(ctx context.Context, req *CreateViewRequest) error {

	url := fmt.Sprintf("%s/Bucket/createView", c.config.URL)
	log.Debug("Url-->", url)
	log.Debug("req-->", req)
	bodyAsBytes, err := marshalCreateViewRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	log.Debug(" adding headers...")
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

func (c *CSClient) ReadView(ctx context.Context, req *ReadViewRequest) (*ReadViewResponse, error) {
	var resp ReadViewResponse
	if err := c.readViewAttr(ctx, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CSClient) readViewAttr(x context.Context, e *ReadViewRequest, r *ReadViewResponse) error {
	url := fmt.Sprintf("%s/Bucket/dataset/name/%s", c.config.URL, e.ID)
	httpReq, err := http.NewRequestWithContext(x, GET, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV2AndDo(e.AuthToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", GET, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var read ReadViewResponse

	if err := c.unmarshalJSONBody(httpResp.Body, &read); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response body: %s", err)
	}

	r.FilterPredicate = read.FilterPredicate
	r.Type = read.Type
	r.MetaData = read.MetaData
	r.RegionAvailability = read.RegionAvailability
	r.ID = read.ID
	r.Bucket = read.Bucket
	r.Pattern = read.Pattern
	r.Transforms = read.Transforms
	r.TimeFieldName = read.TimeFieldName
	r.Sources = read.Sources
	r.Cacheable = read.Cacheable
	r.CaseInsensitive = read.CaseInsensitive
	r.IndexPattern = read.IndexPattern
	return nil
}

func (c *CSClient) DeleteView(ctx context.Context, req *DeleteViewRequest) error {

	safeViewName := url.PathEscape(req.Name)
	deleteViewURL := fmt.Sprintf("%s/V1/%s", c.config.URL, safeViewName)

	httpReq, err := http.NewRequestWithContext(ctx, DELETE, deleteViewURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, deleteViewURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body %s", err)
		}
	}(httpResp.Body)

	return nil
}

func marshalCreateViewRequest(req *CreateViewRequest) ([]byte, error) {
	log.Debug("req.Sources-->", req.Sources)
	body := map[string]interface{}{
		"bucket":          req.Bucket,
		"sources":         req.Sources,
		"indexPattern":    req.IndexPattern,
		"overwrite":       req.Overwrite,
		"caseInsensitive": req.CaseInsensitive,
		"indexRetention":  req.IndexRetention,
		"timeFieldName":   req.TimeFieldName,
		"transforms":      req.Transforms,
		"filter":          req.FilterPredicate,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

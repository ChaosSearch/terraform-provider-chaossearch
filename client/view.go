package client

import (
	"bytes"
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *CSClient) CreateView(ctx context.Context, req *CreateViewRequest) error {
	url := fmt.Sprintf("%s/Bucket/createView", c.config.URL)
	bodyAsBytes, err := marshalCreateViewRequest(req)
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
		return utils.CreateRequestError(err)
	}

	httpResp, err := c.signV2AndDo(e.AuthToken, httpReq, nil)
	if err != nil {
		return utils.SubmitRequestError(GET, url, err)
	}
	defer httpResp.Body.Close()

	var read ReadViewResponse
	if err := c.unmarshalJSONBody(httpResp.Body, &read); err != nil {
		return err
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
		return utils.CreateRequestError(err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)
	if err != nil {
		return utils.SubmitRequestError(DELETE, deleteViewURL, err)
	}
	defer httpResp.Body.Close()

	return nil
}

func marshalCreateViewRequest(req *CreateViewRequest) ([]byte, error) {
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
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

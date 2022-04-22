package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *CSClient) CreateView(ctx context.Context, req *CreateViewRequest) error {
	url := fmt.Sprintf("%s/Bucket/createView", c.config.URL)
	bodyAsBytes, err := marshalCreateViewRequest(req)
	if err != nil {
		return err
	}

	_, err = c.createAndSendReq(ctx, req.AuthToken, url, POST, bodyAsBytes)
	if err != nil {
		return fmt.Errorf("Create View Failure => %s", err)
	}

	return nil
}

func (c *CSClient) ReadView(ctx context.Context, req *ReadViewRequest) (*ReadViewResponse, error) {
	var resp ReadViewResponse
	url := fmt.Sprintf("%s/Bucket/dataset/name/%s", c.config.URL, req.ID)
	httpResp, err := c.createAndSendReq(ctx, req.AuthToken, url, GET, nil)
	if err != nil {
		return nil, fmt.Errorf("Read View Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *CSClient) DeleteView(ctx context.Context, req *DeleteViewRequest) error {
	safeViewName := url.PathEscape(req.Name)
	url := fmt.Sprintf("%s/V1/%s", c.config.URL, safeViewName)
	_, err := c.createAndSendReq(ctx, req.AuthToken, url, DELETE, nil)
	if err != nil {
		return fmt.Errorf("Delete View Failure => %s", err)
	}

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

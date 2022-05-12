package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *CSClient) CreateView(ctx context.Context, req *CreateViewRequest) error {
	bodyAsBytes, err := marshalCreateViewRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/createView", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return fmt.Errorf("Create View Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func (c *CSClient) ReadView(ctx context.Context, req *ReadViewRequest) (*ReadViewResponse, error) {
	var resp ReadViewResponse
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/dataset/name/%s", c.config.URL, req.ID),
		RequestType: GET,
		AuthToken:   req.AuthToken,
	})

	if err != nil {
		return nil, fmt.Errorf("Read View Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

func (c *CSClient) DeleteView(ctx context.Context, req *DeleteViewRequest) error {
	safeViewName := url.PathEscape(req.Name)
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/V1/%s", c.config.URL, safeViewName),
		RequestType: DELETE,
		AuthToken:   req.AuthToken,
	})

	if err != nil {
		return fmt.Errorf("Delete View Failure => %s", err)
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

package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

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

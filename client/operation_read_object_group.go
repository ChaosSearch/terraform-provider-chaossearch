package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *CSClient) ReadObjGroup(x context.Context, e *ReadObjGroupReq) (*ReadObjGroupResp, error) {
	var resp ReadObjGroupResp
	if err := c.read(x, e, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *CSClient) read(x context.Context, e *ReadObjGroupReq, r *ReadObjGroupResp) error {
	url := fmt.Sprintf("%s/Bucket/dataset/name/%s", c.config.URL, e.ID)

	httpReq, err := http.NewRequestWithContext(x, GET, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV2AndDo(e.AuthToken, httpReq, nil)

	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var ReadObjectGroup ReadObjGroupResp
	if err := c.unmarshalJSONBody(httpResp.Body, &ReadObjectGroup); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response body: %s", err)
	}
	r.Format = ReadObjectGroup.Format

	r.ObjectFilter = ReadObjectGroup.ObjectFilter
	r.Interval = ReadObjectGroup.Interval
	r.Metadata = ReadObjectGroup.Metadata
	r.Options = ReadObjectGroup.Options
	r.RegionAvailability = ReadObjectGroup.RegionAvailability
	r.Public = ReadObjectGroup.Public
	r.Realtime = ReadObjectGroup.Realtime
	r.Type = ReadObjectGroup.Type
	r.Bucket = ReadObjectGroup.Bucket
	r.ContentType = ReadObjectGroup.ContentType
	r.ID = ReadObjectGroup.ID
	r.Source = ReadObjectGroup.Source

	r.Compression = ReadObjectGroup.Compression
	r.PartitionBy = ReadObjectGroup.PartitionBy
	r.Pattern = ReadObjectGroup.Pattern
	r.SourceBucket = ReadObjectGroup.SourceBucket
	r.ColumnSelection = ReadObjectGroup.ColumnSelection
	return nil
}

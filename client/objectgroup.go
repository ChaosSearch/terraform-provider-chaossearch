package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *CSClient) CreateObjectGroup(ctx context.Context, req *CreateObjectGroupRequest) error {

	url := fmt.Sprintf("%s/Bucket/createObjectGroup", c.config.URL)

	bodyAsBytes, err := marshalCreateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

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

func (c *CSClient) UpdateObjectGroup(ctx context.Context, req *UpdateObjectGroupRequest) error {

	url := fmt.Sprintf("%s/Bucket/updateObjectGroup", c.config.URL)

	bodyAsBytes, err := marshalUpdateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := c.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

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

func marshalUpdateObjectGroupRequest(req *UpdateObjectGroupRequest) ([]byte, error) {
	body := map[string]interface{}{
		"bucket":                req.Bucket,
		"indexParallelism":      req.IndexParallelism,
		"indexRetention":        req.IndexRetention,
		"targetActiveIndex":     req.TargetActiveIndex,
		"liveEventsParallelism": req.LiveEventsParallelism,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}

func (c *CSClient) DeleteObjectGroup(ctx context.Context, req *DeleteObjectGroupRequest) error {

	safeObjectGroupName := url.PathEscape(req.Name)
	deleteURL := fmt.Sprintf("%s/V1/%s", c.config.URL, safeObjectGroupName)

	httpReq, err := http.NewRequestWithContext(ctx, DELETE, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	sessionToken := req.AuthToken
	httpResp, err := c.signV2AndDo(sessionToken, httpReq, nil)

	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", POST, deleteURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)
	return nil
}

func marshalCreateObjectGroupRequest(req *CreateObjectGroupRequest) ([]byte, error) {
	body := map[string]interface{}{
		"bucket": req.Bucket,
		"source": req.Source,
		"format": map[string]interface{}{
			"_type":           req.Format.Type,
			"columnDelimiter": req.Format.ColumnDelimiter,
			"rowDelimiter":    req.Format.RowDelimiter,
			"headerRow":       req.Format.HeaderRow,
		},
		"filter": []interface{}{
			req.Filter.PrefixFilter, req.Filter.RegexFilter,
		},
		"indexRetention": map[string]interface{}{
			"forPartition": req.IndexRetention.ForPartition,
			"overall":      req.IndexRetention.Overall,
		},
		"options": map[string]interface{}{
			"ignoreIrregular": req.Options.IgnoreIrregular,
		},
		"interval": map[string]interface{}{
			"mode":   req.Interval.Mode,
			"column": req.Interval.Column,
		},
		"realtime": req.Realtime,
	}
	bodyAsBytes, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}

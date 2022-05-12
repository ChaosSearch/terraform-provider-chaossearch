package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *CSClient) CreateObjectGroup(ctx context.Context, req *CreateObjectGroupRequest) error {
	bodyAsBytes, err := marshalCreateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/createObjectGroup", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return fmt.Errorf("Create Object Group Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func (c *CSClient) ReadObjGroup(ctx context.Context, req *ReadObjGroupReq) (*ReadObjGroupResp, error) {
	var resp ReadObjGroupResp
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/dataset/name/%s", c.config.URL, req.ID),
		RequestType: GET,
		AuthToken:   req.AuthToken,
	})

	if err != nil {
		return nil, fmt.Errorf("Read Object Group Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

func (c *CSClient) UpdateObjectGroup(ctx context.Context, req *UpdateObjectGroupRequest) error {
	bodyAsBytes, err := marshalUpdateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/updateObjectGroup", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return fmt.Errorf("Update Object Group Failure => %s", err)
	}

	defer httpResp.Body.Close()
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
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

func (c *CSClient) DeleteObjectGroup(ctx context.Context, req *DeleteObjectGroupRequest) error {
	safeObjectGroupName := url.PathEscape(req.Name)
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/V1/%s", c.config.URL, safeObjectGroupName),
		RequestType: DELETE,
		AuthToken:   req.AuthToken,
	})

	if err != nil {
		return fmt.Errorf("Delete Object Group Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func marshalCreateObjectGroupRequest(req *CreateObjectGroupRequest) ([]byte, error) {
	var filter []interface{}
	var indexRetention, format, options, interval map[string]interface{}
	if req.Filter.PrefixFilter != nil {
		filter = append(filter, req.Filter.PrefixFilter)
	}

	if req.Filter.RegexFilter != nil {
		filter = append(filter, req.Filter.RegexFilter)
	}

	if req.IndexRetention != nil {
		indexRetention = map[string]interface{}{
			"forPartition": req.IndexRetention.ForPartition,
			"overall":      req.IndexRetention.Overall,
		}
	}

	if req.Format != nil {
		format = map[string]interface{}{
			"_type":           req.Format.Type,
			"columnDelimiter": req.Format.ColumnDelimiter,
			"rowDelimiter":    req.Format.RowDelimiter,
			"headerRow":       req.Format.HeaderRow,
		}
	}

	if req.Options != nil {
		options = map[string]interface{}{
			"ignoreIrregular": req.Options.IgnoreIrregular,
		}
	}

	if req.Interval != nil {
		interval = map[string]interface{}{
			"mode":   req.Interval.Mode,
			"column": req.Interval.Column,
		}
	}

	body := map[string]interface{}{
		"bucket":         req.Bucket,
		"source":         req.Source,
		"format":         format,
		"filter":         filter,
		"indexRetention": indexRetention,
		"options":        options,
		"interval":       interval,
		"realtime":       req.Realtime,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *CSClient) CreateObjectGroup(ctx context.Context, req *CreateObjectGroupRequest) error {
	url := fmt.Sprintf("%s/Bucket/createObjectGroup", c.config.URL)
	bodyAsBytes, err := marshalCreateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, req.AuthToken, url, POST, bodyAsBytes)
	if err != nil {
		return fmt.Errorf("Create Object Group Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func (c *CSClient) ReadObjGroup(ctx context.Context, req *ReadObjGroupReq) (*ReadObjGroupResp, error) {
	var resp ReadObjGroupResp
	url := fmt.Sprintf("%s/Bucket/dataset/name/%s", c.config.URL, req.ID)
	httpResp, err := c.createAndSendReq(ctx, req.AuthToken, url, GET, nil)
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
	url := fmt.Sprintf("%s/Bucket/updateObjectGroup", c.config.URL)
	bodyAsBytes, err := marshalUpdateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, req.AuthToken, url, POST, bodyAsBytes)
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
	url := fmt.Sprintf("%s/V1/%s", c.config.URL, safeObjectGroupName)
	httpResp, err := c.createAndSendReq(ctx, req.AuthToken, url, DELETE, nil)
	if err != nil {
		return fmt.Errorf("Delete Object Group Failure => %s", err)
	}

	defer httpResp.Body.Close()
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
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

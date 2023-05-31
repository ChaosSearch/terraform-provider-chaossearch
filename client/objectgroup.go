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
		Url:         fmt.Sprintf("%s/Bucket/createObjectGroup", c.Config.URL),
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

func (c *CSClient) ReadObjGroup(ctx context.Context, req *BasicRequest) (*ReadObjGroupResp, error) {
	var resp ReadObjGroupResp
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/Bucket/dataset/name/%s", c.Config.URL, req.Id),
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
		Url:         fmt.Sprintf("%s/Bucket/updateObjectGroup", c.Config.URL),
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

func (c *CSClient) DeleteObjectGroup(ctx context.Context, req *BasicRequest) error {
	safeObjectGroupName := url.PathEscape(req.Id)
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/V1/%s", c.Config.URL, safeObjectGroupName),
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
	var indexRetention,
		format,
		options,
		interval map[string]interface{}

	filters := []map[string]interface{}{}
	rangeFilters := []string{
		"lastModified",
		"size",
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
			"stripPrefix":     req.Format.StripPrefix,
			"horizontal":      req.Format.Horizontal,
		}

		if req.Format.ArrayFlattenDepth != -1 {
			format["arrayFlattenDepth"] = req.Format.ArrayFlattenDepth
		}

		if req.Format.ArraySelection != nil {
			format["arraySelection"] = req.Format.ArraySelection
		}

		if req.Format.FieldSelection != nil {
			format["fieldSelection"] = req.Format.FieldSelection
		}

		if req.Format.VerticalSelection != nil {
			format["verticalSelection"] = req.Format.VerticalSelection
		}
	}

	if req.Options != nil {
		options = map[string]interface{}{
			"ignoreIrregular": req.Options.IgnoreIrregular,
		}

		if req.Options.Compression != "" {
			options["compression"] = req.Options.Compression
		}

		if req.Options.ColTypes != nil {
			options["colTypes"] = req.Options.ColTypes
		}

		if req.Options.ColRenames != nil {
			options["colRenames"] = req.Options.ColRenames
		}

		if req.Options.ColSelection != nil {
			options["colSelection"] = req.Options.ColSelection
		}
	}

	if req.Interval != nil {
		interval = map[string]interface{}{
			"mode":   req.Interval.Mode,
			"column": req.Interval.Column,
		}
	}

	if len(req.Filter) > 0 {
		for _, filter := range req.Filter {
			filterMap := map[string]interface{}{
				"field": filter.Field,
			}

			if utils.ContainsString(rangeFilters, filter.Field) {
				filterMap["range"] = filter.Range
			} else {
				if filter.Field == "storageClass" {
					filterMap["equals"] = filter.Equals
				}
				if filter.Field == "key" && filter.Prefix != "" {
					filterMap["prefix"] = filter.Prefix
				}
				if filter.Field == "key" && filter.Regex != "" {
					filterMap["regex"] = filter.Regex
				}
			}

			filters = append(filters, filterMap)
		}
	}

	body := map[string]interface{}{
		"bucket":            req.Bucket,
		"source":            req.Source,
		"format":            format,
		"filter":            filters,
		"indexRetention":    indexRetention,
		"options":           options,
		"interval":          interval,
		"realtime":          req.Realtime,
		"targetActiveIndex": req.TargetActiveIndex,
	}

	if req.LiveEvents != "" {
		body["liveEvents"] = req.LiveEvents
	}

	if req.PartitionBy != "" {
		body["partitionBy"] = req.PartitionBy
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

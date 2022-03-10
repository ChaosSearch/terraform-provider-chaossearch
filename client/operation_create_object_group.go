package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

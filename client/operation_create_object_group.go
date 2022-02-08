package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) CreateObjectGroup(ctx context.Context, req *CreateObjectGroupRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/createObjectGroup", client.config.URL)

	bodyAsBytes, err := marshalCreateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := client.signAndDo(httpReq, bodyAsBytes)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()

	return nil
}

func marshalCreateObjectGroupRequest(req *CreateObjectGroupRequest) ([]byte, error) {
	body := map[string]interface{}{
		"bucket": req.Name,
		"source": req.SourceBucket,
		"format": map[string]interface{}{
			"_type":             req.Format,
			"horizontal":        true,
			"stripPrefix":       true,
			"arrayFlattenDepth": req.ArrayFlattenDepth,
			"keepOriginal":      req.KeepOriginal,
		},
		"indexRetention": req.IndexRetention,
		"options": map[string]interface{}{
			"ignoreIrregular": true,
		},
		"interval": map[string]interface{}{
			"mode":   0,
			"column": 0,
		},
	}

	if len(req.ColumnRenames) > 0 {
		var options = body["options"].(map[string]interface{})
		options["colRenames"] = req.ColumnRenames
	}

	if len(req.ColumnSelection) > 0 {
		var options = body["options"].(map[string]interface{})
		// @example
		//"colSelection": [
		//	{
		//	"includes": [
		//		"orig._originalSource",
		//		"attrs.version",
		//		"line.message",
		//		"line.correlation_id",
		//		"Timestamp"
		//	],
		//	"type": "whitelist"
		//	}
		//],
		options["colSelection"] = []map[string]interface{}{req.ColumnSelection}
	}

	if req.Compression != "" {
		var options = body["options"].(map[string]interface{})
		options["compression"] = req.Compression
	}

	if req.FilterJSON != "" {
		filter := make(map[string]interface{})
		if err := json.Unmarshal([]byte(req.FilterJSON), &filter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON string: %s %s", req.FilterJSON, err)
		}
		body["filter"] = filter
	}

	if req.LiveEventsSqsArn != "" {
		body["liveEvents"] = req.LiveEventsSqsArn
	}

	if req.PartitionBy != "" {
		body["partitionBy"] = req.PartitionBy
	}

	if req.Format == "LOG" {
		var format = body["format"].(map[string]interface{})
		format["pattern"] = req.Pattern
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}

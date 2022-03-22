package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

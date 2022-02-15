package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (client *Client) CreateView(ctx context.Context, req *CreateViewRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/createView", client.config.URL)
	log.Debug("Url-->", url)
	log.Debug("req-->", req)
	bodyAsBytes, err := marshalCreateViewRequest(req)
	if err != nil {
		return err
	}

	log.Warn("bodyAsBytes--1")
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	log.Debug(" adding headers...")
	httpReq.Header.Add("Content-Type", "text/plain")

	// var sessionToken = req.AuthToken
	//httpReq.Header.Add("x-amz-security-token", req.AuthToken)

	log.Warn("httpReq-->", httpReq)
	httpResp, err := client.signAndDo(httpReq, bodyAsBytes)
	log.Warn("httpResp-->", httpResp)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()

	return nil
}

func marshalCreateViewRequest(req *CreateViewRequest) ([]byte, error) {
	log.Warn("req.Sources----", req.Sources)
	body := map[string]interface{}{
		"bucket":    req.Bucket,
		"sources":   req.Sources,
		"indexPattern":    req.IndexPattern,
		"overwrite": req.Overwrite,
		"caseInsensitive": req.CaseInsensitive,
		"indexRetention": req.IndexRetention,
		"timeFieldName":req.TimeFieldName,
		"transforms":req.Transforms,
		"filter":req.Filter,
		
		//"cacheable": req.Cacheable,
		// "horizontal":  true,
		// "stripPrefix": true,
		// "arrayFlattenDepth": req.ArrayFlattenDepth,
		// "options": map[string]interface{}{
		// 	"ignoreIrregular": true,
		// },
		// "interval": map[string]interface{}{
		// 	"mode":   0,
		// 	"column": 0,
		// },
	}
	log.Warn("body----", body)

	// if req.FilterJSON != "" {
	// 	filter := make(map[string]interface{})
	// 	if err := json.Unmarshal([]byte(req.FilterJSON), &filter); err != nil {
	// 		return nil, fmt.Errorf("failed to unmarshal JSON string: %s %s", req.FilterJSON, err)
	// 	}
	// 	body["filter"] = filter
	// }

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	log.Warn("marshalling--3")
	return bodyAsBytes, nil
}

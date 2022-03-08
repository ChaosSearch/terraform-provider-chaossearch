package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (csClient *CSClient) CreateView(ctx context.Context, req *CreateViewRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/createView", csClient.config.URL)
	log.Debug("Url-->", url)
	log.Debug("req-->", req)
	bodyAsBytes, err := marshalCreateViewRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	log.Debug(" adding headers...")
	log.Warn("httpReq-->", httpReq)

	httpResp, err := csClient.signV2AndDo(req.AuthToken, httpReq, bodyAsBytes)

	log.Warn("httpResp-->", httpResp)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	return nil
}

func marshalCreateViewRequest(req *CreateViewRequest) ([]byte, error) {
	log.Debug("req.Sources-->", req.Sources)
	body := map[string]interface{}{
		"bucket":          req.Bucket,
		"sources":         req.Sources,
		"indexPattern":    req.IndexPattern,
		"overwrite":       req.Overwrite,
		"caseInsensitive": req.CaseInsensitive,
		"indexRetention":  req.IndexRetention,
		"timeFieldName":   req.TimeFieldName,
		"transforms":      req.Transforms,
		"filter":          req.FilterPredicate,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

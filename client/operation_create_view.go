package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "io/ioutil"
	// "strings"
	log "github.com/sirupsen/logrus"
)

// func (client *Client) CallCreateView()error{
// 	url := "https://ap-south-1-aeternum.chaossearch.io/Bucket/createView"
// 	method := "POST"
//   log.Debug("calling create view.")
// 	payload := strings.NewReader(`{
// 	  "bucket":"test-view-dinesh-123",
// 	  "sources":[],
// 	  "indexPattern":".*",
// 	  "caseInsensitive":false,
// 	  "cacheable":false,
// 	  "overwrite":false,
// 	  "indexRetention":-1,
// 	  "transforms":[]
//   }`)

// 	client := &http.Client {
// 	}
// 	req, err := http.NewRequest(method, url, payload)

// 	if err != nil {
// 	  fmt.Println(err)
// 	  return err
// 	}
// 	req.Header.Add("x-amz-security-token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJTZWNyZXRBY2Nlc3NLZXkiOiJkRlQ0ZncrQkxhRklzd0dWUUp3bHBqWThEdGJYc2RPRTYyMlVzV1M1IiwiZXBvY2giOjE2NDQxMTUwNzM4NzEsImV4dGVybmFsX2lkIjoiYmU0YWViNTMtMjFkNS00OTAyLTg2MmMtOWM5YTE3YWQ2Njc1IiwiQ1NSb2xlIjoidXNlciIsIlJlcXVlc3RVc2VyIjoic2VydmljZV91c2VyQGNoYW9zc2VhcmNoLmNvbSIsImF3c19hY2NvdW50X2lkIjoiNzY3Mzk2NjcxNzMyIiwiVXNlcm5hbWUiOiJzZXJ2aWNlX3VzZXJAY2hhb3NzZWFyY2guY29tIiwiQWNjZXNzS2V5SWQiOiJPS1dHM1lPTE9FTkw2V09ZSE40USIsIlNlcnZpY2VUeXBlIjoiUHJlbWl1bSIsIm5hbWUiOiJzZXJ2aWNlX3VzZXIiLCJQcmltYXJ5VXNlciI6ImFldGVybnVtQGNoYW9zc2VhcmNoLmNvbSIsIndvcmtlci1xdW90YV9hcC1zb3V0aC0xIjoyNDAsIkxvZ2luVHlwZSI6ImFsaWFzIiwiZXhwIjoxNjQ0MzA1ODg1LCJpYXQiOjE2NDQyMTk0ODUsImVtYWlsIjoic2VydmljZV91c2VyQGNoYW9zc2VhcmNoLmNvbSIsImp0aSI6ImFlZmQyNzE1LWNhYjctNDNmMS1iODE3LWIwYzY4MGI4ZDg2MiJ9.kRkXWDjo0VJXuLogXA1VsaiIllyfhgbdlPGgNPryzEo")
// 	req.Header.Add("x-amz-chaossumo-route-token", "login")
// 	req.Header.Add("X-Amz-Content-Sha256", "beaead3198f7da1e70d03ab969765e0821b24fc913697e929e726aeaebf0eba3")
// 	req.Header.Add("X-Amz-Date", "20220207T074248Z")
// 	req.Header.Add("Authorization", "AWS4-HMAC-SHA256 Credential=LCE8T6HRFGJI3ZKBGMGD/20220207/ap-south-1/s3/aws4_request, SignedHeaders=host;x-amz-chaossumo-route-token;x-amz-content-sha256;x-amz-date, Signature=0367288f4acbd838ec5ed706e77e2cbcfd4da64b20877ec7fe4d41fa942d3435")
// 	req.Header.Add("Content-Type", "text/plain")

// 	res, err := client.Do(req)
// 	if err != nil {
// 	  fmt.Println(err)
// 	  return err
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 	  fmt.Println(err)
// 	  return err
// 	}
// 	fmt.Println(string(body))
// }

func (client *Client) CreateView(ctx context.Context, req *CreateViewRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/createView", client.config.URL)
	log.Debug("Url-->", url)
	log.Debug("req-->", req)
	// client.config.
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

	// todo auth token unset as token validatio failure
	// var sessionToken = req.AuthToken
	// httpReq.Header.Add("x-amz-security-token", sessionToken)

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
		"bucket":    req.Name,
		"sources":   req.Sources,
		"cacheable": req.Cacheable,
		"overwrite": req.Overwrite,
		// "transforms":req.Transforms,
		// "horizontal":  true,
		// "stripPrefix": true,
		"indexPattern":    req.Pattern,
		"caseInsensitive": req.CaseInsensitive,
		// "arrayFlattenDepth": req.ArrayFlattenDepth,
		"indexRetention": req.IndexRetention,
		// "options": map[string]interface{}{
		// 	"ignoreIrregular": true,
		// },
		// "interval": map[string]interface{}{
		// 	"mode":   0,
		// 	"column": 0,
		// },
	}
	log.Warn("body----", body)

	if req.FilterJSON != "" {
		filter := make(map[string]interface{})
		if err := json.Unmarshal([]byte(req.FilterJSON), &filter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON string: %s %s", req.FilterJSON, err)
		}
		body["filter"] = filter
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	log.Warn("marshalling--3")
	return bodyAsBytes, nil
}

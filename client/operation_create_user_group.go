package client

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (csClient *CSClient) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/user/groups", csClient.config.URL)
	log.Debug("Url-->", url)
	log.Debug("req-->", req)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return err
	}
	log.Debug("method-->", method)
	log.Debug("bodyAsBytes-->", bodyAsBytes)

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	log.Debug(" adding headers...")
	httpReq.Header.Add("Content-Type", "text/plain")

	var sessionToken = req.AuthToken
	//httpReq.Header.Add("x-amz-security-token", req.AuthToken)

	log.Warn("httpReq-->", httpReq)
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, bodyAsBytes)
	//httpResp, err := client.signV4AndDo(httpReq, bodyAsBytes)
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

func marshalCreateUserGroupRequest(req *CreateUserGroupRequest) ([]byte, error) {
	// todo to be implemented
	return nil, nil
}

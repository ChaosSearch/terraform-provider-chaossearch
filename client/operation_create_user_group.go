package client

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"bytes"
)

func (client *Client) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/user/groups", client.config.URL)
	log.Debug("Url-->", url)
	log.Debug("req-->", req)
	bodyAsBytes, err := marshalCreateUserGroupRequest(req)
	if err != nil {
		return err
	}
	log.Debug("method-->",method)
	log.Debug("bodyAsBytes-->",bodyAsBytes)

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

func marshalCreateUserGroupRequest(req *CreateUserGroupRequest) ([]byte, error) {
	return nil,nil
}

package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
)

func (c *CSClient) CreateMonitor(
	ctx context.Context,
	req *CreateMonitorRequest,
) (*CreateMonitorResponse, error) {
	var resp CreateMonitorResponse
	bodyAsBytes, err := marshalCreateMonitorRequest(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/monitors", c.Config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
		Headers: map[string]string{
			"Cookie": fmt.Sprintf("chaossumo_session_token=%s", req.AuthToken),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Create Monitor Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

func (c *CSClient) UpdateMonitor(
	ctx context.Context,
	req *CreateMonitorRequest,
) error {
	bodyAsBytes, err := marshalCreateMonitorRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/monitors/%s", c.Config.URL, req.Id),
		RequestType: PUT,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
		Headers: map[string]string{
			"Cookie": fmt.Sprintf("chaossumo_session_token=%s", req.AuthToken),
		},
	})

	if err != nil {
		return fmt.Errorf("Update Monitor Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func (c *CSClient) DeleteMonitor(ctx context.Context, req *BasicRequest) error {
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/monitors/%s", c.Config.URL, req.Id),
		RequestType: DELETE,
		AuthToken:   req.AuthToken,
		Headers: map[string]string{
			"Cookie": fmt.Sprintf("chaossumo_session_token=%s", req.AuthToken),
		},
	})

	if err != nil {
		return fmt.Errorf("Delete Monitor Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func marshalCreateMonitorRequest(req *CreateMonitorRequest) ([]byte, error) {
	req.UIMetadata = map[string]interface{}{}
	bodyAsBytes, err := json.Marshal(req)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

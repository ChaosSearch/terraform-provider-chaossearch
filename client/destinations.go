package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
)

func (c *CSClient) CreateDestination(
	ctx context.Context,
	req *CreateDestinationRequest,
) (*CreateDestinationResponse, error) {
	var resp CreateDestinationResponse
	bodyAsBytes, err := marshalCreateDestinationRequest(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/destinations", c.Config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return nil, fmt.Errorf("Create Destination Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	return &resp, nil
}

func (c *CSClient) ReadDestination(ctx context.Context, req *BasicRequest) (*ReadDestinationResponse, error) {
	var resp ReadDestinationResponse
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/destinations/%s", c.Config.URL, req.Id),
		RequestType: GET,
		AuthToken:   req.AuthToken,
	})

	if err != nil {
		return nil, fmt.Errorf("Read Destination Failure => %s", err)
	}

	defer httpResp.Body.Close()
	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *CSClient) UpdateDestination(
	ctx context.Context,
	req *CreateDestinationRequest,
) error {
	var resp CreateDestinationResponse
	bodyAsBytes, err := marshalCreateDestinationRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/destinations/%s", c.Config.URL, req.Id),
		RequestType: PUT,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return fmt.Errorf("Update Destination Failure => %s", err)
	}

	if err := c.unmarshalJSONBody(httpResp.Body, &resp); err != nil {
		return err
	}

	defer httpResp.Body.Close()
	return nil
}

func (c *CSClient) DeleteDestination(ctx context.Context, req BasicRequest) error {
	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/kibana/api/alerting/destinations/%s", c.Config.URL, req.Id),
		RequestType: DELETE,
		AuthToken:   req.AuthToken,
	})

	if err != nil {
		return fmt.Errorf("Delete Destination Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func marshalCreateDestinationRequest(req *CreateDestinationRequest) ([]byte, error) {
	body := map[string]interface{}{
		"name": req.Name,
		"type": req.Type,
	}

	if req.Slack != nil {
		body["slack"] = map[string]interface{}{
			"url": req.Slack.Url,
		}
	}

	if req.CustomWebhook != nil {
		webhookBody := map[string]interface{}{
			"scheme": "",
			"method": req.CustomWebhook.Method,
			"url":    req.CustomWebhook.Url,
			"host":   req.CustomWebhook.Host,
			"port":   req.CustomWebhook.Port,
			"path":   req.CustomWebhook.Path,
		}

		if len(req.CustomWebhook.HeaderParams) > 0 {
			webhookBody["header_params"] = req.CustomWebhook.HeaderParams
		}

		if len(req.CustomWebhook.QueryParams) > 0 {
			webhookBody["query_params"] = req.CustomWebhook.QueryParams
		}

		body["custom_webhook"] = webhookBody
	}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

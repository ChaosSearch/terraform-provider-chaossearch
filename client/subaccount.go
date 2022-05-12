package client

import (
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"fmt"
)

func (c *CSClient) CreateSubAccount(ctx context.Context, req *CreateSubAccountRequest) error {
	bodyAsBytes, err := marshalCreateSubAccountRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/createSubAccount", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return fmt.Errorf("Create SubAccount Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func (c *CSClient) DeleteSubAccount(ctx context.Context, req *DeleteSubAccountRequest) error {
	bodyAsBytes, err := marshalDeleteSubAccountRequest(req)
	if err != nil {
		return err
	}

	httpResp, err := c.createAndSendReq(ctx, ClientRequest{
		Url:         fmt.Sprintf("%s/user/deleteSubAccount", c.config.URL),
		RequestType: POST,
		AuthToken:   req.AuthToken,
		Body:        bodyAsBytes,
	})

	if err != nil {
		return fmt.Errorf("Delete SubAccount Failure => %s", err)
	}

	defer httpResp.Body.Close()
	return nil
}

func marshalCreateSubAccountRequest(req *CreateSubAccountRequest) ([]byte, error) {
	body := map[string]interface{}{
		"UserInfoBlock": map[string]interface{}{
			"Username": req.UserInfoBlock.Username,
			"FullName": req.UserInfoBlock.FullName,
			"Email":    req.UserInfoBlock.Email,
		},
		"GroupIds": req.GroupIds,
		"Password": req.Password,
		"Hocon":    req.HoCon,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

func marshalDeleteSubAccountRequest(req *DeleteSubAccountRequest) ([]byte, error) {
	body := map[string]interface{}{
		"Username": req.Username,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

package client

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

func (csClient *CSClient) ListUsers(ctx context.Context, authToken string) (*ListUsersResponse, error) {
	url := fmt.Sprintf("%s/user/manifest", csClient.config.URL)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := csClient.signV2AndDo(authToken, httpReq, nil)
	//httpResp, err := client.signV4AndDo(httpReq, nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	b, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		panic(err)
	}

	log.Debug("user req body--->>", b)

	var resp ListUsersResponse
	if err := csClient.unmarshalXMLBody(httpResp.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML response body: %s", err)
	}

	return &resp, nil
}

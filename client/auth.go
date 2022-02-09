package client

import (
	"bytes"
	"context"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (client *Client) Auth(ctx context.Context) (token string, err error) {

	url := fmt.Sprintf("%s/user/login", client.config.URL)
	method := "POST"
	login_ := client.Login
	
	log.Warn("url--", url)

	log.Warn("username--", login_.Username)
	log.Warn("password--", login_.Password)
	log.Warn("parentuserid--", login_.ParentUserId)

	bodyAsBytes, err := marshalLoginRequest(login_)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %s", err)
	}
	log.Debug(" adding headers...")

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-amz-chaossumo-route-token", "login")
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	return string(body), nil
}

func marshalLoginRequest(req *Login) ([]byte, error) {
	log.Warn("req.Sources----", req.Username)

	body := map[string]interface{}{
		"Username":  req.Username,
		"Password":  req.Password,
		"ParentUid": req.ParentUserId,
	}
	log.Warn("body---->>", body)

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	log.Warn("marshalling--3")
	return bodyAsBytes, nil
}

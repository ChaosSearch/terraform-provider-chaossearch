package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

func (csClient *CSClient) Auth(ctx context.Context) (token string, err error) {

	url := fmt.Sprintf("%s/user/login", csClient.config.URL)
	method := "POST"
	login_ := csClient.Login

	log.Warn("url-->", url)
	log.Warn("username-->", login_.Username)
	log.Warn("parentuserid-->", login_.ParentUserId)

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

	res, err := csClient.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(res.Body)
	// TODO add a status call once successful login to ensure that the user is actually deployed
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	return string(body), nil
}

func marshalLoginRequest(req *Login) ([]byte, error) {
	var body map[string]interface{}
	if len(req.ParentUserId) == 0 {
		body = map[string]interface{}{
			"Username": req.Username,
			"Password": req.Password,
		}
	} else {
		body = map[string]interface{}{
			"Username":  req.Username,
			"Password":  req.Password,
			"ParentUid": req.ParentUserId,
		}
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

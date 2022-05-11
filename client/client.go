package client

import (
	"bytes"
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Configuration struct {
	URL             string
	AccessKeyID     string
	SecretAccessKey string
	AWSServiceName  string
	Region          string
	Login           *Login
}

type Login struct {
	Username     string
	Password     string
	ParentUserID string `json:"ParentUserID,omitempty"`
}

type CSClient struct {
	config     *Configuration
	httpClient *http.Client
	userAgent  string
	Login      *Login
}

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

func NewConfiguration(url, accessKeyID, secretAccessKey, region string) *Configuration {
	return &Configuration{
		AWSServiceName:  "s3",
		URL:             url,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}
}

func NewClient(config *Configuration, login *Login) *CSClient {
	binaryName := os.Getenv("BINARY")
	hostName := os.Getenv("HOSTNAME")
	version := os.Getenv("VERSION")
	namespace := os.Getenv("NAMESPACE")
	userAgent := fmt.Sprintf("%s/%s %s/%s/%s", binaryName, version, hostName, namespace, binaryName)

	return &CSClient{
		config:     config,
		httpClient: http.DefaultClient,
		userAgent:  userAgent,
		Login:      login,
	}
}

func (c *CSClient) Auth(ctx context.Context) (token string, err error) {
	url := fmt.Sprintf("%s/user/login", c.config.URL)
	login := c.Login
	bodyAsBytes, err := marshalLoginRequest(login)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return "", utils.CreateRequestError(err)
	}

	req.Header.Add("x-amz-chaossumo-route-token", "login")
	req.Header.Add("Content-Type", "text/plain")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", utils.SubmitRequestError(POST, url, err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", utils.ReadResponseError(err)
	}

	return string(body), nil
}

func marshalLoginRequest(req *Login) ([]byte, error) {
	var body map[string]interface{}
	if len(req.ParentUserID) == 0 {
		body = map[string]interface{}{
			"Username": req.Username,
			"Password": req.Password,
		}
	} else {
		body = map[string]interface{}{
			"Username":  req.Username,
			"Password":  req.Password,
			"ParentUid": req.ParentUserID,
		}
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, utils.MarshalJsonError(err)
	}

	return bodyAsBytes, nil
}

func (client *CSClient) createAndSendReq(
	ctx context.Context,
	request ClientRequest,
) (*http.Response, error) {
	httpReq, err := request.constructRequest(ctx)
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := client.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute request => Error: %s", err)
	}

	if httpResp.Body == http.NoBody {
		_, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return nil, utils.ReadResponseError(err)
		}
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		respAsBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return nil, utils.ReadResponseError(err)
		}

		return nil, utils.ResponseCodeError(httpResp, httpReq, request.Body, respAsBytes)
	}

	defer func(httpReq *http.Request) {
		if httpReq.Body != nil {
			httpReq.Body.Close()
		}
	}(httpReq)

	return httpResp, nil
}

func (c *CSClient) unmarshalJSONBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return utils.ReadResponseError(err)
	}

	if err := json.Unmarshal(bodyAsBytes, v); err != nil {
		return utils.UnmarshalJsonError(err)
	}

	return nil
}

func (c *CSClient) unmarshalXMLBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return utils.ReadResponseError(err)
	}

	if err := xml.Unmarshal(bodyAsBytes, v); err != nil {
		return utils.UnmarshalXmlError(err)
	}

	return nil
}

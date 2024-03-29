package client

import (
	"bytes"
	"context"
	"cs-tf-provider/client/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Configuration struct {
	URL             string
	AccessKeyID     string
	SecretAccessKey string
	AWSServiceName  string
	Region          string
	KeyAuthEnabled  bool
	UserID          string
	RetryCount      int
}

type Login struct {
	Username     string
	Password     string
	ParentUserID string `json:"ParentUserID,omitempty"`
}

type CSClient struct {
	Config     *Configuration
	httpClient *http.Client
	userAgent  string
	Login      *Login
	RetryCount int
}

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

func NewConfiguration(url, accessKeyID, secretAccessKey, region string, keyAuth bool) *Configuration {
	return &Configuration{
		AWSServiceName:  "s3",
		URL:             url,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		KeyAuthEnabled:  keyAuth,
	}
}

func NewClient(config *Configuration, login *Login) *CSClient {
	binaryName := os.Getenv("BINARY")
	hostName := os.Getenv("HOSTNAME")
	version := os.Getenv("VERSION")
	namespace := os.Getenv("NAMESPACE")
	userAgent := fmt.Sprintf("%s/%s %s/%s/%s", binaryName, version, hostName, namespace, binaryName)

	return &CSClient{
		Config:     config,
		httpClient: http.DefaultClient,
		userAgent:  userAgent,
		Login:      login,
		RetryCount: config.RetryCount,
	}
}

func (c *CSClient) Auth(ctx context.Context) (token string, err error) {
	if !c.Config.KeyAuthEnabled {
		url := fmt.Sprintf("%s/user/login", c.Config.URL)
		bodyAsBytes, err := marshalLoginRequest(c.Login)
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

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", utils.ReadResponseError(err)
		}

		return string(body), nil
	}

	return "", nil
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
	var httpResp *http.Response
	var backoffSchedule = client.fibBackoff()
	httpReq, err := request.constructRequest(ctx, *client.Config)
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err = client.httpClient.Do(httpReq)
	if httpResp == nil || (httpResp.StatusCode < 200 || httpResp.StatusCode >= 300) || err != nil {
		for _, backoff := range backoffSchedule {
			retryResp, err := client.httpClient.Do(httpReq)
			if err == nil && (retryResp.StatusCode >= 200 && retryResp.StatusCode < 300) {
				defer func(httpReq *http.Request) {
					if httpReq.Body != nil {
						httpResp = retryResp
						httpReq.Body.Close()
					}
				}(httpReq)

				break
			}

			time.Sleep(backoff)
		}
	}

	if err != nil {
		respAsBytes, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf(
			"Failed to execute request (Attempts exceeded) => StatusCode: %v \n Error: %s \n Body: %s",
			httpResp.StatusCode, err, respAsBytes,
		)
	}

	if httpResp.Body == http.NoBody {
		_, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, utils.ReadResponseError(err)
		}
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		respAsBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, utils.ReadResponseError(err)
		}

		return nil, utils.ResponseCodeError(httpResp, httpReq, request.Body, respAsBytes)
	}

	return httpResp, nil
}

func (client *CSClient) fibBackoff() []time.Duration {
	backoff := make([]time.Duration, client.RetryCount)
	backoff[0], backoff[1] = 0*time.Second, 1*time.Second
	for i := 2; i < client.RetryCount; i++ {
		backoff[i] = backoff[i-1] + backoff[i-2]
	}

	return backoff
}

func (c *CSClient) unmarshalJSONBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		return utils.ReadResponseError(err)
	}

	if err := json.Unmarshal(bodyAsBytes, v); err != nil {
		return fmt.Errorf("Error %v, Body: %s", utils.UnmarshalJsonError(err), bodyAsBytes)
	}

	return nil
}

func (c *CSClient) unmarshalXMLBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		return utils.ReadResponseError(err)
	}

	if err := xml.Unmarshal(bodyAsBytes, v); err != nil {
		return utils.UnmarshalXmlError(err)
	}

	return nil
}

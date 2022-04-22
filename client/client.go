package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"cs-tf-provider/client/utils"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func (c *CSClient) Set(data *schema.ResourceData, key string, value interface{}) {
	err := data.Set(key, value)
	if err != nil {
		return
	}
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{
		AWSServiceName: "s3",
	}

	return cfg
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

	// TODO add a status call once successful login to ensure that the user is actually deployed
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
	authToken string,
	url string,
	requestType string,
	body []byte,
) (*http.Response, error) {
	var httpReq *http.Request
	var err error
	if len(body) == 0 {
		httpReq, err = http.NewRequestWithContext(ctx, requestType, url, nil)
	} else {
		httpReq, err = http.NewRequestWithContext(ctx, requestType, url, bytes.NewReader(body))
	}

	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	httpResp, err := client.signV2AndDo(authToken, httpReq, body)
	if err != nil {
		return nil, utils.SubmitRequestError(requestType, url, err)
	}

	defer httpReq.Body.Close()
	return httpResp, nil
}

func (c *CSClient) signV2AndDo(tokenValue string, req *http.Request, bodyAsBytes []byte) (*http.Response, error) {
	var routeToken string

	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenValue, claims)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JWT => Error: %s", err)
	}

	if isAdminAPI(req.URL.Path) {
		routeToken = "login"
	} else {
		routeToken = claims["external_id"].(string)
	}

	dateTime := time.Now().UTC().String()
	msg := fmt.Sprintf("%s\n\n", req.Method) +
		"application/json\n\n" +
		fmt.Sprintf("x-amz-chaossumo-route-token:%s\n", routeToken) +
		fmt.Sprintf("x-amz-date:%s\n", dateTime) +
		fmt.Sprintf("x-amz-security-token:%s\n", tokenValue) +
		req.URL.Path

	auth := fmt.Sprintf(
		"AWS %s:%s",
		claims["AccessKeyId"].(string),
		generateSignature(claims["SecretAccessKey"].(string), msg),
	)

	req.Header.Add("Authorization", auth)
	req.Header.Add("x-amz-cs3-authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-amz-chaossumo-route-token", routeToken)
	req.Header.Add("x-amz-security-token", tokenValue)
	req.Header.Add("X-Amz-Date", dateTime)

	resp, e := c.httpClient.Do(req)
	if e != nil {
		return nil, fmt.Errorf("Failed to execute request => Error: %s", e)
	}

	if resp.Body == http.NoBody {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, utils.ReadResponseError(err)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respAsBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, utils.ReadResponseError(err)
		}

		return nil, utils.ResponseCodeError(resp, req, bodyAsBytes, respAsBytes)
	}

	return resp, nil
}

func generateSignature(secretToken string, payloadBody string) string {
	keyForSign := []byte(secretToken)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(payloadBody))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func isAdminAPI(url string) bool {
	return strings.HasSuffix(url, "/createSubAccount") ||
		strings.HasSuffix(url, "/deleteSubAccount") ||
		strings.HasSuffix(url, "/user/manifest") ||
		strings.HasSuffix(url, "/user/groups") ||
		strings.Contains(url, "/user/group/")
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

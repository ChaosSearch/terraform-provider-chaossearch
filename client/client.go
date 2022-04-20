package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
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

	log "github.com/sirupsen/logrus"
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

	log.Warn("url-->", url)
	log.Warn("username-->", login.Username)
	log.Warn("parent user id-->", login.ParentUserID)

	bodyAsBytes, err := marshalLoginRequest(login)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, POST, url, bytes.NewReader(bodyAsBytes))
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

	res, err := c.httpClient.Do(req)
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
		return nil, err
	}
	return bodyAsBytes, nil
}

func (c *CSClient) signV2AndDo(tokenValue string, req *http.Request, bodyAsBytes []byte) (*http.Response, error) {
	log.Debug("------- AWS V2 Sign Starts------")

	claims := jwt.MapClaims{}
	log.Debug("token-->", tokenValue)

	_, _, err := new(jwt.Parser).ParseUnverified(tokenValue, claims)

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}

	accessKey := claims["AccessKeyId"].(string)
	secretAccessKey := claims["SecretAccessKey"].(string)
	externalID := claims["external_id"].(string)
	dateTime := time.Now().UTC().String()

	req.Header.Add("Content-Type", "application/json")

	var routeToken string
	if isAdminAPI(req.URL.Path) {
		routeToken = "login"
	} else {
		routeToken = externalID
	}

	req.Header.Add("x-amz-chaossumo-route-token", routeToken)
	req.Header.Add("x-amz-security-token", tokenValue)
	req.Header.Add("X-Amz-Date", dateTime)

	log.Debug("headers -->", req.Header)
	log.Debug("req.URL.Path -->", req.URL.Path)

	msgLines := []string{
		req.Method, "",
		"application/json", "",
		"x-amz-chaossumo-route-token:" + routeToken,
		"x-amz-date:" + dateTime,
		"x-amz-security-token:" + tokenValue,
		req.URL.Path,
	}

	msg := strings.Join(msgLines, "\n")
	log.Debug("msg-->", msg)

	auth := "AWS " + accessKey + ":" + generateSignature(secretAccessKey, msg)
	log.Debug("auth-->", auth)

	req.Header.Add("Authorization", auth)
	req.Header.Add("x-amz-cs3-authorization", auth)
	log.Debug("req.Header-->", req.Header)

	resp, e := c.httpClient.Do(req)
	if e != nil {
		return nil, fmt.Errorf("failed to execute request: %s", e)
	}

	if req.Body != nil {
		if req.Body == http.NoBody {
			reqBody, _ := ioutil.ReadAll(req.Body)
			log.Info("Request Body -->", string(reqBody))
		}
	}

	if resp.Body == http.NoBody {
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			return nil, fmt.Errorf("failed to read response body--->: %s", err1)
		}
		log.Info("Response Body -->", string(body))
	}

	log.Warn("Got response:Status code: ", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respAsBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}
		return nil, fmt.Errorf(
			"expected a 2xx status code, but got %d.\nMethod: %s\nURL: %s\nRequest body: %s\nResponse body: %s",
			resp.StatusCode, req.Method, req.URL, bodyAsBytes, respAsBytes)
	}
	log.Debug("------- AWS V2 Sign Ends------")
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
		return fmt.Errorf("failed to read body: %s", err)
	}
	log.Printf("Unmarshalling JSON:-->%s<--", bodyAsBytes)
	if err := json.Unmarshal(bodyAsBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err)
	}
	return nil
}

func (c *CSClient) unmarshalXMLBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return fmt.Errorf("failed to read body: %s", err)
	}
	//log.Warn("Unmarshalling XML:", bodyAsBytes)
	if err := xml.Unmarshal(bodyAsBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal XML: %s", err)
	}
	return nil
}

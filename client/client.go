package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"github.com/dgrijalva/jwt-go"

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

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	log "github.com/sirupsen/logrus"
	// v2 "github.com/aws/aws-sdk-go/private/signer/v2"
)

type CSClient struct {
	config     *Configuration
	httpClient *http.Client
	userAgent  string
	Login      *Login
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

/*  AWS V4 Authentication
Not using with current API implementations
*/
func (csClient *CSClient) signV4AndDo(req *http.Request, bodyAsBytes []byte) (*http.Response, error) {
	var bodyReader io.ReadSeeker
	if bodyAsBytes == nil {
		bodyReader = nil
	} else {
		bodyReader = bytes.NewReader(bodyAsBytes)
	}

	var sessionToken string
	staticCredentials := credentials.NewStaticCredentials(csClient.config.AccessKeyID, csClient.config.SecretAccessKey, sessionToken)
	_, err := v4.NewSigner(staticCredentials).Sign(req, bodyReader, csClient.config.AWSServiceName, csClient.config.Region, time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %s", err)
	}

	log.Warn("Sending request:\nMethod: %s\nURL: %s\nBody: %s", req.Method, req.URL, bodyAsBytes)
	log.Warn("req--------->", req)
	resp, err := csClient.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}

	log.Warn("Got response:\nStatus code: %d", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respAsBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}
		return nil, fmt.Errorf(
			"expected a 2xx status code, but got %d.\nMethod: %s\nURL: %s\nRequest body: %s\nResponse body: %s",
			resp.StatusCode, req.Method, req.URL, bodyAsBytes, respAsBytes)
	}

	return resp, nil
}

func (csClient *CSClient) signV2AndDo(tokenValue string, req *http.Request, bodyAsBytes []byte) (*http.Response, error) {
	log.Debug("------- AWS V2 Sign Starts------")

	claims := jwt.MapClaims{}
	log.Debug("token-->>", tokenValue)

	_, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})

	log.Debug("token err---->", err)

	accessKey := claims["AccessKeyId"].(string)
	secretAccessKey := claims["SecretAccessKey"].(string)
	externalId := claims["external_id"].(string)
	dateTime := time.Now().UTC().String()

	req.Header.Add("Content-Type", "application/json")
	var routeToken string

	if isAdminApi(req.URL.Path) {
		routeToken = "login"
	} else {
		routeToken = externalId
	}
	req.Header.Add("x-amz-chaossumo-route-token", routeToken)

	req.Header.Add("x-amz-security-token", tokenValue)
	req.Header.Add("X-Amz-Date", dateTime)

	log.Debug("headers-->", req.Header)

	msgLines := []string{
		req.Method, "",
		"application/json", "",
		"x-amz-chaossumo-route-token:" + routeToken,
		"x-amz-date:" + dateTime,
		"x-amz-security-token:" + tokenValue,
		req.URL.Path,
	}

	msg := strings.Join(msgLines, "\n")
	log.Debug("msg---->", msg)

	signature := generateSignature(secretAccessKey, msg)
	log.Debug("signature---->", signature)

	auth := "AWS " + accessKey + ":" + signature
	log.Debug("auth---->", auth)

	req.Header.Add("Authorization", auth)
	req.Header.Add("x-amz-cs3-authorization", auth)
	log.Debug("req.Header-->", req.Header)

	for key, val := range req.Header {
		log.Debug("Header -->", key, "  value -->", val)
	}
	log.Debug("req.GetBody-->", req.GetBody)

	resp, e := csClient.httpClient.Do(req)
	if e != nil {
		return nil, fmt.Errorf("failed to execute request: %s", e)
	}

	if req != nil && req.Body != nil {
		reqBody, _ := ioutil.ReadAll(req.Body)
		log.Info("Request Body -->", string(reqBody))
	}
	if resp.Body != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Info("Response Body -->", string(body))
	}

	log.Warn("Got response:\nStatus code: %d", resp.StatusCode)

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

func isAdminApi(url string) bool {
	return strings.HasSuffix(url, "/createSubAccount") ||
		strings.HasSuffix(url, "/deleteSubAccount") ||
		strings.HasSuffix(url, "/user/groups")

}

func (csClient *CSClient) unmarshalJSONBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return fmt.Errorf("failed to read body: %s", err)
	}
	log.Printf("Unmarshalling JSON: %s\n", bodyAsBytes)
	if err := json.Unmarshal(bodyAsBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err)
	}
	return nil
}

func (csClient *CSClient) unmarshalXMLBody(bodyReader io.Reader, v interface{}) error {
	bodyAsBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return fmt.Errorf("failed to read body: %s", err)
	}
	log.Warn("Unmarshalling XML: %s\n", bodyAsBytes)
	if err := xml.Unmarshal(bodyAsBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal XML: %s", err)
	}
	return nil
}

package client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	config     *Configuration
	httpClient *http.Client
	userAgent  string
	Login      *Login
}

func NewClient(config *Configuration, login *Login) *Client {
	binaryName := os.Getenv("BINARY")
	hostName := os.Getenv("HOSTNAME")
	version := os.Getenv("VERSION")
	namespace := os.Getenv("NAMESPACE")
	userAgent := fmt.Sprintf("%s/%s %s/%s/%s", binaryName, version, hostName, namespace, binaryName)

	return &Client{
		config:     config,
		httpClient: http.DefaultClient,
		userAgent:  userAgent,
		Login:      login,
	}
}

func (client *Client) signAndDo(req *http.Request, bodyAsBytes []byte) (*http.Response, error) {
	var bodyReader io.ReadSeeker
	if bodyAsBytes == nil {
		bodyReader = nil
	} else {
		bodyReader = bytes.NewReader(bodyAsBytes)
	}

	var nilvalue string
	credentials := credentials.NewStaticCredentials(client.config.AccessKeyID, client.config.SecretAccessKey, nilvalue)
	_, err := v4.NewSigner(credentials).Sign(req, bodyReader, client.config.AWSServiceName, client.config.Region, time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %s", err)
	}

	log.Warn("Sending request:\nMethod: %s\nURL: %s\nBody: %s", req.Method, req.URL, bodyAsBytes)
	log.Warn("req--------->", req)
	resp, err := client.httpClient.Do(req)
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

func (client *Client) unmarshalXMLBody(bodyReader io.Reader, v interface{}) error {
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

func (client *Client) unmarshalJSONBody(bodyReader io.Reader, v interface{}) error {
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

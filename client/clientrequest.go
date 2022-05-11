package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"cs-tf-provider/client/utils"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type ClientRequest struct {
	RequestType string
	Url         string
	AuthToken   string
	Body        []byte
	Headers     map[string]string
}

func (cr *ClientRequest) constructRequest(ctx context.Context) (*http.Request, error) {
	var httpReq *http.Request
	var routeToken string
	var msg string
	var err error

	if cr.Body == nil {
		httpReq, err = http.NewRequestWithContext(ctx, cr.RequestType, cr.Url, nil)
	} else {
		httpReq, err = http.NewRequestWithContext(ctx, cr.RequestType, cr.Url, bytes.NewReader(cr.Body))
	}

	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	claims := jwt.MapClaims{}
	_, _, err = new(jwt.Parser).ParseUnverified(cr.AuthToken, claims)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JWT => Error: %s", err)
	}

	if cr.isAdminAPI(httpReq.URL.Path) {
		routeToken = "login"
	} else {
		routeToken = claims["external_id"].(string)
	}

	if cr.Headers == nil {
		cr.Headers = make(map[string]string)
	}

	dateTime := time.Now().UTC().String()
	cr.Headers["Content-Type"] = "*/*"
	cr.Headers["x-amz-chaossumo-route-token"] = routeToken
	cr.Headers["x-amz-date"] = dateTime
	cr.Headers["x-amz-security-token"] = cr.AuthToken
	keys := make([]string, 0, len(cr.Headers))
	for key := range cr.Headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	msg = fmt.Sprintf("%s\n\n", cr.RequestType)
	for _, key := range keys {
		if key == "Content-Type" {
			msg += fmt.Sprintf("%s\n\n", cr.Headers[key])
		} else {
			msg += fmt.Sprintf("%s:%s\n", key, cr.Headers[key])
		}
	}

	msg += httpReq.URL.Path
	signature := fmt.Sprintf(
		"AWS %s:%s",
		claims["AccessKeyId"].(string),
		cr.generateSignature(claims["SecretAccessKey"].(string), msg),
	)

	cr.Headers["Authorization"] = signature
	cr.Headers["x-amz-cs3-authorization"] = signature
	cr.Headers["x-correlation-id"] = uuid.New().String()
	for header, value := range cr.Headers {
		httpReq.Header.Add(header, value)
	}

	return httpReq, nil
}

func (cr *ClientRequest) generateSignature(secretToken string, payloadBody string) string {
	keyForSign := []byte(secretToken)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(payloadBody))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (cr *ClientRequest) isAdminAPI(url string) bool {
	return strings.HasSuffix(url, "/createSubAccount") ||
		strings.HasSuffix(url, "/deleteSubAccount") ||
		strings.HasSuffix(url, "/user/manifest") ||
		strings.HasSuffix(url, "/user/groups") ||
		strings.Contains(url, "/user/group/")
}

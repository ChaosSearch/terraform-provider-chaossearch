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

func (cr *ClientRequest) constructRequest(ctx context.Context, config Configuration) (*http.Request, error) {
	var routeToken string
	httpReq, err := cr.newRequest(ctx)
	if err != nil {
		return nil, utils.CreateRequestError(err)
	}

	if config.KeyAuthEnabled {
		return cr.constructHeaders(
			config.AccessKeyID,
			config.SecretAccessKey,
			config.UserID,
			httpReq,
		), nil
	} else {
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

		return cr.constructHeaders(
			claims["AccessKeyId"].(string),
			claims["SecretAccessKey"].(string),
			routeToken,
			httpReq,
		), nil
	}
}

func (cr *ClientRequest) newRequest(ctx context.Context) (*http.Request, error) {
	if cr.Body == nil {
		return http.NewRequestWithContext(ctx, cr.RequestType, cr.Url, nil)
	} else {
		return http.NewRequestWithContext(ctx, cr.RequestType, cr.Url, bytes.NewReader(cr.Body))
	}
}

func (cr *ClientRequest) constructHeaders(
	accessKey,
	secretAccessKey,
	routeToken string,
	httpReq *http.Request,
) *http.Request {
	if cr.Headers == nil {
		cr.Headers = make(map[string]string)
	}

	if cr.AuthToken != "" {
		cr.Headers["x-amz-security-token"] = cr.AuthToken
	}

	dateTime := time.Now().UTC().String()
	cr.Headers["x-amz-date"] = dateTime
	cr.Headers["Content-Type"] = "*/*"
	cr.Headers["x-amz-chaossumo-route-token"] = routeToken
	keys := make([]string, 0, len(cr.Headers))
	for key := range cr.Headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	msg := fmt.Sprintf("%s\n\n", cr.RequestType)
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
		accessKey,
		cr.generateSignature(secretAccessKey, msg),
	)

	cr.Headers["Authorization"] = signature
	cr.Headers["x-amz-cs3-authorization"] = signature
	cr.Headers["x-correlation-id"] = uuid.New().String()
	for header, value := range cr.Headers {
		httpReq.Header.Add(header, value)
	}

	return httpReq
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

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
Testing class for API end point
need to remove this by end of initial developments
*/

func createObjectGroup() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/Bucket/createObjectGroup"
	return false, "POST", url, strings.NewReader(`{
    "bucket": "dinesh-og-121",
    "source": "chaos-test-data-aps1",
    "format": {
        "_type": "CSV",
        "columnDelimiter": ",",
        "rowDelimiter": "\n",
        "headerRow": true
    },
    "interval": {
        "mode": 0,
        "column": 0
    },
    "indexRetention": {
        "overall": -1,
        "forPartition": []
    },
    "filter": [
        {
            "field": "key",
            "prefix": "bluebike"
        },
        {
            "field": "key",
            "regex": ".*"
        }
    ],
    "options": {
        "ignoreIrregular": true
    },
    "realtime": false
}`)

}

func putObjectGroup() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/Bucket/createObjectGroup"
	return false, "PUT", url, strings.NewReader(`{
    "bucket": "dinesh-og-2143",
    "source": "chaos-test-data-aps1",
    "format": {
        "_type": "CSV",
        "columnDelimiter": ",",
        "rowDelimiter": "\n",
        "headerRow": true
    },
    "interval": {
        "mode": 0,
        "column": 0
    },
    "indexRetention": {
        "overall": -1,
        "forPartition": []
    },
    "filter": [
        {
            "field": "key",
            "prefix": "bluebike"
        },
        {
            "field": "key",
            "regex": ".*"
        }
    ],
    "options": {
        "ignoreIrregular": true
    },
    "realtime": false
}`)

}

func createView() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/Bucket/createView"
	return false, "POST", url, strings.NewReader(`{
	   "bucket": "dinesh-view-name121",
	   "sources": [],
	   "indexPattern": ".*",
	   "overwrite": true,
	   "caseInsensitive": true,
	   "indexRetention": -1,
	   "timeFieldName": "@timestamp",
	   "transforms": [],
	   "filter": {
	           "predicate": {
	                   "field": "cs_partition_key_0",
	                   "query": "bluebike",
	                   "state": {
	                       "_type": "chaossumo.query.QEP.Predicate.TextMatchState.Exact"
	                   },
	                   "_type": "chaossumo.query.NIRFrontend.Request.Predicate.TextMatch"
	               }
	       }
	   }
	`)

}
func createSubAccount() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/createSubAccount"
	return true, "POST", url, strings.NewReader(`{
	    "UserInfoBlock": {
	        "Username": "dineshkj",
	        "FullName": "dinesh k j",
	        "Email": "dineshkj@gmail.com"
	    },
	    "GroupIds": [
	        "default"
	    ],
	    "Password": "dineshkj",
	    "Hocon": [
	        "override.Services.worker.quota=50"
	    ]
	}`)
}

func deleteSubAccountByUser() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/deleteSubAccount"
	return true, "POST", url, strings.NewReader(`{"Username": "nibras_41"}`)
}

func importBucket() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/Bucket/importBucket"
	return false, "POST", url, strings.NewReader(`{
		"bucket":"chaos-tera-test-111",
		"hideBucket":false
}`)
}
func retrieveUserGroups() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/groups"
	return true, "GET", url, nil
}
func retrieveObjectGroups() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/V1/"
	return false, "GET", url, nil
}

func retrieveUserGroupByGroupId() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/group/9436aed9-e994-4dba-a25b-7d950d7f3623"
	return true, "GET", url, nil
}

func retrieveUsers() (control bool, method string, url string, reader io.Reader) {
	url = "https://aeternum.chaossearch.io/user/manifest"
	return true, "POST", url, nil
}

func deleteUserGroupByGroupId() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/group/f4da25ba-5e8e-4cba-b13a-c068bdfa07d7"
	return true, "DELETE", url, nil
}

func createUserGroup() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/groups"
	return true, "POST", url, strings.NewReader(`[
    {
        "id": "100042",
        "name": "dinesh-jayddasinsdghe",
        "permissions": [
            {
                "Effect": "Allow",
                "Version": "1.2",
                "Actions": [
                    "*"
                ],
                "Resources": [
                    "*"
                ]
            }
        ]
    }

]`)

}

func putUserGroup() (control bool, method string, url string, reader io.Reader) {
	url = "https://ap-south-1-aeternum.chaossearch.io/user/groups"
	return true, "PUT", url, strings.NewReader(`[{"id":"12bdf713-f019-42c0-9185-fb1b16b001cf"
,"name":"dinesh-view-1","permissions":[{"Version":"2.0","Effect":"Allow","Actions":["*"],"Resources":["*"],"Condition":null}]}]`)

}

func main() {
	//control, method, url, payload := createView()
	//control, method, url, payload := createObjectGroup()

	//control, method, url, payload := createSubAccount()
	//control, method, url, payload := deleteSubAccountByUser()

	//control, method, url, payload := createUserGroup()
	//control, method, url, payload := putUserGroup()
	//control, method, url, payload := retrieveUserGroups()
	//control, method, url, payload := retrieveUserGroupByGroupId()
	//control, method, url, payload := retrieveUserGroupByGroupId()
	//control, method, url, payload := deleteSubAccountByUser()
	//control, method, url, payload := deleteUserGroupByGroupId()
	//control, method, url, payload := importBucket()
	//control, method, url, payload := retrieveUsers()
	control, method, url, payload := retrieveObjectGroups()

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	authToken, err := Auth()
	if err != nil {
		return
	}
	fmt.Println("url-->", url)
	fmt.Println("payload-->", payload)

	var bodyBytes []byte
	if payload != nil {
		if method == "POST" || method == "PUT" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(payload)
			bodyBytes = buf.Bytes()
		} else {
			bodyBytes = nil
		}
	} else {
		bodyBytes = nil
	}

	httpResp, err := signV2AndDo(control, authToken, req, bodyBytes)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func signV2AndDo(control bool, tokenValue string, req *http.Request, bodyAsBytes []byte) (*http.Response, error) {
	fmt.Println("------- AWS V2 Sign Starts------")
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}

	accessKey := claims["AccessKeyId"].(string)
	secretAccessKey := claims["SecretAccessKey"].(string)
	externalId := claims["external_id"].(string)
	dateTime := time.Now().UTC().String()
	fmt.Println("externalId-->", externalId)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-amz-security-token", tokenValue)

	var routeToken string

	if control {
		routeToken = "login"
	} else {
		routeToken = externalId
	}
	req.Header.Add("x-amz-chaossumo-route-token", routeToken)

	req.Header.Add("X-Amz-Date", dateTime)

	msgLines := []string{
		req.Method, "",
		"application/json", "",
		"x-amz-chaossumo-route-token:" + routeToken,
		"x-amz-date:" + dateTime,
		"x-amz-security-token:" + tokenValue,
		req.URL.Path,
	}

	msg := strings.Join(msgLines, "\n")
	signature := generateSignature(secretAccessKey, msg)
	auth := "AWS " + accessKey + ":" + signature

	req.Header.Add("Authorization", auth)
	req.Header.Add("x-amz-cs3-authorization", auth)

	for key, val := range req.Header {
		fmt.Println("Header -->", key, "  value -->", val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}

	if resp.Body != nil {
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			return nil, fmt.Errorf("failed to read response body--->: %s", err1)
		}
		fmt.Println("Response Body -->", string(body))
	}

	fmt.Println("Got response:\nStatus code: %d", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respAsBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}
		return nil, fmt.Errorf(
			"expected a 2xx status code, but got %d.\nMethod: %s\nURL: %s\nRequest body: %s\nResponse body: %s",
			resp.StatusCode, req.Method, req.URL, bodyAsBytes, respAsBytes)
	}
	fmt.Println("------- AWS V2 Sign Ends------")
	return resp, nil
}

const loginUrl = "https://ap-south-1-aeternum.chaossearch.io/user/login"

func Auth() (token string, err error) {

	method := "POST"

	login := Login{
		Username: "",
		Password: "",
	}

	bodyAsBytes, err := marshalLoginRequest(&login)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, loginUrl, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %s", err)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-amz-chaossumo-route-token", "login")
	req.Header.Add("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Response jwt token-->", string(body))

	if err != nil {
		diag.Errorf("Token generation fail..")
		return "", nil
	} else {
		tokenData := AuthResponse{}
		if err := json.Unmarshal([]byte(string(body)), &tokenData); err != nil {
			fmt.Errorf("failed to unmarshal JSON: %s", err)
		}
		return tokenData.Token, nil
	}

}

func marshalLoginRequest(req *Login) ([]byte, error) {

	body := map[string]interface{}{
		"Username": req.Username,
		"Password": req.Password,
	}

	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyAsBytes, nil
}

func generateSignature(secretToken string, payloadBody string) string {
	keyForSign := []byte(secretToken)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(payloadBody))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type Login struct {
	Username     string
	Password     string
	ParentUserId string
}
type AuthResponse struct {
	Token string
}

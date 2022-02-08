package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Auth() (token string, err error) {

	url := "https://ap-south-1-aeternum.chaossearch.io/user/login"
	method := "POST"

	payload := strings.NewReader(`{
	  "Username": "service_user@chaossearch.com",
	  "Password": "thisIsAnEx@mple1!",
	  "ParentUid": "be4aeb53-21d5-4902-862c-9c9a17ad6675"
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-amz-chaossumo-route-token", "login")
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	return string(body),nil
}

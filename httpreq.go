package goeureka

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	// "strconv"
)

// Accepts a Httpaction and a one-way channel to write the results to.
func DoHttpRequest(httpAction HttpAction) bool {
	http.DefaultClient.Timeout = 10 * time.Second
	fmt.Println("Begin Request")
	req := buildHttpRequest(httpAction)
	fmt.Println("Request built")
	var DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	fmt.Println("Transport compleetet")
	//resp, err := DefaultTransport.RoundTrip(req)
	client := &http.Client{Transport: DefaultTransport}
	fmt.Println("Client set up")
	resp, err := client.Do(req)
	fmt.Println("request compelete")
	fmt.Println(resp)
	if resp != nil {
		defer resp.Body.Close()
		body := getBody(resp)
		if err != nil {
			log.Printf("HTTP request failed: %s", err)
			log.Println("Response body: ", body)
			log.Println("Response: ", resp.StatusCode)
			return false
		} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
			return true
		}
	}
	return false
}

func getBody(resp *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("Unable to read response body")
	}
	return string(bodyBytes)
}

func buildHttpRequest(httpAction HttpAction) *http.Request {
	var req *http.Request
	var err error
	if httpAction.Body != "" {
		reader := strings.NewReader(httpAction.Body)
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, reader)
	} else if httpAction.Template != "" {
		reader := strings.NewReader(httpAction.Template)
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, reader)
	} else {
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Add headers
	req.Header.Add("Accept", httpAction.Accept)
	if httpAction.ContentType != "" {
		req.Header.Add("Content-Type", httpAction.ContentType)
	}
	return req
}

/**
 * Trims leading and trailing byte r from string s
 */
func trimChar(s string, r byte) string {
	sz := len(s)

	if sz > 0 && s[sz-1] == r {
		s = s[:sz-1]
	}
	sz = len(s)
	if sz > 0 && s[0] == r {
		s = s[1:sz]
	}
	return s
}

package goeureka

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	// "strconv"
)

// Accepts a Httpaction and a one-way channel to write the results to.
func DoHttpRequest(httpAction HttpAction) bool {
	req := buildHttpRequest(httpAction)
	var DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: DefaultTransport, Timeout: time.Duration(10 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("HTTP request failed: %s", err)
		if resp != nil {
			logger.Error("Response: ", resp.StatusCode)
			logger.Error(resp)
		}
		return false
	}
	logger.Info("Eureka Response:", resp.StatusCode)
	if resp != nil {
		defer resp.Body.Close()
		body, err := getBody(resp)
		if err != nil {
			logger.Printf("Error reading response body : %s", err)
			logger.Error("Response body: ", body)
			logger.Error("Response: ", resp.StatusCode)
			return false
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
			return true
		} else {
			logger.Error("Response body: ", body)
			logger.Error("Response: ", resp.StatusCode)
			return false
		}
	}
	return false
}

func getBody(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Unable to read response body")
		return "", err
	}
	return string(bodyBytes), nil
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
		logger.Error(err)
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

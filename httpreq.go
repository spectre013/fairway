package fairway

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	// "strconv"
)

// DoHTTPRequest Accepts a Httpaction and a one-way channel to write the results to.
func DoRegRequest(httpAction HTTPAction) bool {
	req := buildHTTPRequest(httpAction)
	resp, err := doRequest(req)
	if err != nil {
		logger.Error("HTTP request failed:", err)
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
		}
		logger.Error("Response body: ", body)
		logger.Error("Response: ", resp.StatusCode)
		return false

	}
	return false
}

func doRequest(req *http.Request) (*http.Response, error) {
	var DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: DefaultTransport, Timeout: time.Duration(10 * time.Second)}
	return  client.Do(req)
}

func getBody(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Unable to read response body")
		return "", err
	}
	return string(bodyBytes), nil
}

func buildHTTPRequest(httpAction HTTPAction) *http.Request {
	var req *http.Request
	var err error
	if httpAction.Body != "" {
		reader := strings.NewReader(httpAction.Body)
		req, err = http.NewRequest(httpAction.Method, httpAction.URL, reader)
	} else if httpAction.Template != "" {
		reader := strings.NewReader(httpAction.Template)
		req, err = http.NewRequest(httpAction.Method, httpAction.URL, reader)
	} else {
		req, err = http.NewRequest(httpAction.Method, httpAction.URL, nil)
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

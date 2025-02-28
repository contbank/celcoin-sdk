package celcoin

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggingHTTPClient struct {
	client *http.Client
}

func NewLoggingHTTPClient(client *http.Client) *LoggingHTTPClient {
	return &LoggingHTTPClient{client: client}
}

func (c *LoggingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	var reqBody []byte
	// Log request details
	if req != nil && req.Body != nil {
		reqBody, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Restore the body for further use
	}

	logrus.WithFields(logrus.Fields{
		"celcoin_request": logrus.Fields{
			"method":     req.Method,
			"url":        req.URL.String(),
			"header":     req.Header,
			"body":       string(reqBody),
			"user-agent": req.UserAgent(),
		},
	}).Info("HTTP Request Celcoin")

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("HTTP request failed")
		return nil, err
	}

	// Log response details
	duration := time.Since(start)

	var respBody []byte
	if resp != nil && resp.Body != nil {
		respBody, _ = ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody)) // Restore the body for further use
	}

	logrus.WithFields(logrus.Fields{
		"celcoin_response": logrus.Fields{
			"header":   resp.Header,
			"status":   resp.StatusCode,
			"duration": duration,
			"body":     string(respBody),
		},
	}).Info("HTTP Response Celcoin")

	return resp, nil
}

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

	// Log request details
	logrus.WithFields(logrus.Fields{
		"method": req.Method,
		"url":    req.URL.String(),
		"header": req.Header,
		"body":   req.Body,
	}).Info("HTTP request")

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("HTTP request failed")
		return nil, err
	}

	// Log response details
	duration := time.Since(start)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Restore the body for further use

	logrus.WithFields(logrus.Fields{
		"header":   resp.Header,
		"status":   resp.StatusCode,
		"duration": duration,
		"body":     string(body),
	}).Info("HTTP response")

	return resp, nil
}

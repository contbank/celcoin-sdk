package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Dda ...
type Dda struct {
	session        Session
	authentication *Authentication
	httpClient     *LoggingHTTPClient
}

// NewDda ...
func NewDda(httpClient *http.Client, session Session) *Dda {
	return &Dda{
		session:        session,
		httpClient:     NewLoggingHTTPClient(httpClient),
		authentication: NewAuthentication(httpClient, session),
	}
}

// CreateRegisterUser ...
func (s *Dda) CreateRegisterUser(ctx context.Context, correlationID string,
	model DdaRegisterUserRequest) (*DdaRegisterUserResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create register user")
	return s.createRegisterUser(ctx, correlationID, model)
}

// createRegisterUser ...
func (s *Dda) createRegisterUser(ctx context.Context, requestID string,
	model DdaRegisterUserRequest) (*DdaRegisterUserResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
	}
	logrus.WithFields(fields).Info("Create Register User")
	req := &model

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(DdaSubscriptionPath, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for CreateDdaSubscription")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling CreateDdaSubscription")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *DdaRegisterUserResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultDda
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErroDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultDda
	}

	if errResponse.Error != nil {
		err := FindDdaError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultDda

}

// DeleteRegisterUser ...
func (s *Dda) DeleteRegisterUser(ctx context.Context, correlationID string,
	model DdaDeleteUserRequest) (*DdaRegisterUserResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).Info("delete register user")
	return s.deleteRegisterUser(ctx, correlationID, model)
}

// deleteRegisterUser ...
func (s *Dda) deleteRegisterUser(ctx context.Context, requestID string,
	model DdaDeleteUserRequest) (*DdaRegisterUserResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
	}
	logrus.WithFields(fields).Info("Create Register User")
	req := &model

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(DdaSubscriptionPath, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for CreateDdaSubscription")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling CreateDdaSubscription")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *DdaRegisterUserResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultDda
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErroDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultDda
	}

	if errResponse.Error != nil {
		err := FindDdaError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultDda
}

func (s *Dda) BuildEndpoint(basePath string, queryParams map[string]string, pathParams ...string) (*string, error) {
	u, err := url.Parse(s.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).Error("Error parsing API endpoint")
		return nil, err
	}

	// Construção do caminho completo
	fullPath := path.Join(basePath, path.Join(pathParams...))
	u.Path = path.Join(u.Path, fullPath)

	// Adição de query parameters
	if queryParams != nil {
		q := u.Query()
		for key, value := range queryParams {
			if value != "" {
				q.Set(key, value)
			}
		}
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()
	logrus.WithField("endpoint", endpoint).Info("Endpoint built successfully")
	return &endpoint, nil
}

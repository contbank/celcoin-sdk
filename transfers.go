package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// Transfers ...
type Transfers struct {
	session        Session
	authentication *Authentication
	httpClient     *http.Client
}

// NewTransfers ...
func NewTransfers(httpClient *http.Client, session Session) *Transfers {
	return &Transfers{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// CreateTransfer ...
func (t *Transfers) CreateTransfer(ctx context.Context, correlationID string, model TransfersRequest) (*TransfersResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create transfer")
	return t.createTransferOperation(ctx, correlationID, model)
}

// createTransferOperation ...
func (t *Transfers) createTransferOperation(ctx context.Context, requestID string,
	model TransfersRequest) (*TransfersResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}
	model.ClientCode = requestID

	// Description is required only when client finality is Others
	if strings.Compare(string(model.ClientFinality), string(OthersClientFinality)) == 0 && len(model.Description) <= 0 {
		return nil, grok.NewError(http.StatusBadRequest, "DESCRIPTION_MISSING", "description is required")
	}

	var isInternalTransfer bool = false
	if model.CreditParty.BankISPB == model.CreditParty.BankISPB {
		isInternalTransfer = true
	}

	endpoint, err := t.getTransfersAPIEndpoint(requestID, nil, nil, isInternalTransfer)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error getting api endpoint")
		return nil, err
	}

	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var body *TransfersResponse

	err = json.Unmarshal(respBody, &body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return body, nil
	}

	// error
	if body.Error != nil {
		logrus.WithFields(fields).Error("body error - createTransferOperation")
		var httpErrorStatus int
		httpErrorStatus, err = strconv.Atoi(body.Status)
		if err != nil {
			httpErrorStatus = http.StatusBadRequest
		}
		return nil, grok.NewError(httpErrorStatus, "TRANSFERS_ERROR_"+body.Error.ErrorCode,
			body.Error.ErrorCode+" - "+body.Error.Message)
	}

	logrus.WithFields(fields).
		Error("default error transfer - createTransferOperation")

	return nil, ErrDefaultTransfers
}

// FindTransferByCode ...
func (t *Transfers) FindTransferByCode(ctx context.Context, requestID *string,
	transferAuthenticationCode string, transferRequestID string, isInternalTransfer *bool) (*TransfersResponse, error) {

	if requestID == nil {
		return nil, ErrInvalidCorrelationID
	} else if len(transferAuthenticationCode) == 0 || len(transferRequestID) == 0 {
		return nil, ErrInvalidTransferAuthenticationCode
	}

	fields := logrus.Fields{
		"request_id":                   requestID,
		"transfer_authentication_code": transferAuthenticationCode,
		"transfer_request_id":          transferRequestID,
	}

	if isInternalTransfer != nil && *isInternalTransfer == true {
		logrus.WithFields(fields).Info("find transfer by code - internal flag")
		return t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, true)
	} else if isInternalTransfer != nil && *isInternalTransfer == false {
		logrus.WithFields(fields).Info("find transfer by code - external flag")
		return t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, false)
	} else {
		logrus.WithFields(fields).Info("find transfer by code - internal")
		var resp *TransfersResponse = nil
		resp, err := t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, true)
		if err != nil {
			logrus.WithFields(fields).Info("find transfer by code - external")
			resp, err = t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, false)
			if err != nil {
				logrus.WithFields(fields).WithError(err).
					Error("default error transfer - FindTransferByCode")
				return nil, ErrDefaultFindTransfers
			}
		}
		return resp, nil
	}

}

// findInternalOrExternalTransferByCode ...
func (t *Transfers) findInternalOrExternalTransferByCode(ctx context.Context, requestID *string,
	transferAuthenticationCode string, transferRequestID string, isInternalTransfer bool) (*TransfersResponse, error) {

	if requestID == nil {
		return nil, ErrInvalidCorrelationID
	} else if len(transferAuthenticationCode) == 0 || len(transferRequestID) == 0 {
		return nil, ErrInvalidTransferAuthenticationCode
	}

	fields := logrus.Fields{
		"request_id":                   requestID,
		"transfer_authentication_code": transferAuthenticationCode,
		"transfer_request_id":          transferRequestID,
		"is_internal_transfer":         isInternalTransfer,
	}

	endpoint, err := t.getTransfersAPIEndpoint(*requestID, &transferAuthenticationCode, &transferRequestID, isInternalTransfer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := t.authentication.Token(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", *requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	// response not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var body *TransfersResponse

	err = json.Unmarshal(respBody, &body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK {
		return body, nil
	}

	// error
	if body.Error != nil {
		logrus.WithFields(fields).Error("body error - FindTransferByCode")
		var httpErrorStatus int
		httpErrorStatus, err = strconv.Atoi(body.Status)
		if err != nil {
			httpErrorStatus = http.StatusBadRequest
		}
		return nil, grok.NewError(httpErrorStatus, "TRANSFERS_ERROR_"+body.Error.ErrorCode,
			body.Error.ErrorCode+" - "+body.Error.Message)
	}

	logrus.WithFields(fields).
		Error("default error transfer - findInternalOrExternalTransferByCode")

	return nil, ErrDefaultFindTransfers
}

// getTransfersAPIEndpoint
func (t *Transfers) getTransfersAPIEndpoint(correlationID string,
	transferAuthenticationCode *string, transferRequestID *string, isInternalTransfer bool) (*string, error) {

	u, err := url.Parse(t.session.APIEndpoint)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"correlation_id":               correlationID,
				"transfer_authentication_code": transferAuthenticationCode,
				"transfer_request_id":          transferRequestID,
			}).
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}

	if isInternalTransfer {
		u.Path = path.Join(u.Path, InternalTransfersPath)
	} else {
		u.Path = path.Join(u.Path, ExternalTransfersPath)
	}

	if transferAuthenticationCode != nil && transferRequestID != nil {
		u.Path = path.Join(u.Path, "status")
		q := u.Query()
		q.Set("id", *transferAuthenticationCode)
		q.Set("clientCode", *transferRequestID)
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()
	return &endpoint, nil
}

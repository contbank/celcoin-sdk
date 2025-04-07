package celcoin

import (
	"context"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Cards ...
type Cards struct {
	session        Session
	authentication *Authentication
	httpClient     *http.Client
}

// NewCards ...
func NewCards(httpClient *http.Client, session Session) *Cards {
	return &Cards{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// CreateCard ...
func (t *Cards) CreateCard(ctx context.Context, requestID string,
	model CreateCardRequest) (*CreateCardResponse, error) {

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
	/*
		AAAAA //// parei aqui

		// Description is required only when client finality is Others
		if strings.Compare(string(model.ClientFinality), string(OthersClientFinality)) == 0 && len(model.Description) <= 0 {
			return nil, grok.NewError(http.StatusBadRequest, "DESCRIPTION_MISSING", "description is required")
		}

		var isInternalTransfer bool = false
		if grok.OnlyDigits(model.DebitParty.BankISPB) == grok.OnlyDigits(model.CreditParty.BankISPB) {
			isInternalTransfer = true
			model.ClientRequestId = model.ClientCode
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
			if isInternalTransfer && len(body.Body.ClientRequestId) > 0 && len(body.Body.ClientCode) == 0 {
				body.Body.ClientCode = body.Body.ClientRequestId
			}
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
	*/
	return nil, ErrDefaultTransfers
}

//// FindTransferByCode ...
//func (t *Transfers) FindTransferByCode(ctx context.Context, requestID *string,
//	transferAuthenticationCode string, transferRequestID string, isInternalTransfer *bool) (*TransfersResponse, error) {
//
//	if requestID == nil {
//		return nil, ErrInvalidCorrelationID
//	} else if len(transferAuthenticationCode) == 0 || len(transferRequestID) == 0 {
//		return nil, ErrInvalidTransferAuthenticationCode
//	}
//
//	fields := logrus.Fields{
//		"request_id":                   requestID,
//		"transfer_authentication_code": transferAuthenticationCode,
//		"transfer_request_id":          transferRequestID,
//	}
//
//	if isInternalTransfer != nil && *isInternalTransfer == true {
//		logrus.WithFields(fields).Info("find transfer by code - internal flag")
//		return t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, true)
//	} else if isInternalTransfer != nil && *isInternalTransfer == false {
//		logrus.WithFields(fields).Info("find transfer by code - external flag")
//		return t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, false)
//	} else {
//		logrus.WithFields(fields).Info("find transfer by code - internal")
//		var resp *TransfersResponse = nil
//		resp, err := t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, true)
//		if err != nil || (resp != nil && resp.Error != nil) {
//			logrus.WithFields(fields).Info("find transfer by code - external")
//			resp, err = t.findInternalOrExternalTransferByCode(ctx, requestID, transferAuthenticationCode, transferRequestID, false)
//			if err != nil || (resp != nil && resp.Error != nil) {
//				logrus.WithFields(fields).WithError(err).
//					Error("default error transfer - FindTransferByCode")
//				return nil, ErrDefaultFindTransfers
//			}
//		}
//		return resp, nil
//	}
//
//}
//
//// findInternalOrExternalTransferByCode ...
//func (t *Transfers) findInternalOrExternalTransferByCode(ctx context.Context, requestID *string,
//	transferAuthenticationCode string, transferRequestID string, isInternalTransfer bool) (*TransfersResponse, error) {
//
//	if requestID == nil {
//		return nil, ErrInvalidCorrelationID
//	} else if len(transferAuthenticationCode) == 0 || len(transferRequestID) == 0 {
//		return nil, ErrInvalidTransferAuthenticationCode
//	}
//
//	fields := logrus.Fields{
//		"request_id":                   requestID,
//		"transfer_authentication_code": transferAuthenticationCode,
//		"transfer_request_id":          transferRequestID,
//		"is_internal_transfer":         isInternalTransfer,
//	}
//
//	endpoint, err := t.getTransfersAPIEndpoint(*requestID, &transferAuthenticationCode, &transferRequestID, isInternalTransfer)
//	if err != nil {
//		return nil, err
//	}
//
//	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	token, err := t.authentication.Token(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	req.Header.Add("Authorization", token)
//	req.Header.Add("Content-type", "application/json")
//	req.Header.Add("api-version", t.session.APIVersion)
//	req.Header.Add("x-correlation-id", *requestID)
//
//	resp, err := t.httpClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//
//	respBody, _ := ioutil.ReadAll(resp.Body)
//
//	// response not found
//	if resp.StatusCode == http.StatusNotFound {
//		return nil, ErrEntryNotFound
//	}
//
//	var body *TransfersResponse
//
//	err = json.Unmarshal(respBody, &body)
//	if err != nil {
//		logrus.WithFields(fields).WithError(err).
//			Error("error unmarshal")
//		return nil, err
//	}
//
//	// response ok
//	if resp.StatusCode == http.StatusOK {
//		return body, nil
//	}
//
//	// error
//	if body.Error != nil {
//		logrus.WithFields(fields).Error("body error - FindTransferByCode")
//		var httpErrorStatus int
//		httpErrorStatus, err = strconv.Atoi(body.Status)
//		if err != nil {
//			httpErrorStatus = http.StatusBadRequest
//		}
//		return nil, grok.NewError(httpErrorStatus, "TRANSFERS_ERROR_"+body.Error.ErrorCode,
//			body.Error.ErrorCode+" - "+body.Error.Message)
//	}
//
//	logrus.WithFields(fields).
//		Error("default error transfer - findInternalOrExternalTransferByCode")
//
//	return nil, ErrDefaultFindTransfers
//}

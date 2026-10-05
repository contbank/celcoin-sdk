package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// ReversePixCashIn devolve (total ou parcialmente) um Pix recebido. A resposta vem PROCESSING;
// o resultado final chega pelo webhook pix-reversal-out. Prazo: até 90 dias do recebimento.
func (s *Pix) ReversePixCashIn(ctx context.Context, req PixReversalRequest) (*PixReversalResponse, error) {
	fields := logrus.Fields{
		"id":          req.ID,
		"end_to_end":  req.EndToEndID,
		"client_code": req.ClientCode,
		"amount":      req.Amount,
		"reason":      req.Reason,
	}
	logrus.WithFields(fields).Info("Reverse Pix CashIn")

	if err := grok.Validator.Struct(req); err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixReversalPath, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for ReversePixCashIn")
		return nil, err
	}

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
	httpReq.Header.Set("Accept", "application/json")

	return s.doPixReversal(httpReq, fields)
}

// GetPixReversalStatus consulta uma devolução por id, clientCode ou returnIdentification (E2E da devolução).
func (s *Pix) GetPixReversalStatus(ctx context.Context, id, clientCode, returnIdentification string) (*PixReversalResponse, error) {
	fields := logrus.Fields{"id": id, "client_code": clientCode, "return_identification": returnIdentification}
	logrus.WithFields(fields).Info("Get Pix Reversal Status")

	if id == "" && clientCode == "" && returnIdentification == "" {
		return nil, ErrPixReversalStatusParameters
	}

	params := map[string]string{
		"id":                   id,
		"clientCode":           clientCode,
		"returnIdentification": returnIdentification,
	}
	endpoint, err := s.BuildEndpoint(PixReversalPath, params, "status")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for GetPixReversalStatus")
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Accept", "application/json")

	return s.doPixReversal(httpReq, fields)
}

func (s *Pix) doPixReversal(httpReq *http.Request, fields logrus.Fields) (*PixReversalResponse, error) {
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixReversalResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil || errResponse == nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil && errResponse.Error.ErrorCode != nil {
		err := FindPixErrorWithMessage(*errResponse.Error.ErrorCode, &resp.StatusCode, errResponse.Error.Message)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin pix reversal error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// ErrPixReversalStatusParameters consulta de devolução sem nenhum identificador.
var ErrPixReversalStatusParameters = grok.NewError(http.StatusUnprocessableEntity, "MISSING_REVERSAL_IDENTIFIER",
	"Informe id, clientCode ou returnIdentification da devolução.")

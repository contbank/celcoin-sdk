package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
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
		"method":     "CreateCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin create card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// post request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *CreateCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("body error")
		var httpErrorStatus int
		httpErrorStatus, err = strconv.Atoi(responseBody.Status)
		if err != nil {
			httpErrorStatus = http.StatusBadRequest
		}
		return nil, grok.NewError(httpErrorStatus, fmt.Sprint("CREATE_CARDS_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// GetCard ...
func (t *Cards) GetCard(ctx context.Context, requestID string,
	model GetCardRequest) (*GetCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "GetCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card"))

	// params
	q := u.Query()
	if model.CardID != nil {
		q.Set("cardId", *model.CardID)
	}
	if model.Identifier != nil {
		q.Set("document", grok.OnlyDigits(*model.Identifier))
	}
	if model.Status != nil {
		q.Set("status", string(*model.Status))
	}
	if model.Type != nil {
		q.Set("type", string(*model.Type))
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin get card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// GET request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *GetCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("GET_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// ListCards ...
func (t *Cards) ListCards(ctx context.Context, requestID string,
	model ListCardsRequest) (*ListCardsResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "ListCards",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/cards"))

	// params
	q := u.Query()
	q.Set("page", string(model.Page))
	q.Set("perPage", string(model.PerPage))

	if model.Status != nil {
		q.Set("status", string(*model.Status))
	}
	if model.CardModel != nil {
		q.Set("modes", string(*model.CardModel))
	}
	if model.CardType != nil {
		q.Set("type", string(*model.CardType))
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin list cards request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// GET request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ListCardsResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("LIST_CARDS_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// UpdateCard
func (t *Cards) UpdateCard(ctx context.Context, requestID string,
	model UpdateCardRequest) (*UpdateCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "UpdateCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin update card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PATCH request
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *UpdateCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(http.StatusBadRequest, fmt.Sprint("UPDATE_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// ActivateCard ...
func (t *Cards) ActivateCard(ctx context.Context, requestID string,
	model ActivateCardRequest) (*ActivateCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "ActivateCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/activate"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin activate card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PATCH request
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ActivateCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("ACTIVATE_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// UpdateCardStatus
func (t *Cards) UpdateCardStatus(ctx context.Context, requestID string,
	model UpdateCardStatusRequest) (*UpdateCardStatusResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "UpdateCardStatus",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/status"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin update card status request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PUT request
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *UpdateCardStatusResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("UPDATE_CARD_STATUS_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// BlockCard
func (t *Cards) BlockCard(ctx context.Context, requestID string,
	model BlockCardRequest) (*BlockCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "BlockCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/block"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin block card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PUT request
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *BlockCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("BLOCK_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// UnblockCard
func (t *Cards) UnblockCard(ctx context.Context, requestID string,
	model UnblockCardRequest) (*UnblockCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "UnblockCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/unblock"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin unblock card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PUT request
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *UnblockCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("UNBLOCK_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// CancelCard
func (t *Cards) CancelCard(ctx context.Context, requestID string,
	model CancelCardRequest) (*CancelCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "CancelCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/cancel"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin cancel card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PUT request
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *CancelCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("CANCEL_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// ReissueCard
func (t *Cards) ReissueCard(ctx context.Context, requestID string,
	model ReissueCardRequest) (*ReissueCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "ReissueCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/reissue"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin reissue card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// POST request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ReissueCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("REISSUE_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// ChangePasswordCard
func (t *Cards) ChangePasswordCard(ctx context.Context, requestID string,
	model ChangeCardPasswordRequest) (*ChangeCardPasswordResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "ChangePasswordCard",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/changeCardPassword"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin change password card request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// POST request
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ChangeCardPasswordResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("CHANGE_PASSWORD_CARD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// Embossing - GetCardEmbossing
func (t *Cards) GetCardEmbossing(ctx context.Context, requestID string,
	model GetCardEmbossingRequest) (*GetCardEmbossingResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "GetCardEmbossing",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/embossing"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin get card embossing request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// GET request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *GetCardEmbossingResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(responseBody.Status, fmt.Sprint("GET_CARD_EMBOSSING_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// Embossing - UpdateCardEmbossingAddress
func (t *Cards) UpdateCardEmbossingAddress(ctx context.Context, requestID string,
	model UpdateCardEmbossingAddressRequest) (*UpdateCardEmbossingAddressResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "UpdateCardEmbossingAddress",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/embossing/address"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin update card embossing address request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PATCH request
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *UpdateCardEmbossingAddressResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(http.StatusInternalServerError, fmt.Sprint("UPDATE_CARD_EMBOSSING_ADDRESS_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// Embossing - ReissueCardEmbossing
func (t *Cards) ReissueCardEmbossing(ctx context.Context, requestID string,
	model ReissueCardEmbossingRequest) (*ReissueCardEmbossingResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "ReissueCardEmbossing",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/embossing/reembossing"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin reissue card embossing request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// POST request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ReissueCardEmbossingResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(http.StatusInternalServerError, fmt.Sprint("REISSUE_CARD_EMBOSSING_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// Embossing - GetCurrentCardPassword
func (t *Cards) GetCurrentCardPassword(ctx context.Context, requestID string,
	model ViewCardPasswordRequest) (*ViewCardPasswordResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "GetCurrentCardPassword",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/password"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin get current card pin request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// GET request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ViewCardPasswordResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(http.StatusInternalServerError, fmt.Sprint("GET_CURRENT_CARD_PASSWORD_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// Embossing - GetCardInfo
func (t *Cards) GetCardInfo(ctx context.Context, requestID string,
	model InfoCardRequest) (*InfoCardResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "GetCardInfo",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/info"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin get card info request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// GET request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *InfoCardResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(http.StatusInternalServerError, fmt.Sprint("GET_CARD_INFO_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

// ResetCardPasswordTries
func (t *Cards) ResetCardPasswordTries(ctx context.Context, requestID string,
	model ResetCardPasswordTriesRequest) (*ResetCardPasswordTriesResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
		"method":     "ResetCardPasswordTries",
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	apiEndpoint := t.session.APIEndpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CardsPath)
	u.Path = path.Join(u.Path, fmt.Sprint("/accounts/", model.AccountID))
	u.Path = path.Join(u.Path, fmt.Sprint("/customers/", model.CustomerID))
	u.Path = path.Join(u.Path, fmt.Sprint("/card/", model.CardID))
	u.Path = path.Join(u.Path, fmt.Sprint("/resetPasswordTries"))

	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("celcoin reset card password tries request")

	// model marshal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshal model")
		return nil, err
	}

	// PATCH request
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	// token
	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error authentication")
		return nil, err
	}

	// request header
	req.Header.Add("Authorization", token)
	// req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var responseBody *ResetCardPasswordTriesResponse

	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	// response ok
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return responseBody, nil
	}

	// error
	if responseBody.Error != nil {
		logrus.WithFields(fields).Error("response body error")
		return nil, grok.NewError(http.StatusInternalServerError, fmt.Sprint("RESET_CARD_PASSWORD_TRIES_ERROR_", responseBody.Error.ErrorCode),
			fmt.Sprint(responseBody.Error.ErrorCode, " - ", responseBody.Error.Message))
	}

	logrus.WithFields(fields).Error("default cards error")
	return nil, ErrDefaultCards
}

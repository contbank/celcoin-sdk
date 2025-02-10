package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"
)

// StatementService ... implementa a interface para o extrato da conta.
type Statement struct {
	session    Session
	httpClient *http.Client
}

// NewStatement ... cria uma nova instância de StatementService.
func NewStatement(httpClient *http.Client, session Session) *Statement {
	return &Statement{
		session:    session,
		httpClient: httpClient,
	}
}

// GetStatements ... realiza a requisição para obter os movimentos da carteira.
func (s *Statement) GetStatements(ctx context.Context,
	request *StatementRequest) (*StatementResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"service":    "Statement",
		"interface":  "GetStatements",
	}

	baseURL := s.session.APIEndpoint

	url, err := url.Parse(baseURL)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing base url")
		return nil, err
	}

	url.Path = path.Join(url.Path, StatementPath)

	if request != nil {
		fields["request"] = request

		q := url.Query()
		if request.Account != nil {
			q.Set("Account", *request.Account)
		}
		if request.DateFrom != nil {
			q.Set("DateFrom", *request.DateFrom)
		}
		if request.DateTo != nil {
			q.Set("DateTo", *request.DateTo)
		}
		if request.DocumentNumber != nil {
			q.Set("DocumentNumber", *request.DocumentNumber)
		}
		if request.LimitPerPage != nil {
			q.Set("LimitPerPage", *request.LimitPerPage)
		}
		url.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	req.Header.Set("accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error making request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error reading response body")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse *ErrorDefaultResponse
		if err := json.Unmarshal(respBody, &errResponse); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshalling error response")
			return nil, err
		}

		if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
			return nil, FindStatementError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		}
	}

	var walletMovementResponse StatementResponse
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&walletMovementResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding response body")
		return nil, err
	}

	return &walletMovementResponse, nil
}

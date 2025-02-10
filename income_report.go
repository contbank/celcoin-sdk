package celcoin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"
)

// IncomeReport ...
type IncomeReport struct {
	session    Session
	httpClient *http.Client
}

// NewIncomeReport ...
func NewIncomeReport(httpClient *http.Client, session Session) *IncomeReport {
	return &IncomeReport{
		session:    session,
		httpClient: httpClient,
	}
}

// GetIncomeReport ...
func (r *IncomeReport) GetIncomeReport(ctx context.Context,
	calendarYear *string, accountNumber *string) (*IncomeReportResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	if calendarYear != nil {
		fields["calendarYear"] = calendarYear
	}

	if accountNumber != nil {
		fields["account"] = accountNumber
	}

	u, err := url.Parse(r.session.APIEndpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, IncomeReportPath)

	q := u.Query()

	if calendarYear != nil && len(*calendarYear) > 0 {
		q.Set("calendarYear", *calendarYear)
	}

	if accountNumber != nil && len(*accountNumber) > 0 {
		q.Set("account", *accountNumber)
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	req.Header.Set("accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error making request")
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error reading response body")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse *ErrorDefaultResponse
		if err := json.Unmarshal(bodyBytes, &errResponse); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
			return nil, FindIncomeReportError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		}
	}

	var incomeReportResponse *IncomeReportResponse
	if err := json.Unmarshal(bodyBytes, &incomeReportResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return incomeReportResponse, nil
}

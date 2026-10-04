package celcoin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/contbank/celcoin-sdk"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/require"
)

// claimTransport answers every request with the same status/body and records what was sent.
type claimTransport struct {
	status   int
	body     string
	requests []*http.Request
	payloads []map[string]interface{}
}

func (t *claimTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.requests = append(t.requests, req)
	payload := map[string]interface{}{}
	if req.Body != nil {
		data, _ := ioutil.ReadAll(req.Body)
		_ = json.Unmarshal(data, &payload)
	}
	t.payloads = append(t.payloads, payload)
	return &http.Response{StatusCode: t.status, Body: ioutil.NopCloser(bytes.NewBufferString(t.body)), Header: make(http.Header)}, nil
}

func newClaimPix(t *testing.T, status int, body string) (*celcoin.Pix, *claimTransport) {
	transport := &claimTransport{status: status, body: body}
	endpoint := "https://sandbox.openfinance.celcoin.dev"
	session, err := celcoin.NewSession(celcoin.Config{
		ClientID:      celcoin.String("id"),
		ClientSecret:  celcoin.String("secret"),
		Mtls:          celcoin.Bool(false),
		APIEndpoint:   &endpoint,
		LoginEndpoint: &endpoint,
	})
	require.NoError(t, err)
	return celcoin.NewPix(&http.Client{Transport: transport}, *session), transport
}

const claimBody = `{"version":"1.0.0","status":"SUCCESS","body":{"id":"c1","claimType":"PORTABILITY","key":"a@b.com",
"keyType":"EMAIL","status":"WAITING_RESOLUTION","donorParticipant":"60701190",
"claimerAccount":{"participant":"13935893","branch":"0001","account":"300541976902"},
"donorAccount":{"account":"123","branch":"1"},"resolutionPeriodEnd":"2026-10-11T10:00:00Z"}}`

// Regression: the SDK only accepted OWNERSHIP, so portability could never be requested.
func TestCreatePixClaimAcceptsPortabilityOnV2(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, claimBody)

	response, err := pix.CreatePixClaim(context.Background(), celcoin.PixClaimRequest{
		Key: "a@b.com", KeyType: "EMAIL", Account: "300541976902", ClaimType: string(celcoin.Portability),
	})
	require.NoError(t, err)
	require.Equal(t, "c1", response.Body.ID)
	require.Equal(t, string(celcoin.WaitingResolution), response.Body.Status, "claim status is in body, not in the envelope")

	req := transport.requests[0]
	require.Equal(t, http.MethodPost, req.Method)
	require.Equal(t, "/baas/v2/pix/dict/claim", req.URL.Path)
	require.Equal(t, "PORTABILITY", transport.payloads[0]["claimType"])
}

func TestCreatePixClaimRejectsInvalidRequestWithoutCallingCelcoin(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, claimBody)

	for _, request := range []celcoin.PixClaimRequest{
		{Key: "0c4e990c-6a2b-44fd-a7cc-0db8c3fa5659", KeyType: "EVP", Account: "1", ClaimType: "OWNERSHIP"},
		{Key: "a@b.com", KeyType: "EMAIL", Account: "1", ClaimType: "TRANSFER"},
		{Key: "a@b.com", KeyType: "EMAIL", Account: "123456789012345678901", ClaimType: "OWNERSHIP"},
	} {
		_, err := pix.CreatePixClaim(context.Background(), request)
		require.Error(t, err, request)
	}
	require.Empty(t, transport.requests)
}

func TestConfirmAndCancelPixClaimPaths(t *testing.T) {
	for _, item := range []struct {
		name string
		call func(*celcoin.Pix) (*celcoin.PixClaimResponse, error)
		path string
	}{
		{"confirm", func(p *celcoin.Pix) (*celcoin.PixClaimResponse, error) {
			return p.ConfirmPixClaim(context.Background(), celcoin.PixClaimActionRequest{ID: "c1"})
		}, "/baas/v2/pix/dict/claim/confirm"},
		{"cancel", func(p *celcoin.Pix) (*celcoin.PixClaimResponse, error) {
			return p.CancelPixClaim(context.Background(), celcoin.PixClaimActionRequest{ID: "c1", Reason: string(celcoin.Fraud)})
		}, "/baas/v2/pix/dict/claim/cancel"},
	} {
		pix, transport := newClaimPix(t, http.StatusOK, claimBody)
		_, err := item.call(pix)
		require.NoError(t, err, item.name)
		require.Equal(t, http.MethodPost, transport.requests[0].Method, item.name)
		require.Equal(t, item.path, transport.requests[0].URL.Path, item.name)
		require.Equal(t, "c1", transport.payloads[0]["id"], item.name)
	}
}

func TestPixClaimActionReasonIsOptionalButValidated(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, claimBody)

	_, err := pix.ConfirmPixClaim(context.Background(), celcoin.PixClaimActionRequest{ID: "c1"})
	require.NoError(t, err)
	_, hasReason := transport.payloads[0]["reason"]
	require.False(t, hasReason, "empty reason is omitted (Celcoin default USER_REQUESTED)")

	_, err = pix.CancelPixClaim(context.Background(), celcoin.PixClaimActionRequest{ID: "c1", Reason: string(celcoin.DonorRequest)})
	require.Error(t, err, "DONOR_REQUEST is not accepted by the v2 API")
	require.Len(t, transport.requests, 1)
}

func TestCancelPixClaimNotPendingError(t *testing.T) {
	pix, _ := newClaimPix(t, http.StatusBadRequest,
		`{"version":"1.0.0","status":"ERROR","error":{"errorCode":"CBE307","message":"Não foi possível cancelar essa Claim"}}`)

	_, err := pix.CancelPixClaim(context.Background(), celcoin.PixClaimActionRequest{ID: "c1"})
	grokErr, ok := err.(*grok.Error)
	require.True(t, ok)
	require.Equal(t, "CLAIM_NOT_PENDING", grokErr.Key)
	require.Equal(t, http.StatusBadRequest, grokErr.Code)
}

func TestGetPixClaimAndList(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, claimBody)
	response, err := pix.GetPixClaim(context.Background(), "c1")
	require.NoError(t, err)
	require.Equal(t, "/baas/v2/pix/dict/claim/c1", transport.requests[0].URL.Path)
	require.Equal(t, "60701190", response.Body.DonorParticipant)

	pix, transport = newClaimPix(t, http.StatusOK,
		`{"version":"1.0.0","status":"SUCCESS","body":{"claims":[{"id":"c1","status":"CANCELLED","cancelledBy":"DONOR"}]}}`)
	list, err := pix.GetPixClaimList(context.Background(), "2026-10-01", "", 0, 0, string(celcoin.CanceledClaim), "")
	require.NoError(t, err)
	require.Equal(t, string(celcoin.CanceledClaim), list.Body.Claims[0].Status)

	url := transport.requests[0].URL
	require.Equal(t, "/baas/v2/pix/dict/claim/list", url.Path)
	require.Equal(t, "2026-10-01", url.Query().Get("DateFrom"))
	require.Equal(t, "CANCELLED", url.Query().Get("Status"))
	require.False(t, url.Query().Has("LimitPerPage"), "zero limit uses the Celcoin default")
	require.False(t, url.Query().Has("Page"))
}

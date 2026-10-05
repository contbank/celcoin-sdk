package celcoin_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/contbank/celcoin-sdk"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/require"
)

const reversalBody = `{"status":"PROCESSING","version":"1.0.0","body":{"id":"rev-1","amount":100,"clientCode":"cc-1",
"originalPaymentId":"7b839966-49ac-4749-9ec0-ee51bb84ce00","endToEndId":"E13935893202610032327ioMIdLbtDKM",
"returnIdentification":"D13935893202610040000abcdefghijk","reason":"MD06","reversalDescription":"Devolução"}}`

func reversalRequest() celcoin.PixReversalRequest {
	return celcoin.PixReversalRequest{
		ID:         "7b839966-49ac-4749-9ec0-ee51bb84ce00",
		EndToEndID: "E13935893202610032327ioMIdLbtDKM",
		ClientCode: "cc-1",
		Amount:     100,
		Reason:     string(celcoin.ReversalRequestedByCustomer),
	}
}

func TestReversePixCashInOnV2(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, reversalBody)

	response, err := pix.ReversePixCashIn(context.Background(), reversalRequest())
	require.NoError(t, err)
	require.Equal(t, "PROCESSING", response.Status)
	require.Equal(t, "D13935893202610040000abcdefghijk", response.Body.ReturnIdentification)
	require.Equal(t, "7b839966-49ac-4749-9ec0-ee51bb84ce00", response.Body.OriginalPaymentID)

	req := transport.requests[0]
	require.Equal(t, http.MethodPost, req.Method)
	require.Equal(t, "/baas/v2/pix/reverse", req.URL.Path)
	payload := transport.payloads[0]
	require.Equal(t, "cc-1", payload["clientCode"])
	require.Equal(t, "MD06", payload["reason"])
	require.Equal(t, float64(100), payload["amount"])
	require.Equal(t, "E13935893202610032327ioMIdLbtDKM", payload["endToEndId"])
	_, hasDescription := payload["reversalDescription"]
	require.False(t, hasDescription, "empty description is omitted")
}

func TestReversePixCashInValidatesBeforeCallingCelcoin(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, reversalBody)

	invalid := []func(r *celcoin.PixReversalRequest){
		func(r *celcoin.PixReversalRequest) { r.ID, r.EndToEndID = "", "" },                      // CBE204
		func(r *celcoin.PixReversalRequest) { r.ClientCode = "" },                                // CBE001
		func(r *celcoin.PixReversalRequest) { r.Amount = 0 },                                     // CBE094/095
		func(r *celcoin.PixReversalRequest) { r.Reason = "XX01" },                                // CBE155
		func(r *celcoin.PixReversalRequest) { r.ReversalDescription = strings.Repeat("a", 141) }, // CBE156
	}
	for i, change := range invalid {
		request := reversalRequest()
		change(&request)
		_, err := pix.ReversePixCashIn(context.Background(), request)
		require.Error(t, err, i)
	}

	request := reversalRequest()
	request.ID = ""
	_, err := pix.ReversePixCashIn(context.Background(), request)
	require.NoError(t, err, "endToEndId alone identifies the receipt")
	require.Len(t, transport.requests, 1)
}

func TestReversePixCashInErrors(t *testing.T) {
	for code, key := range map[string]string{
		"CBE157": "REVERSAL_AMOUNT_EXCEEDED",
		"CBE158": "REVERSAL_DEADLINE_EXPIRED",
		"CBE123": "INSUFFICIENT_BALANCE",
		"CBE101": "DUPLICATE_TRANSACTION",
	} {
		pix, _ := newClaimPix(t, http.StatusBadRequest,
			`{"version":"1.0.0","status":"ERROR","error":{"errorCode":"`+code+`","message":"x"}}`)
		_, err := pix.ReversePixCashIn(context.Background(), reversalRequest())
		grokErr, ok := err.(*grok.Error)
		require.True(t, ok, code)
		require.Equal(t, key, grokErr.Key, code)
		require.Equal(t, http.StatusBadRequest, grokErr.Code, code)
	}
}

func TestGetPixReversalStatus(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK, strings.Replace(reversalBody, "PROCESSING", "CONFIRMED", 1))

	response, err := pix.GetPixReversalStatus(context.Background(), "", "", "D13935893202610040000abcdefghijk")
	require.NoError(t, err)
	require.Equal(t, "CONFIRMED", response.Status)

	url := transport.requests[0].URL
	require.Equal(t, http.MethodGet, transport.requests[0].Method)
	require.Equal(t, "/baas/v2/pix/reverse/status", url.Path)
	require.Equal(t, "D13935893202610040000abcdefghijk", url.Query().Get("returnIdentification"))
	require.False(t, url.Query().Has("id"))

	_, err = pix.GetPixReversalStatus(context.Background(), "", "", "")
	require.Equal(t, celcoin.ErrPixReversalStatusParameters, err)
	require.Len(t, transport.requests, 1, "no call without identifier")
}

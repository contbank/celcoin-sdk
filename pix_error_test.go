package celcoin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/contbank/celcoin-sdk"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/require"
)

// Regression (04/10/2026): EMAIL key refused with CBE197 reached the app as 500 UNKNOWN_ERROR "unknown error".
func TestCreatePixKeyMapsCBE197(t *testing.T) {
	pix, _ := newClaimPix(t, http.StatusBadRequest,
		`{"status":"ERROR","error":{"errorCode":"CBE197","message":"Tipo de chave não permitido."},"version":"1.0.0"}`)

	_, err := pix.CreatePixKey(context.Background(), celcoin.PixKeyRequest{Account: "1", KeyType: "EMAIL", Key: "a@b.com"})
	grokErr, ok := err.(*grok.Error)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, grokErr.Code)
	require.Equal(t, "KEY_TYPE_NOT_ALLOWED", grokErr.Key)
}

// An unmapped Celcoin code keeps its 4xx and its message (already in Portuguese) for the app.
func TestUnmappedCelcoinErrorKeepsPartnerMessage(t *testing.T) {
	status := http.StatusUnprocessableEntity
	message := "Mensagem nova da Celcoin."
	err := celcoin.FindPixErrorWithMessage("CBE999999", &status, &message)
	require.Equal(t, http.StatusUnprocessableEntity, err.Code)
	require.Equal(t, "PIX_PARTNER_ERROR", err.Key)
	require.Equal(t, []string{message}, err.Messages)

	status = http.StatusInternalServerError
	err = celcoin.FindPixErrorWithMessage("CBE999999", &status, nil)
	require.Equal(t, http.StatusBadGateway, err.Code)
	require.NotEqual(t, "unknown error", err.Messages[0])
}

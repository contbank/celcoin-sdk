package celcoin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/contbank/celcoin-sdk"
	"github.com/stretchr/testify/require"
)

// Chaves da conta na rota documentada (BaaS v2); a v1 legada ficou como PixDictPathDeprecated.
func TestPixKeyOperationsUseV2(t *testing.T) {
	pix, transport := newClaimPix(t, http.StatusOK,
		`{"version":"1.0.0","status":"CONFIRMED","body":{"keyType":"EMAIL","key":"a@b.com","account":{"account":"421627290"}}}`)
	_, err := pix.CreatePixKey(context.Background(), celcoin.PixKeyRequest{Account: "421627290", KeyType: "EMAIL", Key: "a@b.com"})
	require.NoError(t, err)
	require.Equal(t, http.MethodPost, transport.requests[0].Method)
	require.Equal(t, "/baas/v2/pix/dict/entry", transport.requests[0].URL.Path)
	require.Equal(t, "EMAIL", transport.payloads[0]["keyType"])

	pix, transport = newClaimPix(t, http.StatusOK, `{"version":"1.0.0","status":"SUCCESS","body":{"listKeys":[]}}`)
	_, err = pix.GetPixKeys(context.Background(), "421627290")
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, transport.requests[0].Method)
	require.Equal(t, "/baas/v2/pix/dict/entry/421627290", transport.requests[0].URL.Path)

	pix, transport = newClaimPix(t, http.StatusOK, `{"version":"1.0.0","status":"SUCCESS"}`)
	require.NoError(t, pix.DeletePixKey(context.Background(), "421627290", "a@b.com"))
	require.Equal(t, http.MethodDelete, transport.requests[0].Method)
	require.Equal(t, "/baas/v2/pix/dict/entry/a@b.com", transport.requests[0].URL.Path)
	require.Equal(t, "421627290", transport.payloads[0]["account"])
}

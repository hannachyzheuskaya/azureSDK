package util

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"net/http"
)

const (
	SessionName        = "sessionAzure"
	CtxKeyConn  ctxKey = iota
)

var (
	errNotAuthenticated = errors.New("not authenticated")
	errAuthentication   = errors.New("unable to retrieve authentication information")
)

type AzureAuthInfo struct {
	ClientID       string
	ClientSecret   string
	TenantID       string
	SubscriptionId string
}

type ctxKey int8

func Credential(r *http.Request) (*azidentity.ClientSecretCredential, string, error) {
	authInfo, ok := r.Context().Value(CtxKeyConn).(AzureAuthInfo)
	if !ok {
		return nil, "-", errAuthentication
	}
	cred, err := azidentity.NewClientSecretCredential(authInfo.TenantID, authInfo.ClientID, authInfo.ClientSecret, nil)
	if err != nil {
		return nil, "-", errAuthentication
	}
	return cred, authInfo.SubscriptionId, nil
}

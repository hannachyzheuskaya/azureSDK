package connectionManager

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"io"
	"log/slog"
	"net/http"
	"x/internal/app/controller/handlers/util"
	sl "x/internal/lib/logger"
	resp "x/internal/lib/response"
)

func init() {
	gob.Register(util.AzureAuthInfo{})
}

type Response struct {
	resp.Response
	Message string `json:"message"`
}

func NewConnection(log *slog.Logger, sessionStore sessions.Store) http.HandlerFunc {
	type request struct {
		ClientID       string `json:"client_id"`
		ClientSecret   string `json:"client_secret"`
		TenantID       string `json:"tenant_id"`
		SubscriptionId string `json:"subscription_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.connectionManager.NewConnection"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		req := &request{}

		err := render.DecodeJSON(r.Body, req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded")

		session, err := sessionStore.Get(r, util.SessionName)
		if err != nil {
			log.Error("failed to create session", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to create session"))
			return
		}

		_, err = azidentity.NewClientSecretCredential(req.TenantID, req.ClientID, req.ClientSecret, nil)
		if err != nil {
			log.Error("unauthorized", sl.Err(err))
			render.JSON(w, r, resp.Error("unauthorized"))
			return
		}

		authInfo := &util.AzureAuthInfo{
			ClientID:       req.ClientID,
			ClientSecret:   req.ClientSecret,
			TenantID:       req.TenantID,
			SubscriptionId: req.SubscriptionId,
		}

		session.Values["authInfo"] = authInfo
		if err := sessionStore.Save(r, w, session); err != nil {
			log.Error("failed to save session", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save session"))
			return
		}

		responseOK(w, r, "connected successfully to azure")

	}
}

func NewAuthenticateUser(log *slog.Logger, sessionStore sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			const op = "handlers.connectionManager.NewAuthenticateUser"

			log := log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			session, err := sessionStore.Get(r, util.SessionName)
			if err != nil {
				log.Error("not authenticated", sl.Err(err))
				render.JSON(w, r, resp.Error("not authenticated"))
				return
			}
			authInfo, ok := session.Values["authInfo"]
			if !ok {
				log.Error("unable to retrieve authentication information", sl.Err(err))
				render.JSON(w, r, resp.Error("unable to retrieve authentication information"))
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.CtxKeyConn, authInfo)))
		}
		return http.HandlerFunc(fn)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, msg string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Message:  msg,
	})
}

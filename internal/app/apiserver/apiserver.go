package apiserver

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"x/internal/app/config"
	"x/internal/app/controller"
)

func Start(log *slog.Logger, config *config.Config) error {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	router := controller.NewRouter(log, sessionStore)
	srv := &http.Server{
		Addr:         config.BindAddr,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	//...
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"x/internal/app/controller/handlers/connectionManager"
	"x/internal/app/controller/handlers/virtualMachine"
	"x/internal/app/controller/logger"
)

type Router struct {
	router       *chi.Mux
	logger       *slog.Logger
	sessionStore sessions.Store
}

func NewRouter(log *slog.Logger, sessionStore sessions.Store) *Router {
	r := &Router{
		router:       chi.NewRouter(),
		logger:       log,
		sessionStore: sessionStore,
	}
	r.configureRouter()
	return r
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.router.ServeHTTP(w, r)
}

func (rt *Router) configureRouter() {
	rt.router.Use(middleware.RequestID)
	rt.router.Use(logger.New(rt.logger))
	rt.router.Post("/createSession", connectionManager.NewConnection(rt.logger, rt.sessionStore))
	rt.router.Route("/virtualMachine", func(r chi.Router) {
		r.Use(connectionManager.NewAuthenticateUser(rt.logger, rt.sessionStore))
		r.Post("/create", virtualMachine.New(rt.logger))
	})

	//rt.router.HandleFunc("/createSession", rt.handleCreateSession()).Methods(http.MethodPost)

	//vm := rt.router.PathPrefix("/virtualMachine").Subrouter()
	//vm.Use(s.authenticateUser)
	//vm.HandleFunc("/create", s.handleTestRequest()).Methods(http.MethodPost)

	//rg := s.router.PathPrefix("/resourcegroup").Subrouter()
	//rg.Use(s.authenticateUser)
	//rg.HandleFunc("/create", s.handleCreateRG()).Methods(http.MethodPost)
	//rg.HandleFunc("/delete", s.handleDeleteRG()).Methods(http.MethodPost)
	//rg.HandleFunc("/list", s.listAllRG()).Methods(http.MethodPost)

}

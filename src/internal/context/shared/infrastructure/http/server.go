package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_infrastructure_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer         *http.Server
	defaultMiddlewares []Middleware
	logger             shared_domain_context.Logger
}

func NewServer(logger shared_domain_context.Logger, router http.Handler, conf *shared_infrastructure_context.Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", conf.Server.Port),
			Handler:      router,
			ReadTimeout:  time.Duration(conf.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(conf.Server.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(conf.Server.IdleTimeout) * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s.logger.Info(context.Background(), fmt.Sprintf("Starting HTTP server on %s", s.httpServer.Addr))
		s.httpServer.ListenAndServe()
	}()

	<-shutdown
	s.logger.Info(context.Background(), "Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info(ctx, "Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}

type Route interface {
	Handler(w http.ResponseWriter, r *http.Request)
	Route() string
	Method() string
}

type Router struct {
	*http.ServeMux
}

func NewRouter(middlewares []Middleware, routes ...Route) *Router {
	router := &Router{ServeMux: http.NewServeMux()}

	router.Handle("/metrics", promhttp.Handler())

	for _, route := range routes {
		router.RegisterRoute(route, middlewares...)
	}

	return router
}

func (r *Router) RegisterRoute(route Route, middlewares ...Middleware) {
	handler := http.HandlerFunc(route.Handler)

	for _, m := range middlewares {
		handler = m.Handle(route)
	}

	r.HandleFunc(route.Route(), handler)
}

type Middleware interface {
	Handle(next Route) http.HandlerFunc
}

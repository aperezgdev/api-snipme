package main

import (
	"context"
	"log"

	"github.com/aperezgdev/api-snipme/db/generated"
	link_visit_creator "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/application"
	link_visit_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/infrastructure"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_infrastructure_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http/middleware"
	short_link_application "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	short_link_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure"
	short_link_http "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conf := shared_infrastructure_context.Load()
	logger := shared_domain_context.NewConsoleLogger()
	eventBus := shared_domain_context.NewEventBusInMemory()

	ctx := context.Background()
	config, err := pgxpool.ParseConfig(conf.Database.Url)
	if err != nil {
		logger.Error(ctx, "Error parsing database URL", shared_domain_context.NewField("error", err))
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Error creating pool:", err)
	}
	defer pool.Close()

	queries := generated.New(pool)

	shortLinkRepository := short_link_infrastructure.NewSqlcShortLinkRepository(queries)
	linkVisitRepository := link_visit_infrastructure.NewSqlcLinkVisitRepository(queries)

	shortLinkFinderByCode := short_link_application.NewShortLinkFinderByCode(logger, shortLinkRepository)
	shortLinkCreator := short_link_application.NewShortLinkCreator(logger, shortLinkRepository, &eventBus)

	linkVisitCreator := link_visit_creator.NewLinkVisitCreator(logger, linkVisitRepository)

	getShortLink := short_link_http.NewGetShortLinkByCodeHTTPHandler(logger, *shortLinkFinderByCode, *linkVisitCreator)
	postShortLink := short_link_http.NewPostShortLinkHTTPHandler(logger, *shortLinkCreator)

	router := http.NewRouter([]http.Middleware{middleware.NewRecoveryMiddleware(logger), middleware.NewRequestIDMiddleware(logger), middleware.NewLoggerMiddleware(logger), middleware.NewPrometheusMiddleware()}, getShortLink, postShortLink)

	server := http.NewServer(logger, router, conf)

	server.Start()
}

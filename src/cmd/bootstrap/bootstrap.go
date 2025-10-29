package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	link_visit_creator "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/application"
	link_visit_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/infrastructure"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_infrastructure_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure"
	shared_cache "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/cache"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	shared_infrastructure_http_handler "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http/handler"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http/middleware"
	short_link_application "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	short_link_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure"
	short_link_cache "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure/cache"
	short_link_http "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

func Run() error {
	conf := shared_infrastructure_context.Load()
	var logger shared_domain_context.Logger = shared_domain_context.NewConsoleLogger()
	eventBus := shared_domain_context.NewEventBusInMemory()

	if conf.Loki.Url != "" {
		lokiLogger := shared_infrastructure_context.NewLokiLogger(conf.Loki.Url)
		consoleLogger := shared_domain_context.NewConsoleLogger()
		logger = shared_domain_context.NewCompositeLogger(consoleLogger, lokiLogger)
	}

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Url,
		Password: conf.Redis.Password,
		DB:       0,
	})
	defer redisClient.Close()

	cache := shared_cache.NewRedisCache(redisClient)
	shortLinkRepository := short_link_cache.NewRedisShortLinkRepository(
		short_link_infrastructure.NewSqlcShortLinkRepository(logger, queries),
		cache,
		5*time.Minute,
		logger,
	)

	linkVisitRepository := link_visit_infrastructure.NewSqlcLinkVisitRepository(logger, queries)

	shortLinkFinderByCode := short_link_application.NewShortLinkFinderByCode(logger, shortLinkRepository)
	shortLinkCreator := short_link_application.NewShortLinkCreator(logger, shortLinkRepository, &eventBus)

	linkVisitCreator := link_visit_creator.NewLinkVisitCreator(logger, linkVisitRepository)

	getStatus := shared_infrastructure_http_handler.NewGetStatusHTTPHandler()

	getShortLink := short_link_http.NewGetShortLinkByCodeHTTPHandler(logger, *shortLinkFinderByCode, *linkVisitCreator)
	postShortLink := short_link_http.NewPostShortLinkHTTPHandler(logger, *shortLinkCreator)

	router := http.NewRouter([]http.Middleware{
		middleware.NewRecoveryMiddleware(logger),
		middleware.NewLoggerMiddleware(logger),
		middleware.NewPrometheusMiddleware(),
		middleware.NewRequestIDMiddleware(logger),
		middleware.NewRateLimitMiddleware(logger, rate.Every(100*time.Millisecond), 5),
	}, getStatus, getShortLink, postShortLink)

	server := http.NewServer(logger, router, conf)

	return server.Start()
}

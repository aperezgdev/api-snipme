package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	auth_application "github.com/aperezgdev/api-snipme/src/internal/context/authentication/application"
	auth_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	auth_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	auth_http "github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure/http"
	geo_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/infrastructure"
	link_analytics_application "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/application"
	link_analytics_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/infrastructure"
	link_analytics_http "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/infrastructure/http"
	link_visit_creator "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/application"
	link_visit_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/infrastructure"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_infrastructure_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure"
	shared_cache "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/cache"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	shared_infrastructure_http_handler "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http/handler"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http/middleware"
	client_application "github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/application"
	client_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/infrastructure"
	short_link_application "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	short_link_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure"
	"github.com/oschwald/maxminddb-golang/v2"
	"github.com/robfig/cron/v3"

	link_country_view_counter "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/application"
	link_country_view_counter_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/infrastructure"

	link_visit_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
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

	db, err := maxminddb.Open(conf.GEOFilePath)
	if err != nil {
		logger.Error(ctx, "Error opening MaxMind DB", shared_domain_context.NewField("error", err))
	}
	defer db.Close()

	cache := shared_cache.NewRedisCache(logger, redisClient)
	shortLinkRepository := short_link_cache.NewRedisShortLinkRepository(
		short_link_infrastructure.NewSqlcShortLinkRepository(logger, queries),
		cache,
		5*time.Minute,
		logger,
	)

	clientRepo := client_infrastructure.NewSqlcClientRepository(logger, queries)
	linkVisitRepository := link_visit_infrastructure.NewSqlcLinkVisitRepository(logger, queries)
	linkAnalytics := link_analytics_infrastructure.NewSqlcLinkAnalyticsRepository(logger, queries)
	userRepository := auth_infrastructure.NewSqlcUserRepository(logger, queries)
	refreshTokenRepository := auth_infrastructure.NewSqlcRefreshTokenRepository(logger, queries)
	linkCountryViewCounterRepository := link_country_view_counter_infrastructure.NewSqlcLinkCountryViewCounterRepository(logger, queries)
	geoRepository := geo_infrastructure.NewMMDBRepository(logger, *db)

	jwtManager := auth_infrastructure.NewJWTManager(conf.JWT.Secret, conf.JWT.ExpirationMinutes)

	googleOAuthClient := auth_infrastructure.NewGoogleOAuthClient(
		conf.OAuth.GoogleClientID,
		conf.OAuth.GoogleClientSecret,
		conf.OAuth.GoogleRedirectURL,
	)
	githubOAuthClient := auth_infrastructure.NewGitHubOAuthClient(
		conf.OAuth.GitHubClientID,
		conf.OAuth.GitHubClientSecret,
		conf.OAuth.GitHubRedirectURL,
	)

	authenticator := auth_application.NewAuthenticator(
		logger,
		userRepository,
		refreshTokenRepository,
		jwtManager,
		&eventBus,
		conf.JWT.RefreshTokenTTLDays,
		conf.JWT.ExpirationMinutes,
	)

	tokenValidator := auth_application.NewTokenValidator(logger, jwtManager, userRepository)
	tokenRefresher := auth_application.NewTokenRefresher(logger, refreshTokenRepository, userRepository, jwtManager, conf.JWT.ExpirationMinutes)

	shortLinkFinderByCode := short_link_application.NewShortLinkFinderByCode(logger, shortLinkRepository)
	shortLinkFinderByClient := short_link_application.NewShortLinkFinderByClient(logger, shortLinkRepository, clientRepo)
	shortLinkCreator := short_link_application.NewShortLinkCreator(logger, shortLinkRepository, &eventBus)
	publShortLinkCreator := short_link_application.NewPublicShortLinkCreator(logger, shortLinkRepository, &eventBus)
	shortLinkRemover := short_link_application.NewShortLinkRemover(logger, shortLinkRepository)

	linkVisitCreator := link_visit_creator.NewLinkVisitCreator(logger, linkVisitRepository)
	linkVisitProcessor := link_visit_creator.NewLinkVisitProcessor(logger, linkVisitRepository, &eventBus)

	getStatus := shared_infrastructure_http_handler.NewGetStatusHTTPHandler()

	getShortLink := short_link_http.NewGetShortLinkByCodeHTTPHandler(logger, *shortLinkFinderByCode, *linkVisitCreator)
	getLinkAnalyticsByLink := link_analytics_http.NewGetLinkAnalyticsByLinkHTTPHandler(logger, *link_analytics_application.NewLinkAnalyticsFinder(logger, linkAnalytics))
	postShortLink := short_link_http.NewPostShortLinkHTTPHandler(logger, *shortLinkCreator)
	deleteShortLink := short_link_http.NewDeleteShortLinkHTTPHandler(logger, *shortLinkRemover)
	getShortLinkByClient := short_link_http.NewGetShortLinkByClientHTTPHandler(logger, *shortLinkFinderByClient)
	postPublicShortLink := short_link_http.NewPostPublicShortLinkHTTPHandler(logger, *publShortLinkCreator)

	googleLoginHandler := auth_http.NewGetOAuthLoginHandler(logger, googleOAuthClient, "google", conf.OAuth.StateSecret)
	googleCallbackHandler := auth_http.NewGetOAuthCallbackHandler(logger, googleOAuthClient, authenticator, auth_domain.OAuthProviderGoogle, conf.OAuth.StateSecret)
	githubLoginHandler := auth_http.NewGetOAuthLoginHandler(logger, githubOAuthClient, "github", conf.OAuth.StateSecret)
	githubCallbackHandler := auth_http.NewGetOAuthCallbackHandler(logger, githubOAuthClient, authenticator, auth_domain.OAuthProviderGitHub, conf.OAuth.StateSecret)
	refreshTokenHandler := auth_http.NewPostRefreshTokenHandler(logger, tokenRefresher)

	incrementerOnLinkVisitProcessed := link_country_view_counter.NewIncremeterOnLinkVisitProcessed(logger, geoRepository, linkCountryViewCounterRepository)
	updaterOnLinkVisitProcessed := link_analytics_application.NewUpdaterOnLinkVisitProcessed(logger, linkAnalytics)
	eventBus.AddSubscribers(link_visit_domain.LinkVisitsProcessedEventName, updaterOnLinkVisitProcessed)
	creatorOnShorLinkCreated := link_analytics_application.NewCreatorOnShortLinkCreated(logger, linkAnalytics)
	eventBus.AddSubscribers("ShortLinkCreated", creatorOnShorLinkCreated)
	creatorOnUserCreated := client_application.NewCreatorOnUserCreated(logger, clientRepo)
	eventBus.AddSubscribers(auth_domain.UserCreatedEventName, creatorOnUserCreated)
	eventBus.AddSubscribers(link_visit_domain.LinkVisitsProcessedEventName, incrementerOnLinkVisitProcessed)

	c := cron.New()
	c.AddFunc("@every 5m", func() {
		err := linkVisitProcessor.Run(context.Background())
		if err != nil {
			logger.Error(ctx, "Error running LinkVisitProcessor cron job", shared_domain_context.NewField("error", err))
		}
	})

	c.Start()

	globalMiddlewares := []http.Middleware{
		middleware.NewRecoveryMiddleware(logger),
		middleware.NewLoggerMiddleware(logger),
		middleware.NewPrometheusMiddleware(),
		middleware.NewRequestIDMiddleware(logger),
		middleware.NewRateLimitMiddleware(logger, rate.Every(100*time.Millisecond), 5),
	}

	authMiddleware := middleware.NewAuthenticationMiddleware(logger, tokenValidator)
	protectedMiddlewares := append(globalMiddlewares, authMiddleware)

	publicRoutes := []http.Route{
		getStatus,
		getShortLink,
		googleLoginHandler,
		googleCallbackHandler,
		githubLoginHandler,
		githubCallbackHandler,
		refreshTokenHandler,
		postPublicShortLink,
	}

	protectedRoutes := []http.Route{
		postShortLink,
		deleteShortLink,
		getShortLinkByClient,
		getLinkAnalyticsByLink,
	}

	router := http.NewRouter(globalMiddlewares, publicRoutes...)

	for _, route := range protectedRoutes {
		router.RegisterRoute(route, protectedMiddlewares...)
	}

	server := http.NewServer(logger, router, conf)

	return server.Start()
}

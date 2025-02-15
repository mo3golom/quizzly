package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"quizzly/internal/quizzly"
	"quizzly/pkg/cookie"
	"quizzly/pkg/files"
	"quizzly/pkg/logger"
	"quizzly/pkg/structs"
	"quizzly/pkg/supabase"
	variablesRepo "quizzly/pkg/variables"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/handlers/admin/game"
	"quizzly/web/frontend/handlers/admin/question"
	"quizzly/web/frontend/handlers/admin/static/faq"
	files2 "quizzly/web/frontend/handlers/files"
	gamePublic "quizzly/web/frontend/handlers/public/game"
	"quizzly/web/frontend/handlers/public/login"
	"quizzly/web/frontend/services/link"
	playerService "quizzly/web/frontend/services/player"
	sessionService "quizzly/web/frontend/services/session"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
)

const (
	publicPath = "web/frontend/public"
	serverAddr = ":3000"
)

const (
	ServerTypeHttp   serverType = "http"
	ServerTypeLambda serverType = "lambda"
)

type (
	serverType string

	configuration struct {
		sessions structs.Singleton[sessionService.Service]
		player   structs.Singleton[playerService.Service]
		link     structs.Singleton[link.Service]
	}

	serverSettings struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	muxExtended struct {
		mux        *http.ServeMux
		middleware middleware.Middleware
	}

	ServerInstance struct {
		serverLambda func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
		serverHTTP   *http.Server
		serverType   serverType

		log logger.Logger
	}
)

func (m *muxExtended) HandleFunc(pattern string, metricsKey string, handler func(http.ResponseWriter, *http.Request)) {
	corsFn := func(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			delegate(w, r)
		}
	}

	m.mux.Handle(pattern, middlewarestd.Handler(metricsKey, m.middleware, http.HandlerFunc(corsFn(handler))))
}

func adminRoutes(
	mux *muxExtended,
	config *configuration,
	log logger.Logger,
	quizzlyConfig *quizzly.Configuration,
	authClient supabase.Auth,
	filesManager files.Manager,
) {
	security := authClient.MiddlewareAuth

	mux.HandleFunc("POST /admin/question", "/admin/question", security(handlers.Templ[question.NewPostData](question.NewPostCreateHandler(
		quizzlyConfig.Game.MustGet(),
		filesManager,
	), log)))
	mux.HandleFunc("DELETE /admin/question", "/admin/question", security(handlers.Templ[question.GetDeleteData](question.NewPostDeleteHandler(
		quizzlyConfig.Game.MustGet(),
	), log)))
	mux.HandleFunc("GET /admin/question/list", "/admin/question/list", security(handlers.Templ[question.GetListData](question.NewGetHandler(quizzlyConfig.Game.MustGet()), log)))

	mux.HandleFunc("GET /admin/game/new", "/admin/game/new", security(handlers.Templ[struct{}](game.NewGetCreateHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("GET /admin/game/{game_id}", "/admin/game/:game_id", security(handlers.Templ[game.GetGamePageData](game.NewGetPageHandler(
		quizzlyConfig.Game.MustGet(),
		config.link.MustGet(),
	), log)))
	mux.HandleFunc("POST /admin/game/{game_id}/update", "/admin/game/:game_id/update", security(handlers.Templ[game.PostUpdateData](game.NewPostUpdateHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("POST /admin/game/start", "/admin/game/start", security(handlers.Templ[game.PostStartData](game.NewPostStartHandler(
		quizzlyConfig.Game.MustGet(),
		config.link.MustGet(),
	), log)))
	mux.HandleFunc("POST /admin/game/finish", "/admin/game/finish", security(handlers.Templ[game.PostFinishData](game.NewPostFinishHandler(
		quizzlyConfig.Game.MustGet(),
		config.link.MustGet(),
	), log)))

	mux.HandleFunc("GET /admin/game/list", "/admin/game/list", security(handlers.Templ[struct{}](game.NewGetListHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("GET /admin/game/session/list", "/admin/game/session/list", security(handlers.Templ[game.GetSessionListData](game.NewGetSessionListHandler(config.sessions.MustGet()), log)))

	mux.HandleFunc("GET /admin/faq", "/admin/faq", security(handlers.Templ[struct{}](faq.NewStaticFAQHandler(), log)))
}

func publicRoutes(
	mux *muxExtended,
	log logger.Logger,
	config *configuration,
	quizzlyConfig *quizzly.Configuration,
	authClient supabase.Auth,
) {
	security := authClient.MiddlewareTrace

	mux.HandleFunc("GET /login", "/login", security(handlers.Templ[login.GetLoginPageData](login.NewGetLoginPageHandler(authClient), log)))
	mux.HandleFunc("GET /logout", "/logout", security(handlers.Templ[struct{}](login.NewGetLogoutPageHandler(authClient), log)))

	gamePlayPageHandler := gamePublic.NewGetPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		config.player.MustGet(),
		config.link.MustGet(),
	)
	gameRestartPageHandler := gamePublic.NewGetRestartPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		config.player.MustGet(),
		config.link.MustGet(),
	)
	gameResultsPagehandler := gamePublic.NewGetPlayResultsPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
		config.player.MustGet(),
		config.link.MustGet(),
	)
	gameRenamePlayerHandler := gamePublic.NewPostRenamePlayerHandler(
		quizzlyConfig.Player.MustGet(),
		config.player.MustGet(),
	)

	mux.HandleFunc("GET /", "/", security(handlers.Templ[gamePublic.GetStartPageData](gamePublic.NewGetStartPageHandler(
		quizzlyConfig.Game.MustGet(),
	), log)))

	mux.HandleFunc("GET /game/{game_id}", "/game/:game_id", security(handlers.Templ[gamePublic.GetPlayPageData](gamePlayPageHandler, log)))
	// backwards compatibility
	mux.HandleFunc("GET /game/play", "/game/:game_id (old)", security(handlers.Templ[gamePublic.GetPlayPageData](gamePlayPageHandler, log)))

	mux.HandleFunc("POST /game/{game_id}", "/game/:game_id", handlers.Templ[gamePublic.PostPlayPageData](gamePublic.NewPostPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		config.player.MustGet(),
		config.link.MustGet(),
	), log))

	mux.HandleFunc("GET /game/{game_id}/restart", "/game/:game_id/restart", security(handlers.Templ[gamePublic.GetRestartPageData](gameRestartPageHandler, log)))
	// backwards compatibility
	mux.HandleFunc("GET /game/restart", "/game/:game_id/restart (old)", security(handlers.Templ[gamePublic.GetRestartPageData](gameRestartPageHandler, log)))

	mux.HandleFunc("GET /game/{game_id}/results/{player_id}", "/game/:game_id/results/:player_id", security(handlers.Templ[gamePublic.GetPlayResultsPageData](gameResultsPagehandler, log)))
	// backwards compatibility
	mux.HandleFunc("GET /game/results", "/game/:game_id/results/:player_id (old)", security(handlers.Templ[gamePublic.GetPlayResultsPageData](gameResultsPagehandler, log)))

	mux.HandleFunc("POST /game/{game_id}/player/{player_id}/rename", "/game/:game_id/player/:player_id/rename", security(handlers.Templ[gamePublic.PostRenamePlayerData](gameRenamePlayerHandler, log)))
}

func NewServer(
	log logger.Logger,
	variables variablesRepo.Repository,
	quizzlyConfig *quizzly.Configuration,
	authClient supabase.Auth,
	cookieService cookie.Service,
	filesManager files.Manager,
	serverType serverType,
) *ServerInstance {
	config := &configuration{
		sessions: structs.NewSingleton(func() (sessionService.Service, error) {
			return sessionService.NewService(
				quizzlyConfig.Session.MustGet(),
				quizzlyConfig.Player.MustGet(),
			), nil
		}),
		player: structs.NewSingleton(func() (playerService.Service, error) {
			return playerService.NewService(
				quizzlyConfig.Player.MustGet(),
				cookieService,
				log,
			), nil
		}),
		link: structs.NewSingleton(func() (link.Service, error) {
			return link.NewService(
				variables,
			), nil
		}),
	}

	settings := serverSettings{
		Port:         serverAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	mux := http.NewServeMux()
	muxExtended := &muxExtended{
		mux: mux,
		middleware: middleware.New(middleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		}),
	}

	// Serve static files
	_, err := os.Stat(publicPath)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory '%s' not found.\n", "web"))
	}
	mux.Handle("GET /files/public/",
		http.StripPrefix(
			"/files/public/",
			http.FileServer(http.Dir(publicPath)),
		),
	)

	// Serve S3 and other files
	mux.HandleFunc("GET /files/images/{image_name}", files2.NewGetImageHandler(filesManager, log).Handle())

	// Serve metrics with basic auth
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		providedUser, providedPassword, ok := r.BasicAuth()
		user := variables.GetString(variablesRepo.MetricsUser)
		password := variables.GetString(variablesRepo.MetricsPassword)

		if !ok || providedUser != user || providedPassword != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		promhttp.Handler().ServeHTTP(w, r)
	})

	adminRoutes(muxExtended, config, log, quizzlyConfig, authClient, filesManager)
	publicRoutes(muxExtended, log, config, quizzlyConfig, authClient)

	server := &http.Server{
		Addr:         settings.Port,
		Handler:      muxExtended.mux, // Implement your handlers function
		ReadTimeout:  settings.ReadTimeout,
		WriteTimeout: settings.WriteTimeout,
		IdleTimeout:  settings.IdleTimeout,
	}

	return &ServerInstance{
		serverLambda: httpadapter.New(muxExtended.mux).ProxyWithContext,
		serverHTTP:   server,
		serverType:   serverType,
		log:          log,
	}
}

func (s *ServerInstance) Start(ctx context.Context) {
	switch s.serverType {
	case ServerTypeLambda:
		s.log.Info("lambda server start")
		s.startLambda()
	default:
		s.log.Info("http server start")
		s.startHTTP(ctx)
	}
}

func (s *ServerInstance) startHTTP(ctx context.Context) {
	go func() {
		fmt.Println("Server listening on port ", s.serverHTTP.Addr)
		if err := s.serverHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("serverHTTP error", err)
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			s.log.Error("panic occurred", fmt.Errorf("%W", err))
		}
	}()

	<-ctx.Done()
}

func (s *ServerInstance) startLambda() {
	lambda.Start(s.serverLambda)
}

func (s *ServerInstance) Stop(ctx context.Context) {
	if s.serverType == ServerTypeLambda {
		return
	}

	s.log.Info("Shutting down serverHTTP...")

	// Create a context with a timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Shutdown the serverHTTP
	if err := s.serverHTTP.Shutdown(ctxWithTimeout); err != nil {
		s.log.Error("serverHTTP shutdown error", err)
	}
	s.log.Info("Server gracefully stopped")
}

package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"quizzly/internal/quizzly"
	"quizzly/pkg/auth"
	"quizzly/pkg/files"
	"quizzly/pkg/logger"
	"quizzly/pkg/structs"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/handlers/admin/game"
	"quizzly/web/frontend/handlers/admin/login"
	"quizzly/web/frontend/handlers/admin/question"
	"quizzly/web/frontend/handlers/admin/static/faq"
	files2 "quizzly/web/frontend/handlers/files"
	gamePublic "quizzly/web/frontend/handlers/public/game"
	playerService "quizzly/web/frontend/services/player"
	sessionService "quizzly/web/frontend/services/session"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
)

const (
	publicPath  = "web/frontend/public"
	serverAddr  = ":3000"
	metricsAddr = ":3333"
)

type (
	configuration struct {
		sessions structs.Singleton[sessionService.Service]
		player   structs.Singleton[playerService.Service]
	}

	serverSettings struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	muxWithMetrics struct {
		mux        *http.ServeMux
		middleware middleware.Middleware
	}
)

func (m *muxWithMetrics) HandleFunc(pattern string, metricsKey string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.Handle(pattern, middlewarestd.Handler(metricsKey, m.middleware, http.HandlerFunc(handler)))
}

func adminRoutes(
	mux *muxWithMetrics,
	config *configuration,
	log logger.Logger,
	quizzlyConfig *quizzly.Configuration,
	simpleAuth auth.SimpleAuth,
	filesManager files.Manager,
) {
	security := simpleAuth.Middleware("/admin/login")

	mux.HandleFunc("GET /admin/login", "/admin/login", handlers.Templ[struct{}](login.NewGetLoginPageHandler(), log))
	mux.HandleFunc("POST /admin/login", "/admin/login", handlers.Templ[login.PostLoginPageData](login.NewPostLoginPageHandler(simpleAuth), log))
	mux.HandleFunc("GET /admin/logout", "/admin/logout", handlers.Templ[struct{}](login.NewGetLogoutPageHandler(), log))

	mux.HandleFunc("GET /admin/question/new", "/admin/question/new", security.WithAuth(handlers.Templ[question.GetFormData](question.NewGetFormHandler(), log)))
	mux.HandleFunc("POST /admin/question", "/admin/question", security.WithAuth(handlers.Templ[question.NewPostData](question.NewPostCreateHandler(
		quizzlyConfig.Question.MustGet(),
		filesManager,
	), log)))
	mux.HandleFunc("DELETE /admin/question", "/admin/question", security.WithAuth(handlers.Templ[question.GetDeleteData](question.NewPostDeleteHandler(
		quizzlyConfig.Question.MustGet(),
	), log)))
	mux.HandleFunc("GET /admin/question/list", "/admin/question/list", security.WithAuth(handlers.Templ[question.GetListData](question.NewGetHandler(quizzlyConfig.Question.MustGet()), log)))

	mux.HandleFunc("GET /admin/game/new", "/admin/game/new", security.WithAuth(handlers.Templ[struct{}](game.NewGetFormHandler(), log)))
	mux.HandleFunc("POST /admin/game", "/admin/game", security.WithAuth(handlers.Templ[game.PostCreateData](game.NewPostCreateHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("GET /admin/game/{game_id}", "/admin/game/:game_id", security.WithAuth(handlers.Templ[game.GetAdminPageData](game.NewGetPageHandler(
		quizzlyConfig.Game.MustGet(),
		config.sessions.MustGet(),
	), log)))
	mux.HandleFunc("POST /admin/game/start", "/admin/game/start", security.WithAuth(handlers.Templ[game.PostStartData](game.NewPostStartHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("POST /admin/game/finish", "/admin/game/finish", security.WithAuth(handlers.Templ[game.PostFinishData](game.NewPostFinishHandler(quizzlyConfig.Game.MustGet()), log)))

	mux.HandleFunc("GET /admin/game/list", "/admin/game/list", security.WithAuth(handlers.Templ[struct{}](game.NewGetListHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("GET /admin/game/session/list", "/admin/game/session/list", security.WithAuth(handlers.Templ[game.GetSessionListData](game.NewGetSessionListHandler(config.sessions.MustGet()), log)))

	mux.HandleFunc("GET /admin/faq", "/admin/faq", security.WithAuth(handlers.Templ[struct{}](faq.NewStaticFAQHandler(), log)))
}

func publicRoutes(
	mux *muxWithMetrics,
	log logger.Logger,
	config *configuration,
	quizzlyConfig *quizzly.Configuration,
	simpleAuth auth.SimpleAuth,
) {
	security := simpleAuth.Middleware()

	gamePlayPageHandler := gamePublic.NewGetPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
		config.player.MustGet(),
	)
	gameRestartPageHandler := gamePublic.NewGetRestartPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		config.player.MustGet(),
	)
	gameResultsPagehandler := gamePublic.NewGetPlayResultsPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
		config.player.MustGet(),
	)
	gameRenamePlayerHandler := gamePublic.NewPostRenamePlayerHandler(
		quizzlyConfig.Player.MustGet(),
		config.player.MustGet(),
	)

	mux.HandleFunc("GET /", "/", security.WithEnrich(handlers.Templ[gamePublic.GetStartPageData](gamePublic.NewGetStartPageHandler(
		quizzlyConfig.Game.MustGet(),
	), log)))

	mux.HandleFunc("GET /game/{game_id}", "/game/:game_id", security.WithEnrich(handlers.Templ[gamePublic.GetPlayPageData](gamePlayPageHandler, log)))
	// backwards compatibility
	mux.HandleFunc("GET /game/play", "/game/:game_id (old)", security.WithEnrich(handlers.Templ[gamePublic.GetPlayPageData](gamePlayPageHandler, log)))

	mux.HandleFunc("POST /game/{game_id}", "/game/:game_id", handlers.Templ[gamePublic.PostPlayPageData](gamePublic.NewPostPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
		config.player.MustGet(),
	), log))

	mux.HandleFunc("GET /game/{game_id}/restart", "/game/:game_id/restart", security.WithEnrich(handlers.Templ[gamePublic.GetRestartPageData](gameRestartPageHandler, log)))
	// backwards compatibility
	mux.HandleFunc("GET /game/restart", "/game/:game_id/restart (old)", security.WithEnrich(handlers.Templ[gamePublic.GetRestartPageData](gameRestartPageHandler, log)))

	mux.HandleFunc("GET /game/{game_id}/results/{player_id}", "/game/:game_id/results/:player_id", security.WithEnrich(handlers.Templ[gamePublic.GetPlayResultsPageData](gameResultsPagehandler, log)))
	// backwards compatibility
	mux.HandleFunc("GET /game/results", "/game/:game_id/results/:player_id (old)", security.WithEnrich(handlers.Templ[gamePublic.GetPlayResultsPageData](gameResultsPagehandler, log)))

	mux.HandleFunc("POST /game/{game_id}/player/{player_id}/rename", "/game/:game_id/player/:player_id/rename", security.WithEnrich(handlers.Templ[gamePublic.PostRenamePlayerData](gameRenamePlayerHandler, log)))
}

func ServerRun(
	log logger.Logger,
	quizzlyConfig *quizzly.Configuration,
	simpleAuth auth.SimpleAuth,
	filesManager files.Manager,
) {
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
				log,
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
	muxWithMetrics := &muxWithMetrics{
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

	adminRoutes(muxWithMetrics, config, log, quizzlyConfig, simpleAuth, filesManager)
	publicRoutes(muxWithMetrics, log, config, quizzlyConfig, simpleAuth)

	server := &http.Server{
		Addr:         settings.Port,
		Handler:      muxWithMetrics.mux, // Implement your handlers function
		ReadTimeout:  settings.ReadTimeout,
		WriteTimeout: settings.WriteTimeout,
		IdleTimeout:  settings.IdleTimeout,
	}

	// Serve metrics.
	// Serve our metrics.
	go func() {
		fmt.Println("metrics listening at", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Error("error while serving metrics", err)
		}
	}()

	// ServerRun server in a goroutine
	go func() {
		fmt.Println("Server listening on port ", settings.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", err)
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Error("panic occurred", errors.New(fmt.Sprintf("%v", err)))
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Server gracefully stopped")
}

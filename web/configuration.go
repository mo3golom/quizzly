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
	files2 "quizzly/web/frontend/handlers/files"
	gamePublic "quizzly/web/frontend/handlers/public/game"
	questionService "quizzly/web/frontend/services/question"
	sessionService "quizzly/web/frontend/services/session"
	"syscall"
	"time"
)

const (
	publicPath = "web/frontend/public"
)

type (
	configuration struct {
		questions structs.Singleton[questionService.Service]
		sessions  structs.Singleton[sessionService.Service]
	}

	serverSettings struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
)

func routes(
	mux *http.ServeMux,
	config *configuration,
	log logger.Logger,
	quizzlyConfig *quizzly.Configuration,
	simpleAuth auth.SimpleAuth,
	filesManager files.Manager,
) {
	security := simpleAuth.Middleware()

	mux.HandleFunc("GET /file/{filename}", files2.NewGetFileHandler(filesManager, log).Handle())

	// ADMIN ROUTES
	mux.HandleFunc("GET /login", handlers.Templ[struct{}](login.NewGetLoginPageHandler(), log))
	mux.HandleFunc("POST /login", handlers.Templ[login.PostLoginPageData](login.NewPostLoginPageHandler(simpleAuth), log))
	mux.HandleFunc("GET /logout", handlers.Templ[struct{}](login.NewGetLogoutPageHandler(), log))

	mux.HandleFunc("GET /question/new", security.WithAuth(handlers.Templ[question.GetFormData](question.NewGetFormHandler(), log)))
	mux.HandleFunc("POST /question", security.WithAuth(handlers.Templ[question.NewPostData](question.NewPostCreateHandler(
		quizzlyConfig.Question.MustGet(),
		filesManager,
	), log)))
	mux.HandleFunc("DELETE /question", security.WithAuth(handlers.Templ[question.GetDeleteData](question.NewPostDeleteHandler(
		quizzlyConfig.Question.MustGet(),
	), log)))
	mux.HandleFunc("GET /question/list", security.WithAuth(handlers.Templ[struct{}](question.NewGetHandler(config.questions.MustGet()), log)))

	mux.HandleFunc("GET /game/new", security.WithAuth(handlers.Templ[struct{}](game.NewGetFormHandler(config.questions.MustGet()), log)))
	mux.HandleFunc("POST /game", security.WithAuth(handlers.Templ[game.PostCreateData](game.NewPostCreateHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("GET /game", security.WithAuth(handlers.Templ[game.GetAdminPageData](game.NewGetPageHandler(
		quizzlyConfig.Game.MustGet(),
		config.questions.MustGet(),
		config.sessions.MustGet(),
	), log)))
	mux.HandleFunc("POST /game/start", security.WithAuth(handlers.Templ[game.PostStartData](game.NewPostStartHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("POST /game/finish", security.WithAuth(handlers.Templ[game.PostFinishData](game.NewPostFinishHandler(quizzlyConfig.Game.MustGet()), log)))
	mux.HandleFunc("GET /game/list", security.WithAuth(handlers.Templ[struct{}](game.NewGetListHandler(quizzlyConfig.Game.MustGet()), log)))

	// PUBLIC ROUTES
	mux.HandleFunc("GET /game/play", handlers.Templ[gamePublic.GetPlayPageData](gamePublic.NewGetPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
	), log))
	mux.HandleFunc("POST /game/play", handlers.Templ[gamePublic.PostPlayPageData](gamePublic.NewPostPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
	), log))
}

func ServerRun(
	log logger.Logger,
	quizzlyConfig *quizzly.Configuration,
	simpleAuth auth.SimpleAuth,
	filesManager files.Manager,
) {
	config := &configuration{
		questions: structs.NewSingleton(func() (questionService.Service, error) {
			return questionService.NewService(quizzlyConfig.Question.MustGet()), nil
		}),
		sessions: structs.NewSingleton(func() (sessionService.Service, error) {
			return sessionService.NewService(
				quizzlyConfig.Session.MustGet(),
				quizzlyConfig.Player.MustGet(),
			), nil
		}),
	}

	settings := serverSettings{
		Port:         ":3000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	mux := http.NewServeMux()

	_, err := os.Stat(publicPath)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory '%s' not found.\n", "web"))
	}
	mux.Handle("/", http.FileServer(http.Dir(publicPath)))

	routes(mux, config, log, quizzlyConfig, simpleAuth, filesManager)

	server := &http.Server{
		Addr:         settings.Port,
		Handler:      mux, // Implement your handlers function
		ReadTimeout:  settings.ReadTimeout,
		WriteTimeout: settings.WriteTimeout,
		IdleTimeout:  settings.IdleTimeout,
	}

	// ServerRun server in a goroutine
	go func() {
		fmt.Printf("Server listening on port %s\n", settings.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error: %v\n", err)
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

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
	"quizzly/pkg/structs"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/handlers/admin/game"
	"quizzly/web/frontend/handlers/admin/login"
	"quizzly/web/frontend/handlers/admin/question"
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

func routes(mux *http.ServeMux, config *configuration, quizzlyConfig *quizzly.Configuration, simpleAuth auth.SimpleAuth) {
	security := simpleAuth.Middleware()

	// ADMIN ROUTES
	mux.HandleFunc("GET /login", handlers.Templ[struct{}](login.NewGetLoginPageHandler()))
	mux.HandleFunc("POST /login", handlers.Templ[login.PostLoginPageData](login.NewPostLoginPageHandler(simpleAuth)))

	mux.HandleFunc("GET /question/new", security.WithAuth(handlers.Templ[question.GetFormData](question.NewGetFormHandler())))
	mux.HandleFunc("POST /question", security.WithAuth(handlers.Redirect[question.NewPostData](question.NewPostCreateHandler(
		quizzlyConfig.Question.MustGet(),
		config.questions.MustGet(),
	))))
	mux.HandleFunc("GET /question/list", security.WithAuth(handlers.Templ[struct{}](question.NewGetHandler(config.questions.MustGet()))))

	mux.HandleFunc("GET /game/new", security.WithAuth(handlers.Templ[struct{}](game.NewGetFormHandler(config.questions.MustGet()))))
	mux.HandleFunc("POST /game", security.WithAuth(handlers.Redirect[game.PostCreateData](game.NewPostCreateHandler(quizzlyConfig.Game.MustGet()))))
	mux.HandleFunc("GET /game", security.WithAuth(handlers.Templ[game.GetAdminPageData](game.NewGetPageHandler(
		quizzlyConfig.Game.MustGet(),
		config.questions.MustGet(),
		config.sessions.MustGet(),
	))))
	mux.HandleFunc("POST /game/start", security.WithAuth(handlers.Templ[game.PostStartData](game.NewPostStartHandler(quizzlyConfig.Game.MustGet()))))
	mux.HandleFunc("POST /game/finish", security.WithAuth(handlers.Templ[game.PostFinishData](game.NewPostFinishHandler(quizzlyConfig.Game.MustGet()))))
	mux.HandleFunc("GET /game/list", security.WithAuth(handlers.Templ[struct{}](game.NewGetListHandler(quizzlyConfig.Game.MustGet()))))

	// PUBLIC ROUTES
	mux.HandleFunc("GET /game/play", handlers.Templ[gamePublic.GetPlayPageData](gamePublic.NewGetPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
	)))
	mux.HandleFunc("POST /game/play", handlers.Templ[gamePublic.PostPlayPageData](gamePublic.NewPostPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
	)))
}

func ServerRun(quizzlyConfig *quizzly.Configuration, simpleAuth auth.SimpleAuth) {
	config := &configuration{
		questions: structs.NewSingleton(func() (questionService.Service, error) {
			return questionService.NewService(quizzlyConfig.Question.MustGet()), nil
		}),
		sessions: structs.NewSingleton(func() (sessionService.Service, error) {
			return sessionService.NewService(quizzlyConfig.Session.MustGet()), nil
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

	routes(mux, config, quizzlyConfig, simpleAuth)

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
	signal.Notify(quit, os.Kill)
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

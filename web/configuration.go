package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"quizzly/internal/quizzly"
	"quizzly/pkg/structs"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/handlers/game/admin"
	"quizzly/web/frontend/handlers/game/public"
	"quizzly/web/frontend/handlers/question"
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

func routes(mux *http.ServeMux, config *configuration, quizzlyConfig *quizzly.Configuration) {
	mux.HandleFunc("GET /question/new", handlers.Wrapper[question.GetFormData](question.NewGetFormHandler()))
	mux.HandleFunc("POST /question", handlers.Wrapper[question.NewPostData](question.NewPostCreateHandler(
		quizzlyConfig.Question.MustGet(),
		config.questions.MustGet(),
	)))
	mux.HandleFunc("GET /question", handlers.Wrapper[struct{}](question.NewGetHandler(config.questions.MustGet())))

	mux.HandleFunc("GET /game/new", handlers.Wrapper[struct{}](admin.NewGetFormHandler(config.questions.MustGet())))
	mux.HandleFunc("POST /game", handlers.Wrapper[admin.PostCreateData](admin.NewPostCreateHandler(quizzlyConfig.Game.MustGet())))
	mux.HandleFunc("GET /game", handlers.Wrapper[admin.GetAdminPageData](admin.NewGetPageHandler(
		quizzlyConfig.Game.MustGet(),
		config.questions.MustGet(),
		config.sessions.MustGet(),
	)))
	mux.HandleFunc("POST /game/start", handlers.Wrapper[admin.PostStartData](admin.NewPostStartHandler(quizzlyConfig.Game.MustGet())))
	mux.HandleFunc("POST /game/finish", handlers.Wrapper[admin.PostFinishData](admin.NewPostFinishHandler(quizzlyConfig.Game.MustGet())))
	mux.HandleFunc("GET /game/play", handlers.Wrapper[public.GetPlayPageData](public.NewGetPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
		quizzlyConfig.Player.MustGet(),
	)))
	mux.HandleFunc("POST /game/play", handlers.Wrapper[public.PostPlayPageData](public.NewPostPlayPageHandler(
		quizzlyConfig.Game.MustGet(),
		quizzlyConfig.Session.MustGet(),
	)))
}

func ServerRun(quizzlyConfig *quizzly.Configuration) {
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

	routes(mux, config, quizzlyConfig)

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

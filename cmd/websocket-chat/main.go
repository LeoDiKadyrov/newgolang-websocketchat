package main

import (
	"log/slog"
	"net/http"
	"new-websocket-chat/internal/config"
	refresh "new-websocket-chat/internal/http_server/handlers/jwt"
	"new-websocket-chat/internal/http_server/handlers/user/delete"
	"new-websocket-chat/internal/http_server/handlers/user/save"
	mwLogger "new-websocket-chat/internal/http_server/middleware/logger"
	jwtAuth "new-websocket-chat/internal/lib/jwt"
	"new-websocket-chat/internal/lib/logger/sl"
	"new-websocket-chat/internal/storage/postgres"
	ws "new-websocket-chat/internal/websocket/handlers"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting websocket-chat", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.User, cfg.Password, cfg.DBname, cfg.Hostname, cfg.Port)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	jwtAuthService := &jwtAuth.JWTAuthService{} // Add error handling if not initialized

	hub := ws.NewHub()
	go hub.Run()

	log.Info("websocket hub was created", slog.Any("hub: ", hub))

	router.Post("/user", save.New(log, storage))
	router.Post("/api/jwt/refresh", refresh.New(log, jwtAuthService))
	router.Delete("/user/delete", delete.New(log, storage))
	router.HandleFunc("/ws", ws.ServeWs(log, hub))
	//router.Group(func(r chi.Router) {
	//	r.Use(jwtAuth.TokenAuthMiddleware)
	//	// r.Get("/Auth", auth.New(log, storage))
	//})

	// 1. JWT auth
	// 1.1 JWT unit testing (for generating token, validating token, etc.)
	// 2. Functional testing on saving and deleting
	// 3. Websocket messaging

	// middleware (цепочка хендлеров выполняется, есть основной и остальные, вроде обработки авторизации или модификации, должен быть middleware проверяющий авторизацию при изменении URLа)

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped") // we shouldn't reach this point
}

func NewHub() {
	panic("unimplemented")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

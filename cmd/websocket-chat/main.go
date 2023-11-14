package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"new-websocket-chat/internal/config"
	mwLogger "new-websocket-chat/internal/http_server/middleware/logger"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage: postgresql // sqlite ???
	//storage, err := postgres.New(cfg.User, cfg.Password, cfg.DBname, cfg.Hostname, cfg.Port)
	//if err != nil {
	//	log.Error("failed to init storage", sl.Err(err))
	//	os.Exit(1)
	//}
	//_ = storage
	// TODO: init router: chi, "chi render"

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// middleware (цепочка хендлеров выполняется, есть основной и остальные, вроде обработки авторизации или модификации, должен быть middleware проверяющий авторизацию при изменении URLа)

	// TODO: init server:

	fmt.Println("main is operating")
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

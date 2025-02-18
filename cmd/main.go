package main

import (
	"log/slog"
	"os"

	"github.com/qwaq-dev/golnag-archive/internal/config"
	"github.com/qwaq-dev/golnag-archive/internal/routes"
	"github.com/qwaq-dev/golnag-archive/internal/structures/server"
	"github.com/qwaq-dev/golnag-archive/pkg/logger/handlers/slogpretty"
	"github.com/qwaq-dev/golnag-archive/pkg/logger/sl"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	router := routes.NewRouter(log)

	srv := &server.Server{}

	log.Info("Server starts on port"+cfg.Port, slog.String("env", cfg.Env))

	if err := srv.Run(cfg.Port, router); err != nil {
		log.Error("Error with starting server", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

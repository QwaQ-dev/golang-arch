package routes

import (
	"log/slog"

	"github.com/go-chi/chi"
	"github.com/qwaq-dev/golnag-archive/internal/handlers/middleware/logger"
	"github.com/qwaq-dev/golnag-archive/internal/handlers/upload"
)

func NewRouter(log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.New(log))
	router.Post("/upload", upload.UploadFileHandler(log))

	return router
}

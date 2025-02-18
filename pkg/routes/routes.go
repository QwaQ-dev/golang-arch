package routes

import (
	"log/slog"

	"github.com/go-chi/chi"
	"github.com/qwaq-dev/golnag-archive/pkg/handlers/middleware/logger"
	"github.com/qwaq-dev/golnag-archive/pkg/handlers/upload"
)

func NewRouter(log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.New(log))
	router.Post("/upload", upload.UploadFile(log))

	return router
}

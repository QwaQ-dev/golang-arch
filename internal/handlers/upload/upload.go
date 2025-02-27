package upload

import (
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/qwaq-dev/golnag-archive/internal/service/comp"
	"github.com/qwaq-dev/golnag-archive/internal/service/upload"
	"github.com/qwaq-dev/golnag-archive/pkg/logger/sl"
)

const (
	formDataKey = "file"
)

func UploadFileHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePaths, err := upload.UploadFileFromRequest(log, w, r, formDataKey)
		if err != nil {
			log.Error("Error with form", sl.Err(err))
		}
		log.Info("Files upload successfully")

		photoExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true,
		}

		videoExts := map[string]bool{
			".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".wmv": true, ".flv": true, ".webm": true,
		}

		for _, path := range filePaths {
			ext := filepath.Ext(path)

			if photoExts[ext] {
				log.Info("Compressing image")
				comp.CompressPhoto(path, log)
			} else if videoExts[ext] {
				log.Info("Compressing video")
				comp.CompressVideo(path, log)
			} else {
				log.Error("Unknown file type", slog.String("file", path))
			}
		}
	}

}

package upload

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/qwaq-dev/golnag-archive/internal/service/archive"
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

		compressedPaths := []string{}

		for _, path := range filePaths {
			ext := filepath.Ext(path)

			if photoExts[ext] {
				log.Info("Compressing image")
				outputPath, err := comp.CompressPhoto(path, log)
				if err == nil {
					removeFile(path, log)
				}
				compressedPaths = append(compressedPaths, outputPath)
			} else if videoExts[ext] {
				log.Info("Compressing video")
				outputPathVideo, err := comp.CompressVideo(path, log)
				if err == nil {
					removeFile(path, log)
				}
				compressedPaths = append(compressedPaths, outputPathVideo)
			} else {
				log.Error("Unknown file type", slog.String("file", path))
			}
		}

		zipPath, err := archive.CreateZipArchive(compressedPaths, log)
		if err != nil {
			log.Error("Error creating zip archive", sl.Err(err))
			http.Error(w, "Failed to create ZIP archive", http.StatusInternalServerError)
			return
		}

		sendFile(w, zipPath)
	}
}

func removeFile(filePath string, log *slog.Logger) {
	if err := os.Remove(filePath); err != nil {
		log.Error("Failed to delete file", sl.Err(err))
	} else {
		log.Info("File was delete successfully", slog.String("file", filePath))
	}
}

func sendFile(w http.ResponseWriter, filePath string) {

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/zip")

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	io.Copy(w, file)

}

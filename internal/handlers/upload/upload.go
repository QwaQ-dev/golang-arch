package upload

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"

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
			http.Error(w, "Failed to upload files", http.StatusInternalServerError)
			return
		}
		log.Info("Files uploaded successfully")

		photoExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
		videoExts := map[string]bool{".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".wmv": true, ".flv": true, ".webm": true}

		var wg sync.WaitGroup
		compressedPaths := make(chan string, len(filePaths)) // Канал для хранения путей сжатых файлов

		for _, path := range filePaths {
			ext := filepath.Ext(path)

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()

				var outputPath string
				var err error

				if photoExts[ext] {
					log.Info("Compressing image", slog.String("file", filePath))
					outputPath, err = comp.CompressPhoto(filePath, log)
				} else if videoExts[ext] {
					log.Info("Compressing video", slog.String("file", filePath))
					outputPath, err = comp.CompressVideo(filePath, log)
				} else {
					log.Error("Unknown file type", slog.String("file", filePath))
					return
				}

				if err == nil {
					removeFile(filePath, log)
					compressedPaths <- outputPath
				}
			}(path)
		}

		wg.Wait()
		close(compressedPaths)

		var finalPaths []string
		for path := range compressedPaths {
			finalPaths = append(finalPaths, path)
		}

		zipPath, err := archive.CreateZipArchive(finalPaths, log)
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

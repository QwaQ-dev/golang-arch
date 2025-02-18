package upload

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/qwaq-dev/golnag-archive/pkg/lib/logger/sl"
)

const formDataKey = "file"

func UploadFile(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(0) // Ограничение на 10MB для файла
		if err != nil {
			http.Error(w, "Cannot parse form", http.StatusBadRequest)
			log.Error("Cannot parse form", sl.Err(err))
			return
		}

		files := r.MultipartForm.File[formDataKey]
		if len(files) == 0 {
			http.Error(w, "No files uploaded", http.StatusBadRequest)
			log.Error("No files uploaded")
			return
		}

		if err := os.MkdirAll("../uploads", os.ModePerm); err != nil {
			log.Error("Cannot create directory", sl.Err(err))
		}

		log.Info("Uploading files", slog.Int("files lenght", len(files)))

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Cannot open file", http.StatusInternalServerError)
				log.Error("Cannot open file", sl.Err(err))
				return
			}

			dstFileName := generateUniqueFilename(fileHeader.Filename)
			dstPath := filepath.Join("../uploads", dstFileName)

			dst, err := os.Create(dstPath)
			if err != nil {
				http.Error(w, "Error with saving file", http.StatusInternalServerError)
				log.Error("Error with saving file", sl.Err(err))
				return
			}

			_, err = io.Copy(dst, file)
			dst.Close()
			file.Close()
			if err != nil {
				http.Error(w, "Error with copy", http.StatusInternalServerError)
				log.Error("Error with copy", sl.Err(err))
				return
			}

			log.Info("File was upload successfully", slog.String("file-name", fileHeader.Filename))
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Files upload successfully")
	}
}

func generateUniqueFilename(filename string) string {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

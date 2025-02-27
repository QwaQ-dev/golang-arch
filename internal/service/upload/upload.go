package upload

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/qwaq-dev/golnag-archive/pkg/generatename"
	"github.com/qwaq-dev/golnag-archive/pkg/logger/sl"
)

func UploadFileFromRequest(log *slog.Logger, w http.ResponseWriter, r *http.Request, formDataKey string) ([]string, error) {
	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		log.Error("Cannot parse form", sl.Err(err))
		return nil, err
	}

	files := r.MultipartForm.File[formDataKey]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		log.Error("No files uploaded")
		return nil, err
	}

	if err := os.MkdirAll("../uploads", os.ModePerm); err != nil {
		log.Error("Cannot create directory", sl.Err(err))
		return nil, err
	}

	log.Info("Uploading files", slog.Int("files lenght", len(files)))

	var filePath []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Cannot open file", http.StatusInternalServerError)
			log.Error("Cannot open file", sl.Err(err))
			continue
		}

		dstFileName := generatename.GenerateUniqueFilename(fileHeader.Filename)
		dstPath := filepath.Join("../uploads", dstFileName)

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Error with saving file", http.StatusInternalServerError)
			log.Error("Error with saving file", sl.Err(err))
			continue
		}

		_, err = io.Copy(dst, file)
		dst.Close()
		file.Close()
		if err != nil {
			http.Error(w, "Error with copy", http.StatusInternalServerError)
			log.Error("Error with copy", sl.Err(err))
			continue
		}

		absPath, err := filepath.Abs(dstPath)

		if err != nil {
			log.Error("Error getting absolute path", sl.Err(err))
			continue
		}

		filePath = append(filePath, absPath)

		log.Info("File was upload successfully", slog.String("file-path", dstPath))
	}

	return filePath, nil
}

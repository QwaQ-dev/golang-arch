package archive

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/qwaq-dev/golnag-archive/pkg/logger/sl"
)

func CreateZipArchive(files []string, log *slog.Logger) (string, error) {
	zipFileName := "compressed_files.zip"
	zipFile, err := os.Create(zipFileName)

	if err != nil {
		log.Error("Error with creating Zip archive", sl.Err(err))
		return "", nil
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			log.Error("Cannot open file")
			return "", err
		}
		defer file.Close()

		w, err := zipWriter.Create(filepath.Base(filePath))
		if err != nil {
			return "", err
		}

		_, err = io.Copy(w, file)
		if err != nil {
			return "", err
		}
	}

	return zipFileName, nil
}

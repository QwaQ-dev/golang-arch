package comp

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/qwaq-dev/golnag-archive/pkg/logger/sl"
)

func CompressVideo(inputPath string, log *slog.Logger) (string, error) {
	if !fileExists(inputPath) {
		log.Error("File is not exists")
	}
	dir := filepath.Dir(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	ext := filepath.Ext(inputPath)

	output := filepath.Join(dir, baseName+"_compressed"+ext)
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx264", "-crf", "28", output)

	err := cmd.Run()
	if err != nil {
		log.Error("Error during compression", slog.String("file", inputPath), slog.String("error", err.Error()))
		return "", err
	}

	log.Info("Video compressed successfully", slog.String("output-file", output))
	return output, nil
}

func CompressPhoto(inputPath string, log *slog.Logger) (string, error) {
	if !fileExists(inputPath) {
		log.Error("File does not exist", slog.String("file", inputPath))
		return "", fmt.Errorf("file does not exist: %s", inputPath)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	path, err := compressImageFile(inputPath, ext, log)

	if err != nil {
		return "", fmt.Errorf("Cannot compress file")
	}

	return path, nil
}

func compressImageFile(inputPath, ext string, log *slog.Logger) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		log.Error("Error opening file", sl.Err(err))
		return "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		log.Error("Error decoding image", sl.Err(err))
		return "", err
	}

	var buffer bytes.Buffer

	switch format {
	case "jpeg":
		err = jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 75})
	case "png":
		err = png.Encode(&buffer, img)
	case "jpg":
		err = png.Encode(&buffer, img)
	default:
		log.Error("Unsupported encoding format", slog.String("format", format))
		return "", fmt.Errorf("unsupported encoding format: %s", format)
	}

	if err != nil {
		log.Error("Error encoding image", sl.Err(err))
		return "", err
	}
	outputPath := strings.TrimSuffix(inputPath, ext) + "_compressed" + ext
	err = os.WriteFile(outputPath, buffer.Bytes(), 0644)
	if err != nil {
		log.Error("Error saving compressed image", sl.Err(err))
		return "", err
	}

	log.Info("Image compressed successfully", slog.String("output-file", outputPath))
	return outputPath, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

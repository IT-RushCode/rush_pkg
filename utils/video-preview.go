package utils

import (
	"log"

	"github.com/mowshon/moviego"
)

// Функция для создания превью из видео
func CreateVideoPreview(videoPath, previewPath string, previewTime float64) (string, error) {
	movie, err := moviego.Load(videoPath)
	if err != nil {
		return "", err
	}

	fileName, err := movie.Screenshot(previewTime, previewPath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fileName, nil
}

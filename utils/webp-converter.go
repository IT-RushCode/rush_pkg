package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/chai2010/webp"
)

// ConvertToWebP конвертирует изображение в формат WebP.
func ConvertToWebP(inputPath, outputPath string, lossless bool, quality int) error {
	// Открываем исходный файл изображения
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Определяем формат изображения
	var img image.Image
	var format string

	if strings.HasSuffix(inputPath, ".jpg") || strings.HasSuffix(inputPath, ".jpeg") {
		img, err = jpeg.Decode(file)
		format = "jpeg"
	} else if strings.HasSuffix(inputPath, ".png") {
		img, err = png.Decode(file)
		format = "png"
	} else {
		img, format, err = image.Decode(file)
	}

	if err != nil {
		return fmt.Errorf("не удалось декодировать изображение: %v", err)
	}

	log.Printf("Формат изображения: %s\n", format)

	// Создаем буфер для хранения данных WebP
	var buf bytes.Buffer

	// Кодируем изображение в формат WebP
	if quality == 0 {
		quality = 85 // Установим качество на 85 для уменьшения размера
	}
	options := &webp.Options{Lossless: lossless, Quality: float32(quality)}
	if err := webp.Encode(&buf, img, options); err != nil {
		return fmt.Errorf("не удалось закодировать изображение в WebP: %v", err)
	}

	// Записываем результат в файл
	if err := os.WriteFile(outputPath, buf.Bytes(), 0666); err != nil {
		return fmt.Errorf("не удалось сохранить файл: %v", err)
	}

	log.Printf("Изображение успешно сохранено в %s\n", outputPath)
	return nil
}

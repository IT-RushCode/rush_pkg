package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Константы ошибок
var (
	ErrFileType       = errors.New("необходимо указать один из допустимых типов (например, PDF, видео, изображения)")
	ErrFileMaxSize    = errors.New("файл слишком большой. Максимальный размер 5 МБ")
	ErrFileTypeFormat = errors.New("недопустимый тип файла (разрешены только JPG, JPEG, PNG, WEBP и другие)")
	ErrFileConvert    = errors.New("ошибка конвертации файла в WebP")
	ErrUnprocessable  = errors.New("не удалось обработать файл")
)

// UploadFile - универсальный метод для загрузки файлов
func UploadFile(fileHeader *multipart.FileHeader, options *FileUploadOptions) (string, string, error) {
	// Установка максимального размера по умолчанию, если не указан
	if options.MaxSize == 0 {
		options.MaxSize = 10 * 1024 * 1024 // 10 MB по умолчанию
	}

	// Проверка размера файла
	if fileHeader.Size > int64(options.MaxSize) {
		return "", "", ErrFileMaxSize // Ошибка при превышении размера файла
	}

	// Проверка расширения файла
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !options.AllowedTypes[ext] {
		return "", "", ErrFileTypeFormat // Ошибка при неверном типе файла
	}

	// Создание директории, если она отсутствует
	err := os.MkdirAll(options.BaseDir, os.ModePerm)
	if err != nil {
		return "", "", ErrCreateDir // Ошибка создания директории
	}

	// Генерация уникального имени файла (UUID)
	fileName := uuid.New().String()
	originalFilePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s%s", fileName, ext))

	// Открытие файла для чтения
	src, err := fileHeader.Open()
	if err != nil {
		return "", "", ErrUnprocessable // Ошибка открытия файла
	}
	defer src.Close()

	// Создание файла на диске
	dest, err := os.Create(originalFilePath)
	if err != nil {
		return "", "", ErrSaveFile // Ошибка сохранения файла
	}
	defer dest.Close()

	// Копирование содержимого файла
	if _, err := dest.ReadFrom(src); err != nil {
		return "", "", ErrSaveFile // Ошибка копирования содержимого
	}

	// Если это изображение, конвертируем его в WebP
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		webpFilePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s.webp", fileName))
		err = ConvertToWebP(originalFilePath, webpFilePath, options.Lossless, options.Quality)
		if err != nil {
			return "", "", fmt.Errorf("%w: %s", ErrFileConvert, err.Error()) // Ошибка конвертации файла
		}

		// Удаление оригинального файла после конвертации
		os.Remove(originalFilePath)

		// Возвращаем путь к WebP и имя файла
		return webpFilePath, fileName, nil
	}

	// Возвращаем путь к файлу и имя файла для других типов файлов
	return originalFilePath, fileName, nil
}

// Обертка для загрузки файла через контекст fiber
func UploadFileFromCtx(ctx *fiber.Ctx, formFieldName string, options *FileUploadOptions) (string, string, error) {
	// Получение файла из контекста
	file, err := ctx.FormFile(formFieldName)
	if file == nil {
		return "", "", err // Ошибка загрузки файла
	}

	// Используем основную функцию для загрузки файла
	return UploadFile(file, options)
}

// FileUploadOptions - структура для хранения параметров загрузки файлов
type FileUploadOptions struct {
	MaxSize      int             // Максимальный размер файла (в байтах)
	AllowedTypes map[string]bool // Допустимые типы файлов (расширения)
	BaseDir      string          // Базовая директория для сохранения файлов
	Lossless     bool            // Использовать ли сжатие без потерь (для изображений)
	Quality      int             // Качество сжатия (для изображений)
}

// DefaultImageOptions - параметры по умолчанию для изображений
func DefaultImageOptions() *FileUploadOptions {
	return &FileUploadOptions{
		MaxSize: 5 * 1024 * 1024, // 5 MB
		AllowedTypes: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
		},
		Lossless: false,
		Quality:  85,
	}
}

// DefaultVideoOptions - параметры по умолчанию для видео
func DefaultVideoOptions() *FileUploadOptions {
	return &FileUploadOptions{
		MaxSize: 100 * 1024 * 1024, // 100 MB
		AllowedTypes: map[string]bool{
			".mp4": true,
			".mov": true,
			".avi": true,
			".mkv": true,
		},
	}
}

// DefaultPDFOptions - параметры по умолчанию для PDF
func DefaultPDFOptions() *FileUploadOptions {
	return &FileUploadOptions{
		MaxSize: 10 * 1024 * 1024, // 10 MB
		AllowedTypes: map[string]bool{
			".pdf": true,
		},
	}
}

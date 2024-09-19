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
	ErrFileMaxSize    = errors.New("файл слишком большой. Максимальный размер: ")
	ErrFileTypeFormat = errors.New("недопустимый тип файла (разрешены только PDF, JPG, JPEG, PNG, WEBP и другие)")
	ErrFileConvert    = errors.New("ошибка конвертации файла в WebP")
	ErrUnprocessable  = errors.New("не удалось обработать файл")
)

// UploadFile - функция для обработки уже загруженного файла (например, конвертация)
func UploadFile(savedFilePath string, options *FileUploadOptions) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(savedFilePath))
	fileName := strings.TrimSuffix(filepath.Base(savedFilePath), ext)

	// Если это изображение, конвертируем его в WebP
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		webpFilePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s.webp", fileName))
		err := ConvertToWebP(savedFilePath, webpFilePath, options.Lossless, options.Quality)
		if err != nil {
			return "", "", fmt.Errorf("%w: %s", ErrFileConvert, err.Error()) // Ошибка конвертации файла
		}

		// Удаление оригинального файла после конвертации
		os.Remove(savedFilePath)

		// Возвращаем путь к WebP и имя файла
		return webpFilePath, fileName, nil
	}

	// Возвращаем путь к файлу и имя файла для других типов файлов
	return savedFilePath, fileName, nil
}

// UploadFileFromCtx - функция для сохранения файла на диск через Fiber контекст
func UploadFileFromCtx(ctx *fiber.Ctx, formFieldName string, options *FileUploadOptions) (string, string, error) {
	// Получение файла из контекста
	file, err := ctx.FormFile(formFieldName)
	if err == fiber.ErrUnprocessableEntity || file == nil {
		return "", "", nil // Возвращаем nil, если файла нет
	}

	// Генерация уникального имени файла (UUID)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileName := uuid.New().String()
	filePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s%s", fileName, ext))

	// Сохраняем файл на диск
	if err := ctx.SaveFile(file, filePath); err != nil {
		return "", "", ErrSaveFile // Ошибка сохранения файла
	}

	// Возвращаем путь к файлу и имя файла
	return filePath, fileName, nil
}

// UploadFile - универсальный метод для загрузки файлов
func OldUploadFile(fileHeader *multipart.FileHeader, options *FileUploadOptions) (string, string, error) {
	// Установка максимального размера по умолчанию, если не указан
	if options.MaxSize == 0 {
		options.MaxSize = 10 * 1024 * 1024 // 10 MB по умолчанию
	}

	// Проверка размера файла
	if fileHeader.Size > int64(options.MaxSize) {
		return "", "", fmt.Errorf("%s%d%s", ErrFileMaxSize, options.MaxSize, " МБ") // Ошибка при превышении размера файла
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
func OldUploadFileFromCtx(ctx *fiber.Ctx, formFieldName string, options *FileUploadOptions) (string, string, error) {
	// Получение файла из контекста
	file, err := ctx.FormFile(formFieldName)
	if err == fiber.ErrUnprocessableEntity || file == nil {
		return "", "", nil // Возвращаем nil, если файла нет
	}
	// Используем основную функцию для загрузки файла
	return OldUploadFile(file, options)
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
func DefaultFileOptions() *FileUploadOptions {
	return &FileUploadOptions{
		MaxSize: 10 * 1024 * 1024, // 10 MB
		AllowedTypes: map[string]bool{
			".pdf":  true,
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

package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Константы ошибок
var (
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
	uuid := uuid.New().String()
	filePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s%s", uuid, ext))

	// Сохраняем файл на диск
	if err := ctx.SaveFile(file, filePath); err != nil {
		return "", "", ErrSaveFile // Ошибка сохранения файла
	}

	// Возвращаем путь к файлу и имя файла
	return filePath, uuid, nil
}

// UploadFile - универсальный метод для загрузки файлов с проверкой MIME-типа
func OldUploadFile(fileHeader *multipart.FileHeader, options *FileUploadOptions) (string, string, error) {
	// Установка максимального размера по умолчанию, если не указан
	if options.MaxSize == 0 {
		options.MaxSize = 10 * 1024 * 1024 // 10 MB по умолчанию
	}

	// Проверка размера файла
	if fileHeader.Size > int64(options.MaxSize) {
		return "", "", fmt.Errorf("превышен максимальный размер файла: %s (допустимо не более %s)",
			formatFileSize(fileHeader.Size), formatFileSize(int64(options.MaxSize))) // Ошибка при превышении размера файла
	}

	// Открытие файла для чтения
	src, err := fileHeader.Open()
	if err != nil {
		return "", "", ErrUnprocessable // Ошибка открытия файла
	}
	defer src.Close()

	// Чтение первых 512 байт файла для определения MIME-типа
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return "", "", ErrUnprocessable // Ошибка чтения файла
	}

	// Определяем MIME-тип файла
	mimeType := http.DetectContentType(buffer)

	// Проверка соответствия MIME-типа разрешённым типам
	if !options.AllowedMimeTypes[mimeType] {
		return "", "", errors.New(generateAllowedFileTypesError(options.AllowedMimeTypes, options.AllowedTypes))
	}

	// Возвращаемся к началу файла для последующего копирования
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", "", ErrUnprocessable // Ошибка перемещения указателя в начало файла
	}

	// Создание директории, если она отсутствует
	err = os.MkdirAll(options.BaseDir, os.ModePerm)
	if err != nil {
		return "", "", ErrCreateDir // Ошибка создания директории
	}

	// Генерация уникального имени файла (UUID)
	uuid := uuid.New().String()
	ext := filepath.Ext(fileHeader.Filename) // Получение расширения из имени файла, если оно есть
	if ext == "" {
		// Если нет расширения, можно попытаться определить его по MIME-типу
		ext = GetExtensionFromMimeType(mimeType)
	}
	originalFilePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s%s", uuid, ext))

	// Создание файла на диске
	dest, err := os.Create(originalFilePath)
	if err != nil {
		return "", "", ErrSaveFile // Ошибка сохранения файла
	}
	defer dest.Close()

	// Копирование содержимого файла
	if _, err := io.Copy(dest, src); err != nil {
		return "", "", ErrSaveFile // Ошибка копирования содержимого
	}

	// Если это изображение, конвертируем его в WebP
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		webpFilePath := filepath.Join(options.BaseDir, fmt.Sprintf("%s.webp", uuid))
		err = ConvertToWebP(originalFilePath, webpFilePath, options.Lossless, options.Quality)
		if err != nil {
			return "", "", fmt.Errorf("%w: %s", ErrFileConvert, err.Error()) // Ошибка конвертации файла
		}

		// Удаление оригинального файла после конвертации
		os.Remove(originalFilePath)

		// Возвращаем путь к WebP и имя файла
		return webpFilePath, uuid, nil
	}

	// Возвращаем путь к файлу и имя файла для других типов файлов
	return originalFilePath, uuid, nil
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
	MaxSize          int             // Максимальный размер файла (в байтах)
	AllowedTypes     map[string]bool // Допустимые типы файлов (расширения)
	AllowedMimeTypes map[string]bool // Допустимые MIME-типы файлов
	BaseDir          string          // Базовая директория для сохранения файлов
	Lossless         bool            // Использовать ли сжатие без потерь (для изображений)
	Quality          int             // Качество сжатия (для изображений)
}

// DefaultFileOptions - параметры по умолчанию для изображений
func DefaultFileOptions() *FileUploadOptions {
	return &FileUploadOptions{
		MaxSize: 10 * 1024 * 1024, // 10 MB
		AllowedTypes: map[string]bool{
			".pdf":  true,
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
		},
		AllowedMimeTypes: map[string]bool{
			"application/pdf": true,
			"image/jpeg":      true,
			"image/png":       true,
			"image/gif":       true,
			"image/webp":      true,
		},
		Lossless: false,
		Quality:  85,
	}
}

// DefaultVideoOptions - параметры по умолчанию для видео
func DefaultVideoOptions() *FileUploadOptions {
	return &FileUploadOptions{
		MaxSize: 1024 * 1024 * 1024, // 1000 GB
		AllowedTypes: map[string]bool{
			".mp4": true,
			".mov": true,
			".avi": true,
			".mkv": true,
		},
		AllowedMimeTypes: map[string]bool{
			"video/mp4":        true,
			"video/quicktime":  true,
			"video/x-msvideo":  true,
			"video/x-matroska": true,
			"video/webm":       true,
		},
	}
}

// Формируем сообщение с разрешёнными типами файлов
func generateAllowedFileTypesError(allowedMimeTypes map[string]bool, allowedTypes map[string]bool) string {
	var mimeTypes []string
	var extensions []string

	// Собираем MIME-типы
	for mime := range allowedMimeTypes {
		mimeTypes = append(mimeTypes, mime)
	}

	// Собираем расширения файлов
	for ext := range allowedTypes {
		extensions = append(extensions, ext)
	}

	return fmt.Sprintf("недопустимый тип файла (разрешены только: MIME-типы: %s, расширения: %s)",
		strings.Join(mimeTypes, ", "), strings.Join(extensions, ", "))
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1024*1024*1024:
		return fmt.Sprintf("%.2f ГБ", float64(size)/(1024*1024*1024))
	case size >= 1024*1024:
		return fmt.Sprintf("%.2f МБ", float64(size)/(1024*1024))
	case size >= 1024:
		return fmt.Sprintf("%.2f КБ", float64(size)/1024)
	default:
		return fmt.Sprintf("%d байт", size)
	}
}

package utils

// Функция для получения расширения файла по MIME-типу
func GetExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	// Изображения
	case "image/jpeg":
		return ".jpg"
	case "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"

	// Документы
	case "application/pdf":
		return ".pdf"

	// Видео форматы
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	case "video/x-msvideo":
		return ".avi"
	case "video/x-matroska":
		return ".mkv"
	case "video/webm":
		return ".webm"

	// Аудио форматы (если необходимо)
	case "audio/mpeg":
		return ".mp3"
	case "audio/wav":
		return ".wav"

	default:
		return "" // Если MIME-тип не распознан, возвращаем пустую строку
	}
}
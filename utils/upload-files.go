package utils

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FileData struct {
	UUID     string
	FileName string
	FilePath string
}

// UploadFile принимает путь к папке и входящее название из FormData и сохраняет файл в соответвующей папке.
// Если папка не существует, то создается новая.
func UploadFile(ctx *fiber.Ctx, pathFile, formName string) (*FileData, error) {
	file, err := ctx.FormFile(formName)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil
		}
		return nil, SendResponse(ctx, false, "ошибка получения файла", nil, http.StatusBadRequest)
	}

	err = os.MkdirAll("./uploads/"+pathFile, os.ModePerm)
	if err != nil {
		return nil, SendResponse(ctx, false, "ошибка создания директории", nil, http.StatusInternalServerError)
	}

	filePath := filepath.Join("./uploads/"+pathFile, file.Filename)
	err = ctx.SaveFile(file, filePath)
	if err != nil {
		return nil, SendResponse(ctx, false, "ошибка сохранения файла", nil, http.StatusInternalServerError)
	}

	return &FileData{
		UUID:     uuid.New().String(),
		FileName: file.Filename,
		FilePath: filePath,
	}, err
}

func GetFile(ctx *fiber.Ctx, fileUrl string) (string, error) {
	if _, err := os.Stat(fileUrl); os.IsNotExist(err) {
		return "", SendResponse(ctx, false, "файл не найден", nil, http.StatusNotFound)
	}

	return "", nil
}

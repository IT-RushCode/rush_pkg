package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// LoadEntitiesFromFile загружает данные из указанного файла и записывает их в указанный срез сущностей.
// T - тип сущности.
// filename - имя файла, содержащего JSON данные.
// entities - ссылка на срез сущностей, куда будут записаны данные.
//
// // Пример использования:
//
//	var entities []Entity
//
//	LoadEntitiesFromFile("entities.json", &entities)
func LoadEntitiesFromJsonFile[T any](filename string, entities *[]T) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(bytes, entities); err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
	}
}

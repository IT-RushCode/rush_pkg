package utils

import (
	"fmt"
	"os"
)

func CheckFlags(migrate, seed *bool, migrateFunc func(), seedFunc func()) {
	if *migrate {
		migrateFunc()
		fmt.Println("Успешная миграция базы")
		os.Exit(0)
	}

	if *seed {
		seedFunc()
		fmt.Println("Успешная загрузка сидов")
		os.Exit(0)
	}
}

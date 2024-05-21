package database

import (
	"fmt"
	"log"

	"gitlab.arvand.tj/conveyor/arvand_pkg/config"
	"gitlab.arvand.tj/conveyor/arvand_pkg/utils"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func MSSQL_CONNECT(cfg *config.DatabaseConfig) *gorm.DB {
	conn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := gorm.Open(sqlserver.Open(conn), &gorm.Config{})
	if err != nil && (cfg.Title == "MSSQL_CONVEYOR" || cfg.Title == "MSSQL_ABS") {
		log.Printf(
			"%sБАЗА %s НЕДОСТУПНА! \n%sПРОВЕРЬТЕ СОЕДИНЕНИЕ/НАСТРОЙКИ ИЛИ ЗАКОМЕНТИРУЙТЕ %s\"database.MSSQL_CONNECT(&cfg.DB.%s)\"%s.",
			utils.Red, cfg.Title, utils.Cyan, utils.Green, cfg.Title, utils.White,
		)
		return nil
	} else if err != nil {
		panic("database is not connected!")
	}

	return db
}

func MSSQL_CLOSE(db *gorm.DB) (error, bool) {
	conn, err := db.DB()
	if err != nil {
		return err, false
	}

	return nil, conn.Close() != nil
}

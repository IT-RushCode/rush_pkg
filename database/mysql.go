package database

import (
	"fmt"

	"gitlab.arvand.tj/conveyor/arvand_pkg/config"

	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MYSQL_CONNECT(cfg *config.DatabaseConfig) *gorm.DB {
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.CHARSET,
	)

	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN:                       dns,
				DefaultStringSize:         256,
				DisableDatetimePrecision:  true,
				DontSupportRenameIndex:    true,
				DontSupportRenameColumn:   true,
				SkipInitializeWithVersion: false,
			},
		), &gorm.Config{})

	if err != nil {
		panic("database is not connected!")
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic("failed to get underlying SQL DB")
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	return db
}

func MYSQL_CLOSE(db *gorm.DB) (error, bool) {
	conn, err := db.DB()

	if err != nil {
		return err, false
	}

	return nil, conn.Close() != nil
}

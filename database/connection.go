package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver string `json:"driver"`

	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`

	Environment string `json:"environment"`
}

var Db *gorm.DB

func Init(config2 Config, migrate bool, entitys []interface{}) {

	var err error

	dblogger := log.New(log.Writer(), "[DB] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	dblogger.Printf("initialization started")

	var dbLogger logger.Interface

	if config2.Environment != "development" {
		dbLogger = logger.New(
			dblogger,
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,          // Disable color
			},
		)
	} else {
		dbLogger = logger.Default
	}

	dialectConfig := &gorm.Config{
		Logger: dbLogger,
	}

	if config2.Driver == "sqlite" {
		dblogger.Printf("sqlite db path hardcoded to ./local.db")
		Db, err = gorm.Open(sqlite.Open("./local.db"), dialectConfig)
		if err != nil {
			panic("unable to start db conn: " + err.Error())
		}
	} else {

		pgDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Kiev", config2.Host, config2.Username, config2.Password, config2.Database, config2.Port)
		Db, err = gorm.Open(postgres.Open(pgDsn), dialectConfig)
	}

	if err != nil {
		panic(err)
	}

	if migrate {
		for _, it := range entitys {
			Db.AutoMigrate(it)
		}

	} else {
		dblogger.Printf("migration were skipped. use --migrate arg to perform migration along startup")
	}

	dblogger.Printf("db initialization finished")
}

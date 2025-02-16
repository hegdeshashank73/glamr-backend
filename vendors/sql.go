package vendors

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hegdeshashank73/glamr-backend/utils"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DBMono *sql.DB
var DBConfig *dbConfig

type dbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func loadDBConfig() dbConfig {
	var config dbConfig

	config.Host = viper.GetString("DATABASE.HOST")
	port, _ := strconv.Atoi(viper.GetString("DATABASE.PORT"))
	config.Port = port
	config.Username = viper.GetString("DATABASE.USER")
	config.Password = viper.GetString("DATABASE.PASSWORD")

	return config
}

func initDB() {
	st := time.Now()
	defer utils.LogTimeTaken("init.initDB", st)

	config := loadDBConfig()
	DBConfig = &config
	DBConfig.Database = viper.GetString("DATABASE.DATABASE")

	dataSourceName := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=require",
		config.Username,
		config.Password,
		config.Host,
		DBConfig.Database,
	)

	var err error
	DBMono, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Caught Error While Connecting to Postgres DB: %v", err)
	}
}

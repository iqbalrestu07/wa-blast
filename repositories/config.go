package repositories

import (
	"fmt"
	"os"

	"wa-blast/configs"
	"wa-blast/loggers"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = loggers.Get()

// DBConnect ...
var DBConnect *gorm.DB

// Config ...
func Config() {
	// Init error codes
	DBConn()
	Ping()
}

// DBConn ...
func DBConn() {
	host := configs.MustGetString("database.master.host")
	user := configs.MustGetString("database.master.username")
	pswd := configs.MustGetString("database.master.password")
	dbnm := configs.MustGetString("database.master.database")

	dsn := user + ":" + pswd + "@" + host + "/" + dbnm + "?charset=utf8mb4&parseTime=true&loc=Local"
	fmt.Println("dsn", dsn)
	// db, _ := sql.Open("sqlite", "file:wa.db?cache=shared")

	// db, _ := sql.Open("sqlite3", "file:wa.db?cache=shared")
	db, err := gorm.Open(sqlite.Open("wa.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("errorrr", err)
		return
	}

	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.User{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.MasterType{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.MasterCategory{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.MasterVendor{}).AddForeignKey("master_type_id", "boma.master_type(id)", "RESTRICT", "RESTRICT").AddForeignKey("user_id", "boma.user(id)", "RESTRICT", "RESTRICT")
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.UserRole{})
	// db.SingularTable(true)
	DBConnect = db
}

// Ping ...
func Ping() {
	pingDB, _ := DBConnect.DB()
	ping := pingDB.Ping()
	if ping != nil {
		log.Info("Failed Connecting Mysql.")
		os.Exit(1)
	} else {
		log.Info("Success Connecting Mysql.")
	}
}

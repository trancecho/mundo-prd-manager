package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-prd-manager/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func DBInit() {
	dsn := fmt.Sprintf(
		"root:%s@tcp(%s:%s)/dev?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s&readTimeout=30s",
		viper.GetString("mysql.pwd"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
	)
	log.Println("dsn:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB = db
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}

	err = AutoMigrate(DB)
	if err != nil {
		log.Fatalln("failed to migrate database", err)
	}
	log.Println("connect database success!!!!!!!!")
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.ChatIDHistory{},
	)
	return err
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatalln("database not initialized")
	}
	return DB
}

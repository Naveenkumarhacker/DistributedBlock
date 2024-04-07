package database

import (
	"DistributedBlock/constants"
	"DistributedBlock/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitMySqlDB() *gorm.DB {

	dsn := constants.DbUser + ":" + constants.DbPassword + "@tcp" + "(" + constants.DbHost + ":" + constants.DbPort + ")/" + constants.DbName + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error while connecting to mysql database : error=%v \n", err)
		return nil
	}

	// Auto Migrate the Block struct
	err = db.AutoMigrate(&models.Block{})
	if err != nil {
		log.Fatalf("Error while auto migrate to mysql database : error=%v \n", err)
		return nil
	}
	return db
}

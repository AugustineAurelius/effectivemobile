package repository

import (
	"effectivemobile/FIO"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Init() *gorm.DB {
	db, err := gorm.Open(postgres.Open("postgres://pgadmin:5432@localhost:5432/effictiveMobile"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	err = db.AutoMigrate(&FIO.FIO{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectionDatabase() {
	dsn := os.Getenv("DATABASE_DEV")

	if dsn == "" {
		log.Fatal("❌ DATABASE_DEV env variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL database")

	// Assign to global DB variable
	DB = db

	// Optionally, test the connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get DB instance: %v", err)
	}

	// Test ping
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping DB: %v", err)
	}

	// Optionally set connection pool config
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
}

// type DbInstance struct {
// 	Db *gorm.DB
// }

// var DB DbInstance

// func ConnectionDatabase() {
// 	dsn := os.Getenv("DATABASE_DEV")

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	// DB = DB.Set("gorm:auto_preload", true)

// 	if err != nil {
// 		log.Println("Failed to connect to database: ", err)
// 	}
// 	db.Logger = logger.Default.LogMode(logger.Info)

// 	// DB = DB.Set("gorm:auto_preload", true) //Setting auto preload untuk relasi

// 	DB = DbInstance{
// 		Db: db,
// 	}
// }

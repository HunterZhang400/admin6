package database

import (
	"admin6/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDB *gorm.DB

// GetGormDB returns the GORM database instance
func GetGormDB() *gorm.DB {
	return gormDB
}

// InitGormDB initializes the GORM database connection
func InitGormDB(database *config.Database) {
	var err error

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		database.User,
		database.Password,
		database.Host,
		database.Port,
		database.Name)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	gormDB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("🚀 GORM Connected Successfully to the Database")

}

// createPerformanceIndexes creates database indexes for better performance
func createPerformanceIndexes(db *gorm.DB) error {
	//indexes := []string{
	//	"CREATE INDEX idx_user_id ON user(id)",
	//	"CREATE INDEX idx_user_score ON user(score)",
	//	"CREATE INDEX idx_user_remain ON user(remain)",
	//	"CREATE INDEX idx_user_in_challenge ON user(in_challenge)",
	//	"CREATE INDEX idx_user_result_id ON user(result_id)",
	//	"CREATE INDEX idx_spin_result_user_id ON spin_result(user_id)",
	//	"CREATE INDEX idx_spin_result_user_id_luck ON spin_result(user_id, luck_result_1)",
	//	"CREATE INDEX idx_spin_result_created_at ON spin_result(created_at)",
	//	"CREATE INDEX idx_spin_result_user_created ON spin_result(user_id, created_at)",
	//	"CREATE INDEX idx_user_score_remain ON user(score, remain)",
	//	"CREATE INDEX idx_user_challenge_remain ON user(in_challenge, remain)",
	//}
	//
	//successCount := 0
	//for _, indexSQL := range indexes {
	//	if err := db.Exec(indexSQL).Error; err != nil {
	//		// Check if it's a "Duplicate key name" error (index already exists)
	//		if strings.Contains(err.Error(), "Duplicate key name") || strings.Contains(err.Error(), "already exists") {
	//			log.Printf("Index already exists, skipping: %s", indexSQL)
	//			successCount++
	//		} else {
	//			log.Printf("Failed to create index: %s, error: %v", indexSQL, err)
	//		}
	//	} else {
	//		successCount++
	//	}
	//}
	//
	//log.Printf("✅ Performance indexes: %d/%d created successfully", successCount, len(indexes))
	return nil
}

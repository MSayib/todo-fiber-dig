package database

import (
	"fmt"
	"log"

	"github.com/msayib/todo-fiber-dig/internal/config" // Ganti dengan path modul Anda
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	log.Println("Koneksi ke Postgres (GORM) berhasil!")
	return db, nil
}
package xiangqin_backend

import (
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type ListenConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type PostgresConfig struct {
	DSN             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

type Logger struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

type Config struct {
	Listen   ListenConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Logger   Logger         `yaml:"logger"`
}

func GetConfig(path string) (Config, error) {
	log.Printf("load config: %s", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func NewPostgres(cfg Config) (*gorm.DB, error) {
	log.Printf("connect postgres: %s", cfg.Postgres.DSN)
	db, err := gorm.Open(
		postgres.Open(cfg.Postgres.DSN),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Postgres.ConnMaxIdleTime) * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Postgres.ConnMaxLifetime) * time.Minute)

	return db, nil
}

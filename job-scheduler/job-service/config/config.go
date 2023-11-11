package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	App struct {
		Name         string        `envconfig:"APP_NAME" default:"job-service"`
		StartTimeout time.Duration `envconfig:"START_TIMEOUT" default:"1m"`
		StopTimeout  time.Duration `envconfig:"STOP_TIMEOUT" default:"1m"`
		JobShardSize uint64        `envconfig:"JOB_SHARD_SIZE" default:"10"`
	}
	Http struct {
		Port         int           `envconfig:"HTTP_PORT" default:"8080"`
		ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"1m"`
		WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"1m"`
	}
	Log struct {
		Level string `envconfig:"LOG_LEVEL" default:"debug"`
	}
	Postgresql struct {
		Host                                     string `envconfig:"POSTGRESQL_HOST" default:"localhost"`
		Port                                     string `envconfig:"POSTGRESQL_PORT" default:"5432"`
		Username                                 string `envconfig:"POSTGRESQL_USERNAME" default:"postgres"`
		Password                                 string `envconfig:"POSTGRESQL_PASSWORD" default:"postgres"`
		DbName                                   string `envconfig:"POSTGRESQL_DBNAME" default:"job-scheduler"`
		SSLMode                                  string `envconfig:"POSTGRESQL_SSL_MODE" default:"disable"`
		MaxIdleConns                             int    `envconfig:"POSTGRESQL_MAX_IDLE_CONNS" default:"10"`
		MaxOpenConns                             int    `envconfig:"POSTGRESQL_MAX_OPEN_CONNS" default:"50"`
		DisableForeignKeyConstraintWhenMigrating bool   `envconfig:"POSTGRES_DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING" default:"true"`
		PrepareStmt                              bool   `envconfig:"POSTGRES_PREPARE_STMT" default:"false"`
		ConnMaxLifeTimeInSecond                  int    `envconfig:"POSTGRES_CONN_MAX_LIFE_TIME_IN_SECOND" default:"600"`
		ConnMaxIdleTimeInSecond                  int    `envconfig:"POSTGRES_CONN_MAX_IDLE_TIME_IN_SECOND" default:"30"`
		DebugQuery                               bool   `envconfig:"POSTGRES_DEBUG_QUERY" default:"true"`
	}
	Jwt struct {
		SecretKey           string        `envconfig:"JWT_SECRET_KEY" default:"secret"`
		AccessTokenTimeOut  time.Duration `envconfig:"JWT_ACCESS_TOKEN_TIMEOUT" default:"1h"`
		RefreshTokenTimeOut time.Duration `envconfig:"JWT_REFRESH_TOKEN_TIMEOUT" default:"4h"`
	}
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	var config Config

	err := envconfig.Process("job-service", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewConfig() *Config {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Error occurred when load config", zap.Error(err))
	}

	return config
}

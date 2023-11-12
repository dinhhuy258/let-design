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
		Name         string        `envconfig:"APP_NAME" default:"job-scheduler-worker"`
		StartTimeout time.Duration `envconfig:"START_TIMEOUT" default:"1m"`
		StopTimeout  time.Duration `envconfig:"STOP_TIMEOUT" default:"1m"`
		ShardIds     []uint64      `envconfig:"SHARD_IDS" default:"1,2,3,4,5,6,7,8,9,10"`
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
	Kafka struct {
		BootstrapServers  string `envconfig:"KAFKA_BOOTSTRAP_SERVERS" default:"kafka:9091"`
		ClientId          string `envconfig:"KAFKA_PRODUCER_CLIENT_ID" default:"job-scheduler-worker"`
		ACKS              string `envconfig:"KAFKA_PRODUCER_ACKS" default:"all"`
		ScheduledJobTopic string `envconfig:"KAFKA_PRODUCER_SCHEDULED_JOB_TOPIC" default:"scheduled_jobs"`
	}
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	var config Config

	err := envconfig.Process("job-scheduler-worker", &config)
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

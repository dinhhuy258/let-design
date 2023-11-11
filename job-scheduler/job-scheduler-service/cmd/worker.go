package cmd

import (
	"context"
	"fmt"
	"job-scheduler-service/config"
	"job-scheduler-service/internal/adapter/worker"
	"job-scheduler-service/internal/infra/kafka"
	"job-scheduler-service/internal/infra/postgresql"
	"job-scheduler-service/internal/usecase"
	"job-scheduler-service/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var createServiceCommand = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Start worker",
		Long:  "Start worker",
		Run: func(_ *cobra.Command, _ []string) {
			runSevice()
		},
	}

	return cmd
}

func runSevice() {
	appCtx, canFunc := context.WithCancel(context.Background())
	conf := config.NewConfig()
	logger := logger.New(conf.Log.Level)
	app := fx.New(
		fx.StartTimeout(conf.App.StartTimeout), fx.StopTimeout(conf.App.StopTimeout),
		fx.Provide(
			worker.NewSchedulerWorker,
			newDatabaseConnection,
			postgresql.NewJobRepository,
			kafka.NewMessageBusRepository,
			usecase.NewJobUsecase,
			usecase.NewFairSchedulerUsecase,
		),
		fx.Supply(conf, logger),
		fx.Invoke(startWorker),
		fx.Decorate(),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	stopChan := make(chan os.Signal)

	err := app.Start(appCtx)
	if err != nil {
		log.Fatalf("error occurred when start app %+v", err)
	}

	go func() {
		val := <-quit

		logger.Info("stopping app")

		err := app.Stop(appCtx)
		if err != nil {
			logger.Error("error occurred when stop app %v", err)
		}

		canFunc()
		stopChan <- val
	}()

	<-stopChan
}

func startWorker(lc fx.Lifecycle,
	logger *logger.Logger,
	worker worker.SchedulerWorker,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Start scheduler worker")

			worker.Start()

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Stop scheduler worker")

			worker.Stop()

			return nil
		},
	})
}

func newDatabaseConnection(lc fx.Lifecycle, conf *config.Config, logger *logger.Logger) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		conf.Postgresql.Host,
		conf.Postgresql.Port,
		conf.Postgresql.Username,
		conf.Postgresql.DbName,
		conf.Postgresql.Password,
	)

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: connectionString,
			}),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: conf.Postgresql.DisableForeignKeyConstraintWhenMigrating,
			PrepareStmt:                              conf.Postgresql.PrepareStmt,
		},
	)
	if err != nil {
		return nil, err
	}

	database, err := db.DB()
	if err != nil {
		return nil, err
	}

	database.SetMaxIdleConns(conf.Postgresql.MaxIdleConns)
	database.SetMaxOpenConns(conf.Postgresql.MaxOpenConns)
	database.SetConnMaxLifetime(time.Duration(conf.Postgresql.ConnMaxLifeTimeInSecond) * time.Second)
	database.SetConnMaxIdleTime(time.Duration(conf.Postgresql.ConnMaxIdleTimeInSecond) * time.Second)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			database, err := db.DB()
			if err != nil {
				return err
			}

			return database.Ping()
		},
		OnStop: func(_ context.Context) error {
			database, err := db.DB()
			if err != nil {
				return err
			}

			logger.Info("Closing database connection")

			return database.Close()
		},
	})

	return db, nil
}

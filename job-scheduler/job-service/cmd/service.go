package cmd

import (
	"context"
	"errors"
	"fmt"
	"job-service/config"
	httpv1 "job-service/internal/adapter/http/v1"
	"job-service/internal/infra/postgresql"
	"job-service/internal/usecase"
	"job-service/pkg/httpserver"
	"job-service/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresMigration "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for gorm driver
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var createServiceCommand = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Start service",
		Long:  "Start service",
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
			newDatabaseConnection,
			newHttpServer,
			postgresql.NewUserRepository,
			usecase.NewUserUsecase,
			usecase.NewAuthUsecase,
			httpv1.NewAuthController,
			httpv1.NewUserController,
		),
		fx.Supply(conf, logger),
		fx.Invoke(startServer),
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

func startServer(lc fx.Lifecycle,
	conf *config.Config,
	logger *logger.Logger,
	server httpserver.Interface,
	authController httpv1.AuthController,
	userController httpv1.UserController,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Http server is listening on %v", conf.Http.Port)

			httpv1.SetRoutes(server, authController, userController)

			server.Start(ctx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Http server is shutting down")

			return server.Stop(ctx)
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

			err = database.Ping()
			if err != nil {
				return err
			}

			// migrate database
			driver, err := postgresMigration.WithInstance(database, &postgresMigration.Config{})
			if err != nil {
				return err
			}

			m, err := migrate.NewWithDatabaseInstance("file://./migration", conf.Postgresql.DbName, driver)
			if err != nil {
				return err
			}

			logger.Info("Migrating database...")

			err = m.Up()
			if err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					return nil
				}

				logger.Error("Unable to migrate database. Error: %v", err)

				return err
			}

			return nil
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

func newHttpServer(_ fx.Lifecycle, conf *config.Config) httpserver.Interface {
	return httpserver.New(conf.Http.Port, conf.Http.ReadTimeout, conf.Http.WriteTimeout)
}

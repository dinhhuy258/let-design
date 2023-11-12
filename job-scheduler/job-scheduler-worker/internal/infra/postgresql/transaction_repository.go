package postgresql

import (
	"context"

	"gorm.io/gorm"
)

type gormTransactionKey struct{}

type transactionRepositoryInterface interface {
	WithTransaction(context.Context, func(context.Context) error) error
	GetTransactionOrCreate(context.Context) *gorm.DB
}

type transactionRepository struct {
	database *gorm.DB
}

func (_self *transactionRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return _self.database.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, gormTransactionKey{}, tx)

		return fn(ctx)
	})
}

func (_self *transactionRepository) GetTransactionOrCreate(ctx context.Context) *gorm.DB {
	txValue := ctx.Value(gormTransactionKey{})
	switch txSession := txValue.(type) {
	case *gorm.DB:
		return txSession
	default:
		return _self.database
	}
}

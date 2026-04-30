package repository

import (
	"context"

	"myapp/internal/repository/generated"
)

type DB struct{}

type ProductRepository struct {
	db *DB
	*generated.ProductRepository
}

func executorFromContext(ctx context.Context, db *DB) generated.Executor {
	return nil
}

func translateError(err error) error { return err }

// GetGood wraps the error through translateError.
func (r *ProductRepository) GetGood(ctx context.Context, id string) (string, error) {
	row, err := r.GetProductByID(ctx, executorFromContext(ctx, r.db), id)
	if err != nil {
		return "", translateError(err)
	}
	return row, nil
}

// GetBad returns the raw err.
func (r *ProductRepository) GetBad(ctx context.Context, id string) (string, error) {
	row, err := r.GetProductByID(ctx, executorFromContext(ctx, r.db), id)
	if err != nil {
		return "", err // want `wrap err through translateError`
	}
	return row, nil
}

// NoGenCall returns err but doesn't touch the generated layer; the analyzer
// should leave it alone.
func (r *ProductRepository) NoGenCall(input error) error {
	return input
}

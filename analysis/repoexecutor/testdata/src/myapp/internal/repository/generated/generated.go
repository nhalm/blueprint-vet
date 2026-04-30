// Package generated stubs the skimatik-generated repository surface.
package generated

import "context"

type Executor interface{}

type ProductRepository struct{}

func (r *ProductRepository) GetProductByID(ctx context.Context, exec Executor, id string) (string, error) {
	return "", nil
}

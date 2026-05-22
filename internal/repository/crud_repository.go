package repository

import "context"

type CrudRepository[T any] interface {
	Create(ctx context.Context, item T) (*T, error)
	Get(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, item T) (*T, error)
	Delete(ctx context.Context, id string) error
}

package domain

import (
	"context"
	"errors"

	"github.com/cgroschupp/go-mail-admin/internal/model"
)

var (
	ErrDomainExists    = errors.New("domain already exists")
	ErrDomainNotExists = errors.New("domain not exists")
)

type DomainService interface {
	Create(ctx context.Context, name string) (model.Domain, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, id int32, name string) (model.Domain, error)
	List(ctx context.Context) ([]model.Domain, error)
	Get(ctx context.Context, id int32) (model.Domain, error)
	Stats(ctx context.Context) (Stats, error)
}

package domain

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/model"
)

type AliasService interface {
	Create(ctx context.Context, sourceUsername *string, sourceDomainID int32, destinationUsername, destinationDomain string) (model.Alias, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, id int32, sourceUsername *string, sourceDomainID *int32, destinationUsername, destinationDomain *string, enabled *bool) (model.Alias, error)
	List(ctx context.Context) ([]model.Alias, error)
	Get(ctx context.Context, id int32) (model.Alias, error)
	Stats(ctx context.Context) (Stats, error)
}

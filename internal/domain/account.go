package domain

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/model"
)

type AccountService interface {
	Create(ctx context.Context, username, password string, quota int32, sendonly, enabled bool, domainId int) (model.Account, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, id int32, username *string, quota *int32, sendonly, enabled *bool) (model.Account, error)
	List(ctx context.Context) ([]model.Account, error)
	Get(ctx context.Context, id int32) (model.Account, error)
	ChangePassword(ctx context.Context, id int32, password string) error
	Stats(ctx context.Context) (Stats, error)
}

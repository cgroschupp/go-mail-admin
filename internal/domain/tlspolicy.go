package domain

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/model"
)

type TLSPolicyService interface {
	Create(ctx context.Context, param *string, policy string, domainId int) (model.TLSPolicy, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, id int32, param, policy *string) (model.TLSPolicy, error)
	List(ctx context.Context) ([]model.TLSPolicy, error)
	Get(ctx context.Context, id int32) (model.TLSPolicy, error)
}

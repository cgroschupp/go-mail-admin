package service

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"gorm.io/gorm"
)

type tlsPolicyService struct {
	db *gorm.DB
}

// Create implements domain.TLSPolicyService.
func (t *tlsPolicyService) Create(ctx context.Context, param *string, policy string, domainId int) (model.TLSPolicy, error) {
	tp := model.TLSPolicy{DomainID: domainId, Policy: policy, Params: param}
	if err := t.db.Save(&tp).Error; err != nil {
		return tp, err
	}
	return t.Get(ctx, int32(tp.ID))
}

// Delete implements domain.TLSPolicyService.
func (t *tlsPolicyService) Delete(ctx context.Context, id int32) error {
	return t.db.Delete(&model.TLSPolicy{}, id).Error
}

// Get implements domain.TLSPolicyService.
func (t *tlsPolicyService) Get(ctx context.Context, id int32) (model.TLSPolicy, error) {
	tp := model.TLSPolicy{}
	if err := t.db.Preload("Domain").First(&tp, id).Error; err != nil {
		return tp, nil
	}
	return tp, nil
}

// List implements domain.TLSPolicyService.
func (t *tlsPolicyService) List(ctx context.Context) ([]model.TLSPolicy, error) {
	tlspolicies := []model.TLSPolicy{}
	if err := t.db.Preload("Domain").Find(&tlspolicies).Error; err != nil {
		return tlspolicies, err
	}
	return tlspolicies, nil
}

// Update implements domain.TLSPolicyService.
func (t *tlsPolicyService) Update(ctx context.Context, id int32, param, policy *string) (model.TLSPolicy, error) {
	tp, err := t.Get(ctx, id)
	if err != nil {
		return tp, err
	}
	if param != nil {
		tp.Params = param
	}

	if policy != nil {
		tp.Policy = *policy
	}
	return tp, t.db.Save(&tp).Error
}

func NewTLSPolicyService(db *gorm.DB) domain.TLSPolicyService {
	return &tlsPolicyService{db: db}
}

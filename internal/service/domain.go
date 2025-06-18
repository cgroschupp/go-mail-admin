package service

import (
	"context"
	"errors"

	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"gorm.io/gorm"
)

type domainService struct {
	db *gorm.DB
}

// Stats implements domain.DomainService.
func (d *domainService) Stats(ctx context.Context) (domain.Stats, error) {
	type result struct {
		Total int32
	}
	results := result{}
	if err := d.db.Model(&model.Domain{}).Select("count(*) as total").Find(&results).Error; err != nil {
		return domain.Stats{}, err
	}

	return domain.Stats{
		Labels:   []string{"Domains"},
		Datasets: []domain.Dataset{{Data: []int32{results.Total}, BackgroundColor: []string{"green"}}},
	}, nil
}

// Get implements domain.DomainService.
func (d *domainService) Get(ctx context.Context, id int32) (model.Domain, error) {
	do := model.Domain{}
	if err := d.db.First(&do, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = domain.ErrDomainNotExists
		}
		return do, err
	}
	return do, nil
}

// List implements domain.DomainService.
func (d *domainService) List(ctx context.Context) ([]model.Domain, error) {
	domains := []model.Domain{}
	if err := d.db.Find(&domains).Error; err != nil {
		return domains, err
	}
	return domains, nil
}

// Create implements domain.DomainService.
func (d *domainService) Create(ctx context.Context, name string) (model.Domain, error) {
	do := model.Domain{Name: name}
	if err := d.db.Create(&do).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = domain.ErrDomainExists
		}
		return do, err
	}
	return d.Get(ctx, int32(do.ID))
}

// Delete implements domain.DomainService.
func (d *domainService) Delete(ctx context.Context, name int32) error {
	return d.db.Delete(&model.Domain{}, name).Error
}

// Update implements domain.DomainService.
func (d *domainService) Update(ctx context.Context, id int32, name string) (model.Domain, error) {
	do := model.Domain{Model: model.Model{ID: uint(id)}, Name: name}
	return do, d.db.Save(&do).Error
}

func NewDomainService(db *gorm.DB) domain.DomainService {
	return &domainService{db: db}
}

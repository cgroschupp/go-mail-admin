package service

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"gorm.io/gorm"
)

type aliasService struct {
	db *gorm.DB
}

func NewAliasService(db *gorm.DB) domain.AliasService {
	return &aliasService{db: db}
}

// Create implements domain.AliasService.
func (a *aliasService) Create(ctx context.Context, sourceUsername *string, sourceDomainID int32, destinationUsername, destinationDomain string) (model.Alias, error) {
	alias := model.Alias{
		SourceDomainID:      uint(sourceDomainID),
		DestinationUsername: destinationUsername,
		DestinationDomain:   destinationDomain,
		Enabled:             true,
	}

	if sourceUsername != nil {
		alias.SourceUsername = *sourceUsername
	}

	if err := a.db.Create(&alias).Error; err != nil {
		return alias, err
	}
	return alias, nil
}

// Delete implements domain.AliasService.
func (a *aliasService) Delete(ctx context.Context, id int32) error {
	if err := a.db.Delete(&model.Alias{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Get implements domain.AliasService.
func (a *aliasService) Get(ctx context.Context, id int32) (model.Alias, error) {
	alias := model.Alias{}
	if err := a.db.Preload("SourceDomain").First(&alias, id).Error; err != nil {
		return alias, err
	}
	return alias, nil
}

// List implements domain.AliasService.
func (a *aliasService) List(ctx context.Context) ([]model.Alias, error) {
	aliases := []model.Alias{}
	if err := a.db.Preload("SourceDomain").Find(&aliases).Error; err != nil {
		return aliases, err
	}
	return aliases, nil
}

// Update implements domain.AliasService.
func (a *aliasService) Update(ctx context.Context, id int32, sourceUsername *string, sourceDomainID *int32, destinationUsername, destinationDomain *string, enabled *bool) (model.Alias, error) {
	al, err := a.Get(ctx, id)
	if err != nil {
		return al, err
	}
	if sourceUsername != nil {
		al.SourceUsername = *sourceUsername
	}
	if sourceDomainID != nil {
		al.SourceDomainID = uint(*sourceDomainID)
	}
	if destinationUsername != nil {
		al.DestinationUsername = *destinationUsername
	}
	if destinationDomain != nil {
		al.DestinationDomain = *destinationDomain
	}
	if enabled != nil {
		al.Enabled = *enabled
	}
	return al, a.db.Save(&al).Error
}

// Stats implements domain.AliasService.
func (a *aliasService) Stats(ctx context.Context) (domain.Stats, error) {
	var enabled, disabled int32
	if err := a.db.Model(&model.Alias{}).Select(`
		SUM(CASE WHEN enabled THEN 1 ELSE 0 END)  AS enabled,
		SUM(CASE WHEN NOT enabled THEN 1 ELSE 0 END) AS disabled
	`).Row().Scan(&enabled, &disabled); err != nil {
		return domain.Stats{}, err
	}
	result := []int32{disabled, enabled}
	stats := domain.Stats{
		Labels:   []string{"Disabled", "Enabled"},
		Datasets: []domain.Dataset{{Data: result, BackgroundColor: []string{"red", "green"}}},
	}
	return stats, nil
}

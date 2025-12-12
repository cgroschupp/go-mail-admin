package service

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/password"
	"gorm.io/gorm"
)

type accountService struct {
	db                  *gorm.DB
	passwordHashBuilder password.PasswordHashBuilder
}

// ChangePassword implements domain.AccountService.
func (a *accountService) ChangePassword(ctx context.Context, id int32, password string) error {
	pwHash, err := a.passwordHashBuilder.Hash(password)
	if err != nil {
		return err
	}
	return a.db.Model(&model.Account{}).Where("id = ?", id).Update("password", pwHash).Error
}

func NewAccountService(db *gorm.DB, passwordHashType string) domain.AccountService {
	return &accountService{
		db: db, passwordHashBuilder: password.GetPasswordHashBuilder(passwordHashType),
	}
}

// Create implements domain.AccountService.
func (a *accountService) Create(ctx context.Context, username, password string, quota int32, sendonly, enabled bool, domainId int) (model.Account, error) {
	acc := model.Account{Username: username, DomainID: uint(domainId), Quota: quota, Enabled: enabled, SendOnly: sendonly}
	pwHash, err := a.passwordHashBuilder.Hash(password)
	if err != nil {
		return acc, err
	}
	acc.Password = pwHash

	if err := a.db.Save(&acc).Error; err != nil {
		return acc, err
	}

	return a.Get(ctx, int32(acc.ID))
}

// Delete implements domain.AccountService.
func (a *accountService) Delete(ctx context.Context, id int32) error {
	if err := a.db.Delete(&model.Account{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Get implements domain.AccountService.
func (a *accountService) Get(ctx context.Context, id int32) (model.Account, error) {
	account := model.Account{}
	if err := a.db.Preload("Domain").First(&account, id).Error; err != nil {
		return account, err
	}
	return account, nil
}

// List implements domain.AccountService.
func (a *accountService) List(ctx context.Context) ([]model.Account, error) {
	accounts := []model.Account{}
	if err := a.db.Preload("Domain").Find(&accounts).Error; err != nil {
		return accounts, err
	}
	return accounts, nil
}

// Update implements domain.AccountService.
func (a *accountService) Update(ctx context.Context, id int32, username *string, quota *int32, sendonly, enabled *bool) (model.Account, error) {
	account, err := a.Get(ctx, id)
	if err != nil {
		return account, err
	}

	if username != nil {
		account.Username = *username
	}
	if quota != nil {
		account.Quota = *quota
	}
	if sendonly != nil {
		account.SendOnly = *sendonly
	}
	if enabled != nil {
		account.Enabled = *enabled
	}
	return account, a.db.Save(&account).Error
}

// Stats implements domain.AccountService.
func (a *accountService) Stats(ctx context.Context) (domain.Stats, error) {
	var enabled, disabled int32
	if err := a.db.Model(&model.Account{}).Select(`
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

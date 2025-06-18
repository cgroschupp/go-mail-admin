package service

import (
	"context"

	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/cgroschupp/go-mail-admin/internal/version"
	"gorm.io/gorm"
)

type dashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) domain.DashboardService {
	return &dashboardService{db: db}
}

func (d dashboardService) Version(ctx context.Context) string {
	return version.Version
}

// Status implements domain.DashboardService.
func (d *dashboardService) Status(ctx context.Context) bool {
	if err := d.db.Select("1").Error; err != nil {
		return false
	}
	return true
}

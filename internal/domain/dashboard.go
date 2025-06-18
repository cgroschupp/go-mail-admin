package domain

import "context"

type DashboardService interface {
	Version(ctx context.Context) string
	Status(ctx context.Context) bool
}

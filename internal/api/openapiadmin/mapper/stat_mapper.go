package mapper

import (
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/domain"
)

func MapStatsToResponse(req domain.Stats) openapiadmin.DashboardStatsItem {
	datasets := []openapiadmin.DashboardDataset{}
	for _, v := range req.Datasets {
		datasets = append(datasets, openapiadmin.DashboardDataset{BackgroundColor: v.BackgroundColor, Data: v.Data})
	}
	return openapiadmin.DashboardStatsItem{Labels: req.Labels, Datasets: datasets}
}

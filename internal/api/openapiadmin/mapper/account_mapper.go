package mapper

import (
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/utils/ptr"
)

func MapAccountToResponse(req model.Account) openapiadmin.Account {
	return openapiadmin.Account{
		Domain:    MapDomainToResponse(*req.Domain),
		CreatedAt: &req.CreatedAt,
		UpdatedAt: &req.UpdatedAt,
		DomainId:  int32(req.DomainID),
		Enabled:   &req.Enabled,
		Id:        ptr.To(int32(req.ID)),
		Quota:     &req.Quota,
		Sendonly:  &req.SendOnly,
		Username:  req.Username,
	}
}

func MapAccountListToResponse(req []model.Account) openapiadmin.AccountList {
	rtn := openapiadmin.AccountList{Items: []openapiadmin.Account{}}
	for _, item := range req {
		rtn.Items = append(rtn.Items, MapAccountToResponse(item))
	}
	return rtn
}

package mapper

import (
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/utils/ptr"
)

func MapDomainToResponse(req model.Domain) *openapiadmin.Domain {
	return &openapiadmin.Domain{
		CreatedAt: &req.CreatedAt,
		UpdatedAt: &req.UpdatedAt,
		Id:        ptr.To(int32(req.ID)),
		Name:      req.Name,
	}
}

func MapDomainListToResponse(req []model.Domain) openapiadmin.DomainList {
	rtn := openapiadmin.DomainList{Items: []openapiadmin.Domain{}}
	for _, item := range req {
		rtn.Items = append(rtn.Items, *MapDomainToResponse(item))
	}
	return rtn
}

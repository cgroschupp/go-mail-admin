package mapper

import (
	"fmt"

	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/utils/ptr"
)

func MapAliasToResponse(req model.Alias) openapiadmin.Alias {
	return openapiadmin.Alias{
		Id:                  ptr.To(int32(req.ID)),
		Enabled:             req.Enabled,
		DestinationDomain:   req.DestinationDomain,
		DestinationUsername: req.DestinationUsername,
		DestinationDisplay:  ptr.To(fmt.Sprintf("%s@%s", req.DestinationUsername, req.DestinationDomain)),
		SourceDomainId:      int32(req.SourceDomainID),
		SourceUsername:      &req.SourceUsername,
		SourceDisplay:       ptr.To(fmt.Sprintf("%s@%s", req.SourceUsername, req.SourceDomain.Name)), // TODO wildcard support and add domain name
		SourceDomain:        MapDomainToResponse(req.SourceDomain),
		CreatedAt:           &req.CreatedAt,
		UpdatedAt:           &req.UpdatedAt,
	}
}
func MapAliasListToResponse(req []model.Alias) openapiadmin.AliasList {
	rtn := openapiadmin.AliasList{Items: []openapiadmin.Alias{}}
	for _, item := range req {
		rtn.Items = append(rtn.Items, MapAliasToResponse(item))
	}
	return rtn
}

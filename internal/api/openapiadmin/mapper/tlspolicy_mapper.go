package mapper

import (
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/utils/ptr"
)

func MapTLSPolicyToResponse(req model.TLSPolicy) openapiadmin.TLSPolicy {
	item := openapiadmin.TLSPolicy{
		Id:       ptr.To(int32(req.ID)),
		DomainId: int32(req.DomainID),
		Policy:   openapiadmin.TLSPolicyPolicy(req.Policy),
		Domain:   MapDomainToResponse(*req.Domain),
	}

	if req.Params != nil {
		item.Params = req.Params
	}

	return item
}

func MapTLSPolcyListToResponse(req []model.TLSPolicy) openapiadmin.TLSPolicyList {
	rtn := openapiadmin.TLSPolicyList{Items: []openapiadmin.TLSPolicy{}}
	for _, item := range req {
		rtn.Items = append(rtn.Items, MapTLSPolicyToResponse(item))
	}
	return rtn
}

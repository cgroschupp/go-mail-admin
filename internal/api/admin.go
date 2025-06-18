package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin/mapper"
	"github.com/cgroschupp/go-mail-admin/internal/config"
	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type serverHandler struct {
	domainService    domain.DomainService
	aliasService     domain.AliasService
	accountService   domain.AccountService
	tlspolicyService domain.TLSPolicyService
	config           config.Config
}

func NewServerHandler(domainService domain.DomainService, aliasService domain.AliasService, accountService domain.AccountService, tlspolicyService domain.TLSPolicyService) *serverHandler {
	return &serverHandler{
		domainService:    domainService,
		aliasService:     aliasService,
		accountService:   accountService,
		tlspolicyService: tlspolicyService,
	}
}

var _ openapiadmin.ServerInterface = (*serverHandler)(nil)

// AccountsChangePassword implements ServerInterface.
func (s *serverHandler) AccountsChangePassword(w http.ResponseWriter, r *http.Request, id int32) {
	body := openapiadmin.AccountsChangePasswordJSONRequestBody{}
	err := render.Bind(r, &body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	err = s.accountService.ChangePassword(r.Context(), id, body.Password)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
}

// DashboardOperationsStats implements ServerInterface.
func (s *serverHandler) DashboardOperationsStats(w http.ResponseWriter, r *http.Request) {
	domainsStats, err := s.domainService.Stats(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	accountStats, err := s.accountService.Stats(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}

	aliasStats, err := s.aliasService.Stats(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	test := openapiadmin.DashboardStats{
		"Domains":  mapper.MapStatsToResponse(domainsStats),
		"Accounts": mapper.MapStatsToResponse(accountStats),
		"Aliases":  mapper.MapStatsToResponse(aliasStats),
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, test)
}

// TLSPoliciesCreate implements ServerInterface.
func (s *serverHandler) TLSPoliciesCreate(w http.ResponseWriter, r *http.Request) {
	body := openapiadmin.TLSPoliciesCreateJSONRequestBody{}
	err := render.Bind(r, &body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	tp, err := s.tlspolicyService.Create(r.Context(), body.Params, string(body.Policy), int(body.DomainId))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, mapper.MapTLSPolicyToResponse(tp))
}

// TLSPoliciesDelete implements ServerInterface.
func (s *serverHandler) TLSPoliciesDelete(w http.ResponseWriter, r *http.Request, id int32) {
	err := s.tlspolicyService.Delete(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
}

// TLSPoliciesList implements ServerInterface.
func (s *serverHandler) TLSPoliciesList(w http.ResponseWriter, r *http.Request) {
	tp, err := s.tlspolicyService.List(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapTLSPolcyListToResponse(tp))
}

// TLSPoliciesRead implements ServerInterface.
func (s *serverHandler) TLSPoliciesRead(w http.ResponseWriter, r *http.Request, id int32) {
	tp, err := s.tlspolicyService.Get(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapTLSPolicyToResponse(tp))
}

// TLSPoliciesUpdate implements ServerInterface.
func (s *serverHandler) TLSPoliciesUpdate(w http.ResponseWriter, r *http.Request, id int32) {
	patch := openapiadmin.TLSPolicyMergePatchUpdate{}
	err := render.DecodeJSON(r.Body, &patch)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	defer r.Body.Close()

	tp, err := s.tlspolicyService.Update(r.Context(), id, patch.Params, (*string)(patch.Policy))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapTLSPolicyToResponse(tp))
}

// UserOperationsLogout implements ServerInterface.
func (s *serverHandler) UserOperationsLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
		Secure:   s.config.Cookie.Secure,
		Path:     "/",
		Domain:   s.config.Host,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})
}

// AliasesCreate implements ServerInterface.
func (s *serverHandler) AliasesCreate(w http.ResponseWriter, r *http.Request) {
	body := openapiadmin.AliasesCreateJSONRequestBody{}
	err := render.Bind(r, &body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	a, err := s.aliasService.Create(r.Context(), body.SourceUsername, body.SourceDomainId, body.DestinationUsername, body.DestinationDomain)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	a, err = s.aliasService.Get(r.Context(), int32(a.ID))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, mapper.MapAliasToResponse(a))
}

// AliasesDelete implements ServerInterface.
func (s *serverHandler) AliasesDelete(w http.ResponseWriter, r *http.Request, id int32) {
	err := s.aliasService.Delete(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AliasesList implements ServerInterface.
func (s *serverHandler) AliasesList(w http.ResponseWriter, r *http.Request) {
	aliases, err := s.aliasService.List(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapAliasListToResponse(aliases))
}

// AliasesRead implements ServerInterface.
func (s *serverHandler) AliasesRead(w http.ResponseWriter, r *http.Request, id int32) {
	a, err := s.aliasService.Get(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapAliasToResponse(a))
}

// AliasesUpdate implements ServerInterface.
func (s *serverHandler) AliasesUpdate(w http.ResponseWriter, r *http.Request, id int32) {
	patch := openapiadmin.AliasMergePatchUpdate{}
	if err := render.DecodeJSON(r.Body, &patch); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}

	al, err := s.aliasService.Update(r.Context(), id, patch.SourceUsername, patch.SourceDomainId, patch.DestinationUsername, patch.DestinationDomain, patch.Enabled)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapAliasToResponse(al))
}

// AccountsCreate implements ServerInterface.
func (s *serverHandler) AccountsCreate(w http.ResponseWriter, r *http.Request) {
	body := openapiadmin.AccountsCreateJSONRequestBody{}
	err := render.Bind(r, &body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	var sendonly, enabled bool
	var quota int32
	if body.Sendonly != nil {
		sendonly = *body.Sendonly
	}
	if body.Enabled != nil {
		enabled = *body.Enabled
	}
	if body.Quota != nil {
		quota = 0
	}
	acc, err := s.accountService.Create(r.Context(), body.Username, body.Password, quota, sendonly, enabled, int(body.DomainId))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	acc, err = s.accountService.Get(r.Context(), int32(acc.ID))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, mapper.MapAccountToResponse(acc))
}

// AccountsDelete implements ServerInterface.
func (s *serverHandler) AccountsDelete(w http.ResponseWriter, r *http.Request, id int32) {
	err := s.accountService.Delete(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AccountsList implements ServerInterface.
func (s *serverHandler) AccountsList(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.accountService.List(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapAccountListToResponse(accounts))
}

// AccountsRead implements ServerInterface.
func (s *serverHandler) AccountsRead(w http.ResponseWriter, r *http.Request, id int32) {
	account, err := s.accountService.Get(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapAccountToResponse(account))
}

// AccountsUpdate implements ServerInterface.
func (s *serverHandler) AccountsUpdate(w http.ResponseWriter, r *http.Request, id int32) {
	patch := openapiadmin.AccountMergePatchUpdate{}
	err := render.DecodeJSON(r.Body, &patch)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	defer r.Body.Close()

	acc, err := s.accountService.Update(r.Context(), id, patch.Username, patch.Quota, patch.Sendonly, patch.Enabled)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapAccountToResponse(acc))
}

// DomainsCreate implements ServerInterface.
func (s *serverHandler) DomainsCreate(w http.ResponseWriter, r *http.Request) {
	d := &openapiadmin.Domain{}
	err := render.Bind(r, d)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	if err := validate.Struct(d); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	do, err := s.domainService.Create(r.Context(), d.Name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		if errors.Is(err, domain.ErrDomainExists) {
			render.Status(r, http.StatusConflict)
		}
		_ = render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, do)
}

// DomainsDelete implements ServerInterface.
func (s *serverHandler) DomainsDelete(w http.ResponseWriter, r *http.Request, id int32) {
	err := s.domainService.Delete(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DomainsList implements ServerInterface.
func (s *serverHandler) DomainsList(w http.ResponseWriter, r *http.Request) {
	domains, err := s.domainService.List(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapDomainListToResponse(domains))
}

// DomainsRead implements ServerInterface.
func (s *serverHandler) DomainsRead(w http.ResponseWriter, r *http.Request, id int32) {
	do, err := s.domainService.Get(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		if errors.Is(err, domain.ErrDomainNotExists) {
			render.Status(r, http.StatusNotFound)
		}
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.JSON(w, r, mapper.MapDomainToResponse(do))
}

// DomainsUpdate implements ServerInterface.
func (s *serverHandler) DomainsUpdate(w http.ResponseWriter, r *http.Request, id int32) {
	patch := openapiadmin.DomainsCreateJSONRequestBody{}

	err := render.DecodeJSON(r.Body, &patch)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	defer r.Body.Close()
	if err := validate.Struct(patch); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	do, err := s.domainService.Update(r.Context(), id, patch.Name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapper.MapDomainToResponse(do))
}

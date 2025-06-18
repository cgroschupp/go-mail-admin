package api

import (
	"net/http"
	"time"

	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiauth"
	"github.com/cgroschupp/go-mail-admin/internal/config"
	"github.com/cgroschupp/go-mail-admin/internal/domain"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type authHandler struct {
	userService domain.UserService
	jwtAuth     *jwtauth.JWTAuth
	config      *config.Config
}

func NewAuthHandler(userService domain.UserService, config *config.Config, jwtAuth *jwtauth.JWTAuth) *authHandler {
	return &authHandler{
		userService: userService,
		config:      config,
		jwtAuth:     jwtAuth,
	}
}

// UserOperationsLogin implements ServerInterface.
func (s *authHandler) UserOperationsLogin(w http.ResponseWriter, r *http.Request) {
	resp := openapiauth.LoginParameter{}
	err := render.Bind(r, &resp)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	ok, err := s.userService.Login(resp.Username, resp.Password)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiadmin.Error{Error: err.Error()})
		return
	}
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{})
		return
	}
	claim := map[string]interface{}{"admin": true}
	jwtauth.SetExpiryIn(claim, time.Hour*5)
	jwtauth.SetIssuedNow(claim)
	_, tokenString, err := s.jwtAuth.Encode(claim)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, openapiauth.Error{Error: err.Error(), Message: "unable to parse jwt token"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: tokenString, HttpOnly: true, Secure: s.config.Cookie.Secure, Path: "/", Domain: s.config.Host, SameSite: http.SameSiteLaxMode})
	render.Status(r, http.StatusOK)
	render.JSON(w, r, openapiauth.LoginResponse{Login: true, Token: tokenString})
}

var _ openapiauth.ServerInterface = (*authHandler)(nil)

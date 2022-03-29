package internal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

var tokenAuth *jwtauth.JWTAuth

func (m *MailServerConfiguratorInterface) login(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("Login new JWT Function")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	loginResult := struct {
		Login bool   `json:"login"`
		Token string `json:"token"`
	}{
		Login: false,
		Token: "",
	}

	json.Unmarshal(body, &loginData)

	if loginData.Username == m.Config.Auth.Username && loginData.Password == m.Config.Auth.Password {
		claim := map[string]interface{}{"admin": true}

		jwtauth.SetIssuedNow(claim)
		jwtauth.SetExpiry(claim, time.Now().Add(m.Config.Auth.Expire))
		_, tokenString, _ := tokenAuth.Encode(claim)
		loginResult.Token = tokenString
		loginResult.Login = true
	}

	render.JSON(w, r, loginResult)
}

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/cgroschupp/go-mail-admin/internal"
	"github.com/cgroschupp/go-mail-admin/internal/api/openapiauth"
	"github.com/cgroschupp/go-mail-admin/internal/config"
	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite
	Server    *internal.MailServerConfiguratorInterface
	Token     string
	crsfToken string
	cookies   []*http.Cookie
}

func (suite *DBSuite) SetupTest() {
	s := internal.NewMailServerConfiguratorInterface(&config.Config{
		Database: config.DatabaseConfig{Type: "sqlite", DSN: "unittest.db"},
		Password: config.PasswordConfig{Scheme: "SSHA512"},
		Auth:     config.AuthConfig{Username: "unittest", Password: "unittest", Secret: "unittest", Expire: 1 * time.Hour},
	})

	suite.NoError(s.ConnectToDb())
	s.MountHandlers()

	req := httptest.NewRequest("GET", "/api/v1/csrf", nil)
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
	csrfResult := openapiauth.CsrfResponse{}
	err := json.NewDecoder(rr.Body).Decode(&csrfResult)
	suite.NoError(err)
	suite.crsfToken = csrfResult.CsrfToken
	req = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBufferString("{\"username\":\"unittest\",\"password\":\"unittest\"}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-CSRF-Token", suite.crsfToken)
	req.Header.Add("Origin", "localhost")
	suite.cookies = rr.Result().Cookies()
	for _, c := range rr.Result().Cookies() {
		req.AddCookie(c)
	}
	rr = httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	result := openapiauth.LoginResponse{}
	err = json.NewDecoder(rr.Body).Decode(&result)
	suite.NoError(err)
	suite.True(result.Login)
	suite.Token = result.Token

	suite.Server = s
}

func (suite *DBSuite) TearDownTest() {
	err := os.Remove("unittest.db")
	suite.NoError(err)
}

func (suite *DBSuite) Request(method, path string, contentType string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", suite.Token))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-CSRF-Token", suite.crsfToken)
	req.Header.Add("Origin", "localhost")
	for _, c := range suite.cookies {
		req.AddCookie(c)
	}
	rr := httptest.NewRecorder()
	suite.Server.Router.ServeHTTP(rr, req)
	return rr
}

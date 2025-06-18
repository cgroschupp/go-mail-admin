package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cgroschupp/go-mail-admin/internal/api/openapiadmin"
	"github.com/cgroschupp/go-mail-admin/internal/model"
	"github.com/cgroschupp/go-mail-admin/internal/test"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AccountTestSuite struct {
	test.DBSuite
}

func (suite *AccountTestSuite) TestListAccountEmpty() {
	rr := suite.Request("GET", "/api/v1/account", "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	var accounts openapiadmin.AccountList
	err := json.NewDecoder(rr.Body).Decode(&accounts)

	if suite.NoError(err, "Error decoding response body") {
		suite.Len(accounts.Items, 0)
	}
}

func (suite *AccountTestSuite) TestListAccounts() {
	user := model.Account{Username: "test", Domain: &model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&user).Error)
	rr := suite.Request("GET", "/api/v1/account", "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	var accounts openapiadmin.AccountList
	err := json.NewDecoder(rr.Body).Decode(&accounts)

	if suite.NoError(err, "Error decoding response body") {
		if suite.Len(accounts.Items, 1) {
			suite.Equal("test", accounts.Items[0].Username)
			suite.Equal("test.de", accounts.Items[0].Domain.Name)
		}
	}
}

func (suite *AccountTestSuite) TestAddAccount() {
	d := model.Domain{Name: "test.de"}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)
	rr := suite.Request("POST", "/api/v1/account", "application/json", bytes.NewBufferString("{\"username\":\"test\",\"domain_id\":1, \"password\":\"test\"}"))

	suite.Require().Equal(http.StatusCreated, rr.Code)

	var a openapiadmin.Account
	err := json.NewDecoder(rr.Body).Decode(&a)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal("test", a.Username)
	}
}

func (suite *AccountTestSuite) TestGetAccount() {
	a := model.Account{Username: "test", Domain: &model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&a).Error)

	rr := suite.Request("GET", fmt.Sprintf("/api/v1/account/%d", a.ID), "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	var a2 openapiadmin.Account
	err := json.NewDecoder(rr.Body).Decode(&a2)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(a.Username, a2.Username)
		suite.Equal(a.Domain.ID, uint(a2.DomainId))
	}
}

func (suite *AccountTestSuite) TestDeleteAccount() {
	a := model.Account{Username: "test", Domain: &model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&a).Error)

	rr := suite.Request("DELETE", fmt.Sprintf("/api/v1/account/%d", a.ID), "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	suite.ErrorIs(suite.Server.DB.First(a).Error, gorm.ErrRecordNotFound)
}

func (suite *AccountTestSuite) TestUpdateAccount() {
	a := model.Account{Username: "test", Domain: &model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&a).Error)

	rr := suite.Request("PATCH", fmt.Sprintf("/api/v1/account/%d", a.ID), "application/merge-patch+json", bytes.NewBufferString("{\"username\":\"test2\"}"))
	suite.Require().Equal(http.StatusOK, rr.Code)
	var a2 openapiadmin.Account
	err := json.NewDecoder(rr.Body).Decode(&a2)
	if suite.NoError(err) {
		suite.Equal("test2", a2.Username)
	}
}

func (suite *AccountTestSuite) TestUpdateAccountPassword() {
	a := model.Account{Username: "test", Domain: &model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&a).Error)

	suite.Require().NoError(suite.Server.DB.First(&a).Error)
	oldPasswordHash := a.Password

	rr := suite.Request("PUT", fmt.Sprintf("/api/v1/account/%d/password", a.ID), "application/json", bytes.NewBufferString("{\"password\":\"test2\"}"))

	suite.Require().Equal(http.StatusOK, rr.Code)
	suite.Empty(rr.Body)
	suite.Require().NoError(suite.Server.DB.First(&a).Error)
	suite.NotEqual(oldPasswordHash, a.Password)
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

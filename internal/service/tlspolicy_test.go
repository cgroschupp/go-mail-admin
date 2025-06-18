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

type TLSPolicyTestSuite struct {
	test.DBSuite
}

func (suite *TLSPolicyTestSuite) TestListTLSPolicyEmpty() {
	rr := suite.Request("GET", "/api/v1/tlspolicy", "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	var tlsPolicies openapiadmin.TLSPolicyList
	err := json.NewDecoder(rr.Body).Decode(&tlsPolicies)

	if suite.NoError(err, "Error decoding response body") {
		suite.Len(tlsPolicies.Items, 0)
	}
}

func (suite *TLSPolicyTestSuite) TestListTLSPolicy() {
	tlsPolicy := model.TLSPolicy{
		Domain: &model.Domain{Name: "test.de"},
		Policy: "may",
		Params: String("test"),
	}
	suite.Require().NoError(suite.Server.DB.Create(&tlsPolicy).Error)
	rr := suite.Request("GET", "/api/v1/tlspolicy", "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	var tlsPolicies openapiadmin.TLSPolicyList
	err := json.NewDecoder(rr.Body).Decode(&tlsPolicies)

	if suite.NoError(err, "Error decoding response body") {
		if suite.Len(tlsPolicies.Items, 1) {
			suite.Equal(openapiadmin.TLSPolicyPolicyMay, tlsPolicies.Items[0].Policy)
			suite.Equal("test", *tlsPolicies.Items[0].Params)
			suite.Equal("test.de", tlsPolicies.Items[0].Domain.Name)
		}
	}
}

func (suite *TLSPolicyTestSuite) TestGetTLSPolicy() {
	tlsPolicy := model.TLSPolicy{
		Domain: &model.Domain{Name: "test.de"},
		Policy: "may",
		Params: String("test"),
	}
	suite.Require().NoError(suite.Server.DB.Create(&tlsPolicy).Error)
	rr := suite.Request("GET", fmt.Sprintf("/api/v1/tlspolicy/%d", tlsPolicy.ID), "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)

	err := json.NewDecoder(rr.Body).Decode(&tlsPolicy)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal("may", tlsPolicy.Policy)
		suite.Equal("test", *tlsPolicy.Params)
		suite.Equal(uint(1), tlsPolicy.DomainID)
	}
}

func String(input string) *string {
	return &input
}

func (suite *TLSPolicyTestSuite) TestDeleteTLSPolicy() {
	tlsPolicy := model.TLSPolicy{
		Domain: &model.Domain{Name: "test.de"},
		Policy: "may",
		Params: String("test"),
	}
	suite.Require().NoError(suite.Server.DB.Create(&tlsPolicy).Error)

	rr := suite.Request("DELETE", fmt.Sprintf("/api/v1/tlspolicy/%d", tlsPolicy.ID), "application/json", nil)

	suite.Require().Equal(http.StatusOK, rr.Code)

	suite.ErrorIs(suite.Server.DB.First(tlsPolicy).Error, gorm.ErrRecordNotFound)
}

func (suite *TLSPolicyTestSuite) TestUpdateTLSPolicy() {
	tlsPolicy := model.TLSPolicy{
		Domain: &model.Domain{Name: "test.de"},
		Policy: "may",
		Params: String("test"),
	}
	suite.Require().NoError(suite.Server.DB.Create(&tlsPolicy).Error)

	rr := suite.Request("PATCH", fmt.Sprintf("/api/v1/tlspolicy/%d", tlsPolicy.ID), "application/merge-patch+json", bytes.NewBufferString("{\"policy\":\"dane\",\"params\":\"\"}"))

	suite.Require().Equal(http.StatusOK, rr.Code)

	result := openapiadmin.TLSPolicy{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(openapiadmin.TLSPolicyPolicyDane, result.Policy)
		suite.Equal("", *result.Params)
		suite.Equal(tlsPolicy.ID, uint(*result.Id))
	}
}

func (suite *TLSPolicyTestSuite) TestAddTLSPolicy() {
	d := model.Domain{Name: "test.de"}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("POST", "/api/v1/tlspolicy", "application/json", bytes.NewBufferString("{\"domain_id\":1, \"policy\":\"may\"}"))

	suite.Require().Equal(http.StatusCreated, rr.Code)

	var tlsPolicy model.TLSPolicy
	err := json.NewDecoder(rr.Body).Decode(&tlsPolicy)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal("may", tlsPolicy.Policy)
		suite.Empty(tlsPolicy.Params)
		suite.Equal(uint(1), d.ID)
	}
}

func TestTLSPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(TLSPolicyTestSuite))
}

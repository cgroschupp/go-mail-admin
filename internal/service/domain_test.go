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
)

const (
	unittest_domain_name = "test.de"
)

type DomainTestSuite struct {
	test.DBSuite
}

func (suite *DomainTestSuite) TestListDomainsEmpty() {
	rr := suite.Request("GET", "/api/v1/domain", "application/json", nil)
	suite.Equal(http.StatusOK, rr.Code)

	var domains openapiadmin.DomainList
	err := json.NewDecoder(rr.Body).Decode(&domains)

	if suite.NoError(err, "Error decoding response body") {
		suite.Len(domains.Items, 0)
	}
}

func (suite *DomainTestSuite) TestListDomains() {
	d := model.Domain{Name: unittest_domain_name}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("GET", "/api/v1/domain", "application/json", nil)
	suite.Equal(http.StatusOK, rr.Code)

	var domains openapiadmin.DomainList
	err := json.NewDecoder(rr.Body).Decode(&domains)

	if suite.NoError(err, "Error decoding response body") {
		if suite.Len(domains.Items, 1) {
			suite.Equal(unittest_domain_name, domains.Items[0].Name)
		}
	}
}

func (suite *DomainTestSuite) TestAddDomain() {
	rr := suite.Request("POST", "/api/v1/domain", "application/json", bytes.NewBufferString("{\"name\":\"test.de\"}"))

	suite.Require().Equal(http.StatusCreated, rr.Code)

	var domain openapiadmin.Domain
	err := json.NewDecoder(rr.Body).Decode(&domain)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(unittest_domain_name, domain.Name)
	}
}

func (suite *DomainTestSuite) TestUpdateDomain() {
	d := model.Domain{Name: unittest_domain_name}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("PATCH", fmt.Sprintf("/api/v1/domain/%d", d.ID), "application/merge-patch+json", bytes.NewBufferString("{\"name\":\"test2.de\"}"))
	suite.Require().Equal(http.StatusOK, rr.Code)

	result := openapiadmin.Domain{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	if suite.NoError(err) {
		suite.Equal(d.ID, uint(*result.Id))
		suite.Equal("test2.de", result.Name)
	}
}

func (suite *DomainTestSuite) TestAddDomainConflict() {
	d := model.Domain{Name: unittest_domain_name}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("POST", "/api/v1/domain", "application/json", bytes.NewBufferString("{\"name\":\"test.de\"}"))
	suite.Require().Equal(http.StatusConflict, rr.Code)
}

func (suite *DomainTestSuite) TestAddDomainInvalid() {
	d := model.Domain{Name: unittest_domain_name}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)
	rr := suite.Request("POST", "/api/v1/domain", "application/json", bytes.NewBufferString("{\"name\":\"-invalid\"}"))
	suite.Require().Equal(http.StatusBadRequest, rr.Code)
}

func (suite *DomainTestSuite) TestGetDomain() {
	d := model.Domain{Name: unittest_domain_name}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("GET", fmt.Sprintf("/api/v1/domain/%d", d.ID), "application/json", nil)
	suite.Require().Equal(http.StatusOK, rr.Code)

	var domain model.Domain
	err := json.NewDecoder(rr.Body).Decode(&domain)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(unittest_domain_name, domain.Name)
	}
}

func (suite *DomainTestSuite) TestGetDomainNotExists() {
	rr := suite.Request("GET", "/api/v1/domain/1", "application/json", nil)

	suite.Require().Equal(http.StatusNotFound, rr.Code)
}

func (suite *DomainTestSuite) TestDeleteDomain() {
	d := model.Domain{Name: unittest_domain_name}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("GET", fmt.Sprintf("/api/v1/domain/%d", d.ID), "application/json", nil)
	suite.Require().Equal(http.StatusOK, rr.Code)
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}

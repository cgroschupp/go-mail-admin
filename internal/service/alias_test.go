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

type AliasTestSuite struct {
	test.DBSuite
}

func (suite *AliasTestSuite) TestListAliasEmpty() {
	rr := suite.Request("GET", "/api/v1/alias", "application/json", nil)
	suite.Equal(http.StatusOK, rr.Code)

	var aliases openapiadmin.AliasList
	err := json.NewDecoder(rr.Body).Decode(&aliases)

	if suite.NoError(err, "Error decoding response body") {
		suite.Len(aliases.Items, 0)
	}
}

func (suite *AliasTestSuite) TestListAlias() {
	d := model.Alias{SourceUsername: "test", SourceDomain: model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)
	rr := suite.Request("GET", "/api/v1/alias", "application/json", nil)
	suite.Equal(http.StatusOK, rr.Code)

	var aliaes openapiadmin.AliasList
	err := json.NewDecoder(rr.Body).Decode(&aliaes)

	if suite.NoError(err, "Error decoding response body") {
		if suite.Len(aliaes.Items, 1) {
			suite.Equal(d.SourceDomainID, uint(aliaes.Items[0].SourceDomainId))
		}
	}
}

func (suite *AliasTestSuite) TestAddAlias() {
	rr := suite.Request("POST", "/api/v1/alias", "application/json", bytes.NewBufferString("{\"source_username\":\"test\", \"source_domain_id\":1, \"destination_username\":\"foo\", \"destination_domain\":\"google.com\", \"enabled\":true}"))
	suite.Require().Equal(http.StatusCreated, rr.Code)

	var alias openapiadmin.Alias
	err := json.NewDecoder(rr.Body).Decode(&alias)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(int32(1), alias.SourceDomainId)
		suite.Equal("test", *alias.SourceUsername)
		suite.Equal("google.com", alias.DestinationDomain)
		suite.Equal("foo", alias.DestinationUsername)
		suite.Equal(true, alias.Enabled)
	}
}

func (suite *AliasTestSuite) TestAddAliasCatchAll() {
	// suite.Server.Config.Feature.CatchAll = true
	rr := suite.Request("POST", "/api/v1/alias", "application/json", bytes.NewBufferString("{\"source_domain_id\":1, \"destination_username\":\"foo\", \"destination_domain\":\"google.com\", \"enabled\":true}"))
	suite.Require().Equal(http.StatusCreated, rr.Code)

	var alias model.Alias
	err := json.NewDecoder(rr.Body).Decode(&alias)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(uint(1), alias.SourceDomainID)
		suite.Equal("", alias.SourceUsername)
		suite.Equal("google.com", alias.DestinationDomain)
		suite.Equal("foo", alias.DestinationUsername)
		suite.Equal(true, alias.Enabled)
	}
}

func (suite *AliasTestSuite) TestGetAlias() {
	d := model.Alias{SourceUsername: "test", SourceDomain: model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&d).Error)

	rr := suite.Request("GET", fmt.Sprintf("/api/v1/alias/%d", d.ID), "application/json", nil)
	suite.Equal(http.StatusOK, rr.Code)

	var alias model.Alias
	err := json.NewDecoder(rr.Body).Decode(&alias)

	if suite.NoError(err, "Error decoding response body") {
		suite.Equal(d.SourceUsername, alias.SourceUsername)
		suite.Equal(d.SourceDomainID, alias.SourceDomainID)
	}
}

func (suite *AliasTestSuite) TestDeleteAlias() {
	a := model.Alias{SourceUsername: "test", SourceDomain: model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&a).Error)

	rr := suite.Request("DELETE", fmt.Sprintf("/api/v1/alias/%d", a.ID), "application/json", nil)

	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *AliasTestSuite) TestUpdateAlias() {
	a := model.Alias{SourceUsername: "test", SourceDomain: model.Domain{Name: "test.de"}}
	suite.Require().NoError(suite.Server.DB.Create(&a).Error)

	rr := suite.Request("PATCH", fmt.Sprintf("/api/v1/alias/%d", a.ID), "application/merge-patch+json", bytes.NewBufferString("{\"source_username\":\"test\", \"source_domain_id\": 1, \"destination_username\":\"foo\", \"destination_domain\":\"google.com\", \"enabled\":true}"))
	suite.Require().Equal(http.StatusOK, rr.Code)

	result := openapiadmin.Alias{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	if suite.NoError(err) {
		suite.Equal(a.ID, uint(*result.Id))
		suite.Equal(int32(1), result.SourceDomainId)
		suite.Equal("test", *result.SourceUsername)
		suite.Equal("google.com", result.DestinationDomain)
		suite.Equal("foo", result.DestinationUsername)
	}
}

func TestAliasTestSuite(t *testing.T) {
	suite.Run(t, new(AliasTestSuite))
}

package infrastructure

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/go-ldap/ldap/v3"
	"net/http"
)

func selfLoginToLDAP(l *log.DXLog, user utils.JSON, userPassword string, organizationAuthSource string, organizationAttribute string) (isSuccess bool, err error) {
	c := *configuration.Manager.Configurations["external_system"].Data

	ldapConfig, ok := c[organizationAuthSource].(utils.JSON)
	if !ok {
		return false, l.WarnAndCreateErrorf("SELF_LOGIN_TO_LDAP:LDAP_CONFIG_NOT_FOUND:%s", organizationAuthSource)
	}

	ldapAddress, ok := ldapConfig["address"].(string)
	if !ok {
		return false, l.WarnAndCreateErrorf("SELF_LOGIN_TO_LDAP:LDAP_ADDRESS_NOT_FOUND:%s", organizationAuthSource)
	}
	ldapConnection, err := ldap.DialURL(ldapAddress)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = ldapConnection.Close()
	}()

	userAttribute := user["attribute"].(string)
	attribute := userAttribute

	if organizationAttribute != "" {
		attribute = attribute + "," + organizationAttribute
	}

	// Bind as the user to verify their password
	err = ldapConnection.Bind(attribute, userPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}

func doOnAuthenticateUser(aepr *api.DXAPIEndPointRequest, userLoginId string, userPassword string, organizationUId string) (verificationResult bool, user utils.JSON, organization utils.JSON /*organizations []utils.JSON,*/, err error) {
	_, user, err = user_management.ModuleUserManagement.User.SelectOne(&aepr.Log, nil, utils.JSON{
		"loginid": userLoginId,
	}, nil, nil)
	if err != nil {
		return false, nil, nil, err
	}
	if user == nil {
		return false, nil, nil, aepr.WriteResponseAndNewErrorf(http.StatusUnauthorized, "", "USER_NOT_FOUND")
	}

	userId := user["id"].(int64)

	us := utils.JSON{
		"user_id": userId,
	}

	if organizationUId != "" {
		us["organization_uid"] = organizationUId
	}

	_, userOrganizationMembership, err := user_management.ModuleUserManagement.UserOrganizationMembership.SelectOne(
		&aepr.Log, nil, us, nil, map[string]string{"order_index": "asc"},
	)
	if err != nil {
		return false, nil, nil, err
	}

	if userOrganizationMembership == nil {
		return false, nil, nil, aepr.WriteResponseAndNewErrorf(http.StatusUnauthorized, "", "USER_ORGANIZATION_MEMBERSHIP_NOT_FOUND")
	}

	organizationId := userOrganizationMembership["organization_id"].(int64)
	_, organization, err = user_management.ModuleUserManagement.Organization.GetById(&aepr.Log, organizationId)
	if err != nil {
		return false, nil, nil, err
	}
	if organization == nil {
		return false, nil, nil, aepr.WriteResponseAndNewErrorf(http.StatusUnauthorized, "", "ORGANIZATION_NOT_FOUND")
	}
	//organizations = []utils.JSON{userOrganizationMembership}

	authSource1, ok := userOrganizationMembership["organization_auth_source1"].(string)
	if !ok {
		authSource1 = ""
	}
	attribute1, ok := userOrganizationMembership["organization_attribute1"].(string)
	if !ok {
		attribute1 = ""
	}
	authSource2, ok := userOrganizationMembership["organization_auth_source2"].(string)
	if !ok {
		authSource2 = ""
	}
	attribute2, ok := userOrganizationMembership["organization_attribute2"].(string)
	if !ok {
		attribute2 = ""
	}

	if authSource1 != "" {
		verificationResult, err = selfLoginToLDAP(&aepr.Log, user, userPassword, authSource1, attribute1)
		if err == nil {
			return verificationResult, user, organization, nil
		}
		aepr.Log.Warn(fmt.Sprintf("AUTHSOURCE1_EXTERNAL_SYSTEM_FAILED:%s", err.Error()))
	}

	if authSource2 != "" {
		verificationResult, err = selfLoginToLDAP(&aepr.Log, user, userPassword, authSource2, attribute2)
		if err == nil {
			return verificationResult, user, organization, nil
		}
		aepr.Log.Warn(fmt.Sprintf("AUTHSOURCE2_EXTERNAL_SYSTEM_FAILED:%s", err.Error()))
	}

	verificationResult, err = user_management.ModuleUserManagement.UserPasswordVerify(&aepr.Log, userId, userPassword)
	if err != nil {
		return false, nil, nil, err
	}

	return verificationResult, user, organization, nil
}

func doOnCreateSessionObject(aepr *api.DXAPIEndPointRequest, user utils.JSON, userLoggedOrganization, originalSessionObject utils.JSON) (newSessionObject utils.JSON, err error) {
	userId := user["id"].(int64)

	_, selfOrganizationMemberships, err := user_management.ModuleUserManagement.UserOrganizationMembership.Select(&aepr.Log, nil, utils.JSON{
		"user_id": userId,
	}, nil, map[string]string{"order_index": "asc"}, 0)
	if err != nil {
		return originalSessionObject, err
	}
	if selfOrganizationMemberships != nil {
		originalSessionObject["user_organization_memberships"] = selfOrganizationMemberships
	}
	return originalSessionObject, nil
}

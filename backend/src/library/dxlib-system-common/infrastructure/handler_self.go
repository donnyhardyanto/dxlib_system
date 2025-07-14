package infrastructure

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/configuration"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/go-ldap/ldap/v3"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
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
		for i, selfOrganizationMembership := range selfOrganizationMemberships {
			organizationId := selfOrganizationMembership["organization_id"].(int64)

			_, organizationExecutor, err := partner_management.ModulePartnerManagement.OrganizationExecutor.SelectOne(&aepr.Log, nil, utils.JSON{
				"organization_id": organizationId,
			}, nil, map[string]string{"id": "asc"})
			if err != nil {
				return originalSessionObject, err
			}

			if organizationExecutor != nil {
				organizationRoleId := organizationExecutor["id"].(int64)

				_, organizationExecutorExpertises, err := partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.Select(&aepr.Log, nil, utils.JSON{
					"organization_role_id": organizationRoleId,
				}, nil, map[string]string{"id": "asc"}, 0, nil)
				if err != nil {
					return originalSessionObject, err
				}

				_, organizationExecutorLocations, err := partner_management.ModulePartnerManagement.OrganizationExecutorLocation.Select(&aepr.Log, nil, utils.JSON{
					"organization_role_id": organizationRoleId,
				}, nil, map[string]string{"id": "asc"}, 0, nil)
				if err != nil {
					return originalSessionObject, err
				}

				_, organizationExecutorAreas, err := partner_management.ModulePartnerManagement.OrganizationExecutorArea.Select(&aepr.Log, nil, utils.JSON{
					"organization_role_id": organizationRoleId,
				}, nil, map[string]string{"id": "asc"}, 0, nil)
				if err != nil {
					return originalSessionObject, err
				}

				organizationExecutor["organization_executor_expertises"] = organizationExecutorExpertises
				organizationExecutor["organization_executor_locations"] = organizationExecutorLocations
				organizationExecutor["organization_executor_areas"] = organizationExecutorAreas
			}

			selfOrganizationMemberships[i]["organization_executor"] = organizationExecutor

			_, organizationSupervisor, err := partner_management.ModulePartnerManagement.OrganizationSupervisor.SelectOne(&aepr.Log, nil, utils.JSON{
				"organization_id": organizationId,
			}, nil, map[string]string{"id": "asc"})
			if err != nil {
				return originalSessionObject, err
			}

			if organizationSupervisor != nil {
				organizationRoleId := organizationSupervisor["id"].(int64)

				_, organizationSupervisorExpertises, err := partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.Select(&aepr.Log, nil, utils.JSON{
					"organization_role_id": organizationRoleId,
				}, nil, map[string]string{"id": "asc"}, 0, nil)
				if err != nil {
					return originalSessionObject, err
				}

				_, organizationSupervisorLocations, err := partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.Select(&aepr.Log, nil, utils.JSON{
					"organization_role_id": organizationRoleId,
				}, nil, map[string]string{"id": "asc"}, 0, nil)
				if err != nil {
					return originalSessionObject, err
				}

				_, organizationSupervisorAreas, err := partner_management.ModulePartnerManagement.OrganizationSupervisorArea.Select(&aepr.Log, nil, utils.JSON{
					"organization_role_id": organizationRoleId,
				}, nil, map[string]string{"id": "asc"}, 0, nil)
				if err != nil {
					return originalSessionObject, err
				}

				organizationSupervisor["organization_supervisor_expertises"] = organizationSupervisorExpertises
				organizationSupervisor["organization_supervisor_locations"] = organizationSupervisorLocations
				organizationSupervisor["organization_supervisor_areas"] = organizationSupervisorAreas
			}

			selfOrganizationMemberships[i]["organization_supervisor"] = organizationSupervisor

			_, selfFieldExecutor, err := partner_management.ModulePartnerManagement.FieldExecutor.SelectOne(&aepr.Log, nil, utils.JSON{
				"user_id":         userId,
				"organization_id": organizationId,
			}, nil, map[string]string{"id": "asc"})
			if err != nil {
				return originalSessionObject, err
			}
			if selfFieldExecutor != nil {
				userRoleMembershipId := selfFieldExecutor["id"].(int64)

				_, selfFieldExecutorExpertise, err := partner_management.ModulePartnerManagement.FieldExecutorExpertise.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, map[string]string{"id": "asc"}, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldExecutor["field_executor_expertises"] = selfFieldExecutorExpertise

				_, selfFieldExecutorLocation, err := partner_management.ModulePartnerManagement.FieldExecutorLocation.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, map[string]string{"id": "asc"}, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldExecutor["field_executor_locations"] = selfFieldExecutorLocation

				_, selfFieldExecutorArea, err := partner_management.ModulePartnerManagement.FieldExecutorArea.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, map[string]string{"id": "asc"}, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldExecutor["field_executor_areas"] = selfFieldExecutorArea

				_, selfFieldExecutorEffectiveExpertise, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveExpertise.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, nil, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldExecutor["field_executor_effective_expertises"] = selfFieldExecutorEffectiveExpertise

				_, selfFieldExecutorEffectiveLocation, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveLocation.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, nil, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldExecutor["field_executor_effective_locations"] = selfFieldExecutorEffectiveLocation

				_, selfFieldExecutorEffectiveArea, err := partner_management.ModulePartnerManagement.FieldExecutorEffectiveArea.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, nil, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldExecutor["field_executor_effective_areas"] = selfFieldExecutorEffectiveArea

				selfOrganizationMemberships[i]["field_executor"] = selfFieldExecutor
			}

			_, selfFieldSupervisor, err := partner_management.ModulePartnerManagement.FieldSupervisor.SelectOne(&aepr.Log, nil, utils.JSON{
				"user_id":         userId,
				"organization_id": organizationId,
			}, nil, map[string]string{"id": "asc"})
			if err != nil {
				return originalSessionObject, err
			}

			if selfFieldSupervisor != nil {
				userRoleMembershipId := selfFieldSupervisor["id"].(int64)

				_, selfFieldSupervisorExpertise, err := partner_management.ModulePartnerManagement.FieldSupervisorExpertise.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, map[string]string{"id": "asc"}, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldSupervisor["field_supervisor_expertises"] = selfFieldSupervisorExpertise

				_, selfFieldSupervisorLocation, err := partner_management.ModulePartnerManagement.FieldSupervisorLocation.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, map[string]string{"id": "asc"}, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldSupervisor["field_supervisor_locations"] = selfFieldSupervisorLocation

				_, selfFieldSupervisorArea, err := partner_management.ModulePartnerManagement.FieldSupervisorArea.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, map[string]string{"id": "asc"}, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldSupervisor["field_supervisor_areas"] = selfFieldSupervisorArea

				_, selfFieldSupervisorEffectiveExpertise, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveExpertise.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, nil, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldSupervisor["field_supervisor_effective_expertises"] = selfFieldSupervisorEffectiveExpertise

				_, selfFieldSupervisorEffectiveLocation, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveLocation.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, nil, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldSupervisor["field_supervisor_effective_locations"] = selfFieldSupervisorEffectiveLocation

				_, selfFieldSupervisorEffectiveArea, err := partner_management.ModulePartnerManagement.FieldSupervisorEffectiveArea.Select(&aepr.Log, nil, utils.JSON{
					"user_role_membership_id": userRoleMembershipId,
				}, nil, nil, nil, nil)
				if err != nil {
					return originalSessionObject, err
				}
				selfFieldSupervisor["field_supervisor_effective_areas"] = selfFieldSupervisorEffectiveArea

				selfOrganizationMemberships[i]["field_supervisor"] = selfFieldSupervisor
			}
		}
		originalSessionObject["user_organization_memberships"] = selfOrganizationMemberships
	}
	return originalSessionObject, nil
}

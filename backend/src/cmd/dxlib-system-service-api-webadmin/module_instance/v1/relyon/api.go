package relyon

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	utilsHTTP "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib-system/common/external/relyon"
	"net/http"
)

func HandleRelyOnRegisterInstallationUpdateStatus(aepr *api.DXAPIEndPointRequest) (err error) {
	_, registrationId, err := aepr.GetParameterValueAsString("registration_id")
	if err != nil {
		return err
	}
	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return err
	}
	_, r, err := relyon.RegisterInstallationUpdateStatus(&aepr.Log, session, registrationId)
	if err != nil {
		return err
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": r})
	return nil
}

func HandleRelyOnRegisterInstallationUpdateType(aepr *api.DXAPIEndPointRequest) (err error) {
	_, registrationId, err := aepr.GetParameterValueAsString("registration_id")
	if err != nil {
		return err
	}
	_, installationTypeId, err := aepr.GetParameterValueAsInt64("installation_type_id")
	if err != nil {
		return err
	}
	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return err
	}
	_, err = relyon.RegisterInstallationUpdateType(&aepr.Log, session, registrationId, installationTypeId)
	if err != nil {
		return err
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func HandleRelyOnCancelSubscriptionMileStone(aepr *api.DXAPIEndPointRequest) (err error) {
	_, registrationId, err := aepr.GetParameterValueAsString("registration_id")
	if err != nil {
		return err
	}
	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return err
	}
	milestones, err := relyon.CancelSubscriptionMileStone(&aepr.Log, session, registrationId)
	if err != nil {
		return err
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"data": utils.JSON{
			"milestones": milestones,
		},
	})
	return nil
}

func HandleRelyOnCancelSubscriptionCreate(aepr *api.DXAPIEndPointRequest) (err error) {
	_, registrationId, err := aepr.GetParameterValueAsString("registration_id")
	if err != nil {
		return err
	}
	_, channel, err := aepr.GetParameterValueAsString("channel")
	if err != nil {
		return err
	}
	_, cancelBy, err := aepr.GetParameterValueAsString("cancel_by")
	if err != nil {
		return err
	}
	_, reason, err := aepr.GetParameterValueAsArrayOfString("reason")
	if err != nil {
		return err
	}
	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return err
	}
	err = relyon.CancelSubscriptionCreate(&aepr.Log, session, registrationId, reason, channel, cancelBy)
	if err != nil {
		return err
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, nil)
	return nil
}

func HandleAuth(aepr *api.DXAPIEndPointRequest) (err error) {
	_, session, err := relyon.Auth(&aepr.Log)
	if err != nil {
		return err
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": session})
	return nil
}

func DefineAPIEndPoints(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("RelyOn",
		"Test Auth",
		"/v1/relyon/auth/", "GET", api.EndPointTypeHTTPJSON, utilsHTTP.ContentTypeApplicationJSON, nil,
		HandleAuth, nil, nil, nil, nil, 0, "default",
	)

	anAPI.NewEndPoint("RelyOn.RegisterInstallationUpdateStatus",
		"Update status of register installation to RelyOn from PGN-Partner",
		"/v1/relyon/register_installation_update_status", "POST", api.EndPointTypeHTTPJSON, utilsHTTP.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "registration_id", Type: "string", Description: "", IsMustExist: true},
		}, HandleRelyOnRegisterInstallationUpdateStatus, nil, nil, nil, []string{"RELYON.REGISTERINSTALLATIONUPDATESTATUS"}, 0, "default",
	)

	anAPI.NewEndPoint("RelyOn.RegisterInstallationUpdateType",
		"Update type of register installation to RelyOn from PGN-Partner",
		"/v1/relyon/register_installation_update_type", "POST", api.EndPointTypeHTTPJSON, utilsHTTP.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "registration_id", Type: "string", Description: "", IsMustExist: true},
			{NameId: "installation_type_id", Type: "int64", Description: "", IsMustExist: true},
		}, HandleRelyOnRegisterInstallationUpdateType, nil, nil, nil, []string{"RELYON.REGISTERINSTALLATIONUPDATESTATUS"}, 0, "default",
	)

	anAPI.NewEndPoint("RelyOn.CancelSubscriptionMileStone",
		"Request Cancel Subscription mileStone status",
		"/v1/relyon/cancel_subscription_milestone", "POST", api.EndPointTypeHTTPJSON, utilsHTTP.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "registration_id", Type: "string", Description: "", IsMustExist: true},
		}, HandleRelyOnCancelSubscriptionMileStone, nil, nil, nil, []string{"RELYON.REGISTERINSTALLATIONUPDATESTATUS"}, 0, "default",
	)

	anAPI.NewEndPoint("RelyOn.Cancel Subscription Create",
		"Request Cancel Subscription Create",
		"/v1/relyon/cancel_subscription_create", "POST", api.EndPointTypeHTTPJSON, utilsHTTP.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "registration_id", Type: "string", Description: "", IsMustExist: true},
			{NameId: "reason", Type: "array-string", Description: "", IsMustExist: true},
			{NameId: "channel", Type: "string", Description: "", IsMustExist: true},
			{NameId: "cancel_by", Type: "string", Description: "", IsMustExist: true},
		}, HandleRelyOnCancelSubscriptionCreate, nil, nil, nil, []string{"RELYON.REGISTERINSTALLATIONUPDATESTATUS"}, 0, "default",
	)

}

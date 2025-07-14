package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/lib"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"net/http"
)

var UserIdentityCard *lib.ImageObjectStorage

func UserIdentityCardUpdate(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.ParameterValues["user_id"].Value.(int64)

	_, _, err = user_management.ModuleUserManagement.User.ShouldGetById(&aepr.Log, userId)
	if err != nil {
		aepr.Log.Warnf("USER_NOT_FOUND:%d:%s", userId, err.Error())
		err = aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "USER_NOT_FOUND:%d", userId)
		return errors.Wrap(err, "error occured")
	}

	idAsString := utils.Int64ToString(userId)

	filename := idAsString + ".png"

	err = UserIdentityCard.Update(aepr, filename, "")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func UserIdentityCardUpdateFileContentBase64(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.ParameterValues["user_id"].Value.(int64)

	_, _, err = user_management.ModuleUserManagement.User.ShouldGetById(&aepr.Log, userId)
	if err != nil {
		aepr.Log.Warnf("USER_NOT_FOUND:%d:%s", userId, err.Error())
		err = aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "USER_NOT_FOUND:%d", userId)
		return errors.Wrap(err, "error occured")
	}

	idAsString := utils.Int64ToString(userId)

	filename := idAsString + ".png"

	_, fileContentBase64, err := aepr.GetParameterValueAsString("content_base64")
	if err != nil {
		return err
	}

	err = UserIdentityCard.Update(aepr, filename, fileContentBase64)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func UserIdentityCardDownloadSource(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.ParameterValues["user_id"].Value.(int64)

	_, _, err = user_management.ModuleUserManagement.User.ShouldGetById(&aepr.Log, userId)
	if err != nil {
		aepr.Log.Warnf("USER_NOT_FOUND:%d:%s", userId, err.Error())
		err = aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "USER_NOT_FOUND:%d", userId)
		return errors.Wrap(err, "error occured")
	}

	idAsString := utils.Int64ToString(userId)

	filename := idAsString + ".png"

	err = UserIdentityCard.DownloadSource(aepr, filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func UserIdentityCardDownloadBig(aepr *api.DXAPIEndPointRequest) (err error) {
	userId := aepr.ParameterValues["user_id"].Value.(int64)

	_, _, err = user_management.ModuleUserManagement.User.ShouldGetById(&aepr.Log, userId)
	if err != nil {
		aepr.Log.Warnf("USER_NOT_FOUND:%d:%s", userId, err.Error())
		err = aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "USER_NOT_FOUND:%d", userId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
		return nil
	}

	idAsString := utils.Int64ToString(userId)

	filename := idAsString + ".png"

	err = UserIdentityCard.DownloadProcessedImage(aepr, "big", filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

func SelfIdentityCardDownloadBig(aepr *api.DXAPIEndPointRequest) (err error) {
	user := aepr.LocalData["user"].(utils.JSON)
	userId := user["id"].(int64)

	idAsString := utils.Int64ToString(userId)

	filename := idAsString + ".png"

	err = UserIdentityCard.DownloadProcessedImage(aepr, "big", filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	return nil
}

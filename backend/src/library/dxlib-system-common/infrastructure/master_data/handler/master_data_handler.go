package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"net/http"
)

func CustomerSegmentList(aepr *api.DXAPIEndPointRequest) (err error) {
	_, rows, err := master_data.ModuleMasterData.CustomerSegment.SelectAll(&aepr.Log)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	data := utils.JSON{
		"data": utils.JSON{
			"list": rows,
		}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func CustomerTypeList(aepr *api.DXAPIEndPointRequest) (err error) {
	_, rows, err := master_data.ModuleMasterData.CustomerType.SelectAll(&aepr.Log)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	data := utils.JSON{
		"data": utils.JSON{
			"list": rows,
		}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

func PaymentSchemeList(aepr *api.DXAPIEndPointRequest) (err error) {
	_, rows, err := master_data.ModuleMasterData.PaymentScheme.SelectAll(&aepr.Log)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	data := utils.JSON{"data": utils.JSON{
		"list": rows,
	}}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, data)

	return nil
}

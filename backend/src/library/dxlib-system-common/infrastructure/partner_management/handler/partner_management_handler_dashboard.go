package handler

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"net/http"
)

func CMSDashboardStatsTaskTypeStatus(aepr *api.DXAPIEndPointRequest) error {
	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	_, tss, err := db.Select("partner_management.mv_task_status_summary", nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"task_status_summary": tss,
	}})

	return nil
}

func CMSDashboardStatsTaskLocationDistribution(aepr *api.DXAPIEndPointRequest) error {
	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]
	_, tld, err := db.Select("partner_management.mv_task_location_distribution", nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"task_location_distribution": tld,
	}})

	return nil
}

func CMSDashboardStatsFieldExecutorStatus(aepr *api.DXAPIEndPointRequest) error {
	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]

	_, fes, err := db.Select("partner_management.mv_field_executor_status_time_series", nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"field_executor_status": fes,
	}})
	return nil
}
func CMSDashboardStatsTaskTimeSeries(aepr *api.DXAPIEndPointRequest) error {
	db := database.Manager.Databases[base.DatabaseNameIdTaskDispatcher]

	_, dts, err := db.Select("partner_management.mv_dashboard_time_series", nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{"data": utils.JSON{
		"list": dts,
	}})
	return nil
}

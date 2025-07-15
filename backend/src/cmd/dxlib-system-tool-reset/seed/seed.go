package seed

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib-system/common"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/configuration_settings"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/json"
	"github.com/donnyhardyanto/dxlib_module/module/general"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"sync"
)

func PopulateDataPartnerTaskDispatcherMasterDataLocation(tx *database.DXDatabaseTx) (err error) {
	// INSERT_CODE_HERE
	return nil
}
func Seed() (err error) {
	err = user_management.ModuleUserManagement.AutoCreateUserSuperAdminPasswordIfNotExist(&log.Log)
	if err != nil {
		return err
	}

	wgMain := &sync.WaitGroup{}

	organizationIdOwner, err := general.ModuleGeneral.Property.GetAsInt64(&log.Log, "CONFIG.ORGANIZATION:OWNER.ID")
	if err != nil {
		return err
	}

	dbConfig := database.Manager.Databases["config"]
	var dtx1 *database.DXDatabaseTx
	dtx1, err = dbConfig.TransactionBegin(sql.LevelReadCommitted)
	if err != nil {
		return err
	}
	{
		defer dtx1.Finish(&log.Log, err)

		err = PopulateDataPartnerTaskDispatcherMasterDataLocation(dtx1)
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "USER_REGISTRATION",
			"content_type": "text/html",
			"subject":      "Pembuatan user baru <fullname>",
			"body":         "Halo <fullname>,<br><br>Terima kasih telah mendaftar di aplikasi kami. Berikut adalah informasi akun Anda:<br><br>Username: <username><br>Password: <password><br><br>Terima kasih,<br>Admin",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "USER_RESET_PASSWORD",
			"content_type": "text/html",
			"subject":      "Reset Password <fullname>",
			"body":         "Halo <fullname>,<br><br>Anda telah meminta reset password. Berikut adalah password baru Anda:<br><br>Password: <password><br><br>Terima kasih,<br>Admin",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "NEW_CONSTRUCTION_TASK_AVAILABLE",
			"content_type": "text/plain",
			"subject":      "Pekerjaan Baru <task_code> tiba",
			"body":         "Telah masuk pekerjaan baru <task_code> kategori Kontruksi yang dapat Anda ambil sekarang.",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "SUB_TASK_FIELD_SUPERVISOR_VERIFICATION_FAILED",
			"content_type": "text/plain",
			"subject":      "Verifikasi Sub Task <sub_task_code> oleh Pengawas Gagal",
			"body":         "Verifikasi Sub Task <sub_task_code> oleh Pengawas gagal. Harap diperbaiki.",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "SUB_TASK_FIELD_SUPERVISOR_VERIFICATION_SUCCESS",
			"content_type": "text/txt",
			"subject":      "Verifikasi Sub Task <sub_task_code> oleh Pengawas Sukses",
			"body":         "Verifikasi Sub Task <sub_task_code> oleh Pengawas sukses..",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "SUB_TASK_CGP_VERIFICATION_FAILED",
			"content_type": "text/plain",
			"subject":      "Verifikasi Sub Task <sub_task_code> oleh CGP Gagal",
			"body":         "Verifikasi Sub Task <sub_task_code> oleh CGP gagal. Harap diperbaiki.",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "SUB_TASK_CGP_VERIFICATION_SUCCESS",
			"content_type": "text/txt",
			"subject":      "Verifikasi Sub Task <sub_task_code> oleh CGP Sukses",
			"body":         "Verifikasi Sub Task <sub_task_code> oleh CGP sukses..",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.EMailTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "TASK_CANCEL_BY_CUSTOMER",
			"content_type": "text/plain",
			"subject":      "Task  <task_code> Dibatalkan oleh Customer",
			"body":         "Task  <task_code> dibatalkan oleh Customer",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.SMSTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "USER_REGISTRATION",
			"content_type": "text/html",
			"subject":      "Pembuatan user baru <fullname>",
			"body":         "Halo <fullname>,<br><br>Terima kasih telah mendaftar di aplikasi kami. Berikut adalah informasi akun Anda:<br><br>Username: <username><br>Password: <password><br><br>Terima kasih,<br>Admin",
		})
		if err != nil {
			return err
		}

		_, err = configuration_settings.ModuleConfigurationSettings.SMSTemplate.TxInsert(dtx1, utils.JSON{
			"nameid":       "USER_RESET_PASSWORD",
			"content_type": "text/html",
			"subject":      "Reset Password <fullname>",
			"body":         "Halo <fullname>,<br><br>Anda telah meminta reset password. Berikut adalah password baru Anda:<br><br>Password: <password><br><br>Terima kasih,<br>Admin",
		})
		if err != nil {
			return err
		}

	}

	_, roleFieldExecutor, err := user_management.ModuleUserManagement.Role.ShouldGetByUtag(&log.Log, "FIELD_EXECUTOR")
	if err != nil {
		return err
	}
	roleIdFieldExecutor := roleFieldExecutor["id"].(int64)

	err = general.ModuleGeneral.Property.SetAsInt64(&log.Log, base.ConfigRoleIdFieldExecutor, roleIdFieldExecutor)
	if err != nil {
		return err
	}
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "ACCESS.MOBILE_APP")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_GAS_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_GAS_APPLIANCE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_GAS_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_GAS_APPLIANCE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_GAS_APPLIANCE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_TAPPING_SADDLE_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_TAPPING_SADDLE_APPLIANCE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_TAPPING_SADDLE_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_TAPPING_SADDLE_APPLIANCE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_TAPPING_SADDLE_APPLIANCE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_METER_APPLIANCE_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_METER_APPLIANCE_TYPE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_METER_APPLIANCE_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_METER_APPLIANCE_TYPE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_METER_APPLIANCE_TYPE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_REGULATOR_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_REGULATOR_APPLIANCE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_REGULATOR_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_REGULATOR_APPLIANCE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_REGULATOR_APPLIANCE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_G_SIZE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_G_SIZE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_G_SIZE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_G_SIZE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_G_SIZE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_ANNOUNCEMENT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_ANNOUNCEMENT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_ANNOUNCEMENT.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_RS_CUSTOMER_SECTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUBTASK_STATUS_SUMMARY.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_CUSTOMER.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_CUSTOMER.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_CUSTOMER_METER.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_FIELD_EXECUTOR_TASK.SEARCH")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_FIELD_EXECUTOR_SUB_TASK.SEARCH")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.SCHEDULE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.PICK")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.START-WORKING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.FINISH-WORKING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.START-REWORKING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.CANCEL-REWORKING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.FINISH-REWORKING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.START-FIXING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.FINISH-FIXING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.PAUSE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.RESUME")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.CANCEL")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_USER_SUB_TASK.CANCEL_CUSTOMER")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_FILE_GROUP.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE.ASSIGN")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE_SOURCE.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE_SMALL.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE_MEDIUM.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_PICTURE_BIG.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_BERITA_ACARA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldExecutor, "MOBILE_SUB_TASK_REPORT_BERITA_ACARA.DOWNLOAD")

	_, roleFieldSupervisor, err := user_management.ModuleUserManagement.Role.ShouldGetByUtag(&log.Log, "FIELD_SUPERVISOR")
	if err != nil {
		return err
	}
	roleIdFieldSupervisor := roleFieldSupervisor["id"].(int64)

	err = general.ModuleGeneral.Property.SetAsInt64(&log.Log, base.ConfigRoleIdFieldSupervisor, roleIdFieldSupervisor)
	if err != nil {
		return err
	}
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "ACCESS.MOBILE_APP")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_USER_SUB_TASK.FINISH_VERIFY_SUCCESS")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_USER_SUB_TASK.FINISH_VERIFY_FAIL")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUBTASK_VERIFICATION_STATUS.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_FIELD_SUPERVISOR_TASK.SEARCH")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_FIELD_SUPERVISOR_SUB_TASK.SEARCH")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_FILE_GROUP.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_PICTURE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_PICTURE.ASSIGN")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_PICTURE_SOURCE.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_PICTURE_SMALL.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_PICTURE_MEDIUM.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_PICTURE_BIG.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_SUB_TASK_REPORT_BERITA_ACARA.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_TAPPING_SADDLE_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_METER_APPLIANCE_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_METER_APPLIANCE_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_G_SIZE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_G_SIZE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_GAS_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_ANNOUNCEMENT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_ANNOUNCEMENT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_ANNOUNCEMENT.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSxMustInsert(&log.Log, roleIdFieldSupervisor, "MOBILE_RS_CUSTOMER_SECTOR.LIST")

	/* Global Changeable Setting */

	err = general.ModuleGeneral.Property.SetAsJSON(&log.Log, base.SettingsTaskConstruction, map[string]any{
		base.SettingsTaskConstructionKeyFreeLengthPipeInMeter: 10,
	})
	if err != nil {
		return err
	}

	roleOMId, err := user_management.ModuleUserManagement.Role.Insert(&log.Log, utils.JSON{
		"organization_types": []string{"OWNER"},
		"nameid":             "OM",
		"name":               "OM",
		"description":        "OM",
	})
	if err != nil {
		return err
	}

	log.Log.Tracef("OM:%d", roleOMId)
	_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
		"organization_id": organizationIdOwner,
		"role_id":         roleOMId,
	})
	if err != nil {
		return err
	}

	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ACCESS.WEB_CMS")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "GAS_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "GAS_APPLIANCE.READ")

	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER_METER.LIST")

	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TAPPING_SADDLE_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "METER_APPLIANCE_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "METER_APPLIANCE_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "REGULATOR_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "REGULATOR_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "G_SIZE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "G_SIZE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ANNOUNCEMENT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ANNOUNCEMENT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER_REF.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER_REF.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER_SEGMENT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "PAYMENT_SCHEMA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "RS_CUSTOMER_SECTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "RS_CUSTOMER_SECTOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TASK_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TASK_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_EFFECTIVE_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_AREA.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_EXECUTOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_LOCATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FIELD_SUPERVISOR_EXPERTISE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_EXECUTOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_LOCATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FCM_APPLICATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FCM_APPLICATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "FCM_APPLICATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "CUSTOMER.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT_FILE_GROUP.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT_PICTURE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT_PICTURE.ASSIGN_EXISTING")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT_PICTURE.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "SUB_TASK_REPORT_PICTURE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TASK.SEARCH")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TASK.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TASK.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "USER_MESSAGE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "USER_MESSAGE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION.READ_BY_NAME")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "USER_ROLE_MEMBERSHIP.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleOMId, "ORGANIZATION.READ_BY_UTAG")

	roleSCMId, err := user_management.ModuleUserManagement.Role.Insert(&log.Log, utils.JSON{
		"organization_types": []string{"OWNER"},
		"nameid":             "SCM",
		"name":               "SCM",
		"description":        "SCM",
	})
	if err != nil {
		return err
	}
	log.Log.Tracef("SCM:%d", roleSCMId)

	_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
		"organization_id": organizationIdOwner,
		"role_id":         roleSCMId,
	})
	if err != nil {
		return err
	}

	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSCMId, "ACCESS.WEB_CMS")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSCMId, "CUSTOMER.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSCMId, "CUSTOMER.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSCMId, "CUSTOMER_METER.LIST")

	roleCGPId, err := user_management.ModuleUserManagement.Role.Insert(&log.Log, utils.JSON{
		"organization_types": []string{"OWNER"},
		"nameid":             "CGP",
		"name":               "CGP",
		"description":        "CGP",
	})
	if err != nil {
		return err
	}
	log.Log.Tracef("CGP:%d", roleCGPId)

	_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
		"organization_id": organizationIdOwner,
		"role_id":         roleCGPId,
	})
	if err != nil {
		return err
	}

	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ACCESS.WEB_CMS")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER_METER.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GAS_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GAS_APPLIANCE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GAS_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GAS_APPLIANCE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GAS_APPLIANCE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TAPPING_SADDLE_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TAPPING_SADDLE_APPLIANCE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TAPPING_SADDLE_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TAPPING_SADDLE_APPLIANCE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TAPPING_SADDLE_APPLIANCE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "METER_APPLIANCE_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "METER_APPLIANCE_TYPE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "METER_APPLIANCE_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "METER_APPLIANCE_TYPE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "METER_APPLIANCE_TYPE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "REGULATOR_APPLIANCE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "REGULATOR_APPLIANCE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "REGULATOR_APPLIANCE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "REGULATOR_APPLIANCE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "REGULATOR_APPLIANCE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "G_SIZE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "G_SIZE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.UPLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ANNOUNCEMENT.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "AREA.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "LOCATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "LOCATION.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER_REF.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER_REF.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER_REF.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER_SEGMENT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "CUSTOMER_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "PAYMENT_SCHEMA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "RS_CUSTOMER_SECTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "RS_CUSTOMER_SECTOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "RS_CUSTOMER_SECTOR.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "RS_CUSTOMER_SECTOR.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK_TYPE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK_TYPE.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK_TYPE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_TYPE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_TYPE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_TYPE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_LOCATION.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_EFFECTIVE_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_AREA.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_EXECUTOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_LOCATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FIELD_SUPERVISOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_LOCATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_EXECUTOR_EXPERTISE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_LOCATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_LOCATION.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_LOCATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_LOCATION.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_AREA.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_AREA.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_AREA.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_AREA.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_EXPERTISE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_EXPERTISE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION_SUPERVISOR_EXPERTISE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FCM_APPLICATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "FCM_APPLICATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_REPORT.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_REPORT.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_REPORT_FILE_GROUP.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_REPORT_PICTURE.DOWNLOAD")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "SUB_TASK_REPORT_PICTURE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK.SEARCH")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "TASK.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ROLE_PRIVILEGE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ROLE_PRIVILEGE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ROLE_PRIVILEGE.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ROLE_PRIVILEGE.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.CREATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.UPDATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.DELETE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.ACTIVATE")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER.SUSPEND")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER_MESSAGE.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER_MESSAGE.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION.READ")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION.READ_BY_NAME")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "ORGANIZATION.READ_BY_UTAG")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER_ROLE_MEMBERSHIP.LIST")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER_SUB_TASK.CGP_VERIFY_SUCCESS")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER_SUB_TASK.CGP_VERIFY_FAIL")
	_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPId, "USER_SUB_TASK.CGP_EDIT_AFTER_CGP_VERIFY_SUCCESS")

	_, sorAreas, err := master_data.ModuleMasterData.Area.Select(&log.Log, nil, utils.JSON{
		"type": "AREA",
	}, nil, map[string]string{"id": "asc"}, nil)
	if err != nil {
		return err
	}
	for _, sorArea := range sorAreas {
		areaCode, err := json.GetString(sorArea, "code")
		if err != nil {
			return err
		}
		areaName, err := json.GetString(sorArea, "name")
		if err != nil {
			return err
		}

		roleSORAdminSalesId, err := common.PartnerInstance.RoleCreate(&log.Log, utils.JSON{
			"organization_types": []string{"OWNER"},
			"nameid":             "SOR/ADMIN_SALES_" + areaCode,
			"name":               "SOR/Admin Sales " + areaName,
			"description":        "SOR/Admin Sales " + areaName,
			"area_code":          areaCode,
			"task_type_id":       base.TaskTypeIdConstruction,
		})
		if err != nil {
			return err
		}
		log.Log.Tracef("SOR/ADMIN_SALES_%s:%d", areaCode, roleSORAdminSalesId)

		_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
			"organization_id": organizationIdOwner,
			"role_id":         roleSORAdminSalesId,
		})
		if err != nil {
			return err
		}

		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSORAdminSalesId, "ACCESS.WEB_CMS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSORAdminSalesId, "CUSTOMER.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleSORAdminSalesId, "CUSTOMER.READ")

		roleCGPPMId, err := common.PartnerInstance.RoleCreate(&log.Log, utils.JSON{
			"organization_types": []string{"OWNER"},
			"nameid":             "CGP/PM_" + areaCode,
			"name":               "CGP/PM " + areaName,
			"description":        "CGP/PM " + areaName,
			"area_code":          areaCode,
			"task_type_id":       base.TaskTypeIdConstruction,
		})
		if err != nil {
			return err
		}
		log.Log.Tracef("CGP/PM_%s:%d", areaCode, roleCGPPMId)

		_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
			"organization_id": organizationIdOwner,
			"role_id":         roleCGPPMId,
		})
		if err != nil {
			return err
		}

		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ACCESS.WEB_CMS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GAS_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GAS_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GAS_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GAS_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GAS_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TAPPING_SADDLE_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TAPPING_SADDLE_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TAPPING_SADDLE_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TAPPING_SADDLE_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TAPPING_SADDLE_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "METER_APPLIANCE_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "METER_APPLIANCE_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "METER_APPLIANCE_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "METER_APPLIANCE_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "METER_APPLIANCE_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "REGULATOR_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "REGULATOR_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "REGULATOR_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "REGULATOR_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "REGULATOR_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "G_SIZE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "G_SIZE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.UPLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ANNOUNCEMENT.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER_REF.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER_REF.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER_REF.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER_SEGMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "CUSTOMER_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "PAYMENT_SCHEMA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "RS_CUSTOMER_SECTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "RS_CUSTOMER_SECTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "RS_CUSTOMER_SECTOR.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "RS_CUSTOMER_SECTOR.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FIELD_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_EXECUTOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION_SUPERVISOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FCM_APPLICATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "FCM_APPLICATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_REPORT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_REPORT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_REPORT_FILE_GROUP.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_REPORT_PICTURE.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "SUB_TASK_REPORT_PICTURE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK.SEARCH.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "TASK.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ROLE_PRIVILEGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ROLE_PRIVILEGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ROLE_PRIVILEGE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ROLE_PRIVILEGE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.ACTIVATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER.SUSPEND")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER_MESSAGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER_MESSAGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION.READ_BY_NAME")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "ORGANIZATION.READ_BY_UTAG")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPPMId, "USER_ROLE_MEMBERSHIP.LIST")

		roleCGPSiteCoordinatorId, err := common.PartnerInstance.RoleCreate(&log.Log, utils.JSON{
			"organization_types": []string{"OWNER"},
			"nameid":             "CGP/SITE_COORDINATOR_" + areaCode,
			"name":               "CGP/Site Coordinator " + areaName,
			"description":        "CGP/Site Coordinator " + areaName,
			"area_code":          areaCode,
			"task_type_id":       base.TaskTypeIdConstruction,
		})
		if err != nil {
			return err
		}
		log.Log.Tracef("CGP/SITE_COORDINATOR_%s:%d", areaCode, roleCGPSiteCoordinatorId)

		_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
			"organization_id": organizationIdOwner,
			"role_id":         roleCGPSiteCoordinatorId,
		})
		if err != nil {
			return err
		}

		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ACCESS.WEB_CMS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GAS_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GAS_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GAS_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GAS_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GAS_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TAPPING_SADDLE_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TAPPING_SADDLE_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TAPPING_SADDLE_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TAPPING_SADDLE_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TAPPING_SADDLE_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "METER_APPLIANCE_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "METER_APPLIANCE_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "METER_APPLIANCE_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "METER_APPLIANCE_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "METER_APPLIANCE_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "REGULATOR_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "REGULATOR_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "REGULATOR_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "REGULATOR_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "REGULATOR_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "G_SIZE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "G_SIZE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.UPLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ANNOUNCEMENT.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER_REF.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER_REF.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER_REF.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER_SEGMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "CUSTOMER_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "PAYMENT_SCHEMA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "RS_CUSTOMER_SECTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "RS_CUSTOMER_SECTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "RS_CUSTOMER_SECTOR.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "RS_CUSTOMER_SECTOR.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FIELD_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_EXECUTOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION_SUPERVISOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FCM_APPLICATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "FCM_APPLICATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_REPORT.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_REPORT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_REPORT_FILE_GROUP.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_REPORT_PICTURE.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "SUB_TASK_REPORT_PICTURE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK.SEARCH")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "TASK.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ROLE_PRIVILEGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ROLE_PRIVILEGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ROLE_PRIVILEGE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ROLE_PRIVILEGE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.ACTIVATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER.SUSPEND")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER_MESSAGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER_MESSAGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION.READ_BY_NAME")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "ORGANIZATION.READ_BY_UTAG")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPSiteCoordinatorId, "USER_ROLE_MEMBERSHIP.LIST")

		roleCGPTraceabilityId, err := common.PartnerInstance.RoleCreate(&log.Log, utils.JSON{
			"organization_types": []string{"OWNER"},
			"nameid":             "CGP/TRACEABILITY_" + areaCode,
			"name":               "CGP/Traceability " + areaName,
			"description":        "CGP/Traceability " + areaName,
			"area_code":          areaCode,
			"task_type_id":       base.TaskTypeIdConstruction,
		})

		log.Log.Tracef("CGP/TRACEABILITY_%s:%d", areaCode, roleCGPTraceabilityId)

		_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
			"organization_id": organizationIdOwner,
			"role_id":         roleCGPTraceabilityId,
		})
		if err != nil {
			return err
		}

		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ACCESS.WEB_CMS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER.UPLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GAS_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GAS_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GAS_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GAS_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GAS_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TAPPING_SADDLE_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TAPPING_SADDLE_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TAPPING_SADDLE_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TAPPING_SADDLE_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TAPPING_SADDLE_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "METER_APPLIANCE_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "METER_APPLIANCE_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "METER_APPLIANCE_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "METER_APPLIANCE_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "METER_APPLIANCE_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "REGULATOR_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "REGULATOR_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "REGULATOR_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "REGULATOR_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "REGULATOR_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "G_SIZE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "G_SIZE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "G_SIZE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "G_SIZE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "G_SIZE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ANNOUNCEMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ANNOUNCEMENT.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ANNOUNCEMENT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ANNOUNCEMENT.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ANNOUNCEMENT.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ANNOUNCEMENT.UPLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_REF.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_REF.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_REF.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_REF.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_REF.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_SEGMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "CUSTOMER_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "PAYMENT_SCHEMA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "RS_CUSTOMER_SECTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "RS_CUSTOMER_SECTOR.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "RS_CUSTOMER_SECTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "RS_CUSTOMER_SECTOR.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "RS_CUSTOMER_SECTOR.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EXPERTISE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_EXECUTOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FIELD_SUPERVISOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_EXECUTOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION_SUPERVISOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FCM_APPLICATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FCM_APPLICATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FCM_APPLICATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FCM_APPLICATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "FCM_APPLICATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK.SCHEDULE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK.ASSIGN")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK.REPLACE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_SUB_TASK.CGP_VERIFY_SUCCESS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_SUB_TASK.CGP_VERIFY_FAIL")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_SUB_TASK.CGP_EDIT_AFTER_CGP_VERIFY_SUCCESS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_REPORT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_REPORT_FILE_GROUP.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_REPORT_PICTURE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_REPORT_PICTURE.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "SUB_TASK_REPORT_PICTURE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK.SEARCH.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "TASK.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ROLE_PRIVILEGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ROLE_PRIVILEGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ROLE_PRIVILEGE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ROLE_PRIVILEGE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.ACTIVATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER.SUSPEND")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_MESSAGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_MESSAGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.READ_BY_NAME")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "ORGANIZATION.READ_BY_UTAG")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_ROLE_MEMBERSHIP.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_ROLE_MEMBERSHIP.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPTraceabilityId, "USER_ROLE_MEMBERSHIP.DELETE")

		roleCGPAdminConstructionId, err := common.PartnerInstance.RoleCreate(&log.Log, utils.JSON{
			"organization_types": []string{"OWNER"},
			"nameid":             "CGP/ADMIN_CONSTRUCTION_" + areaCode,
			"name":               "CGP/Admin Coordinator " + areaName,
			"description":        "CGP/Admin Coordinator " + areaName,
			"area_code":          areaCode,
			"task_type_id":       base.TaskTypeIdConstruction,
		})
		if err != nil {
			return err
		}
		log.Log.Tracef("CGP/ADMIN_CONSTRUCTION_%s:%d", areaCode, roleCGPAdminConstructionId)

		_, err = user_management.ModuleUserManagement.OrganizationRoles.Insert(&log.Log, utils.JSON{
			"organization_id": organizationIdOwner,
			"role_id":         roleCGPAdminConstructionId,
		})
		if err != nil {
			return err
		}

		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ACCESS.WEB_CMS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER.UPLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GAS_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GAS_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GAS_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GAS_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GAS_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TAPPING_SADDLE_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TAPPING_SADDLE_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TAPPING_SADDLE_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TAPPING_SADDLE_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TAPPING_SADDLE_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "METER_APPLIANCE_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "METER_APPLIANCE_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "METER_APPLIANCE_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "METER_APPLIANCE_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "METER_APPLIANCE_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "REGULATOR_APPLIANCE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "REGULATOR_APPLIANCE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "REGULATOR_APPLIANCE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "REGULATOR_APPLIANCE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "REGULATOR_APPLIANCE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "G_SIZE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "G_SIZE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "G_SIZE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "G_SIZE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "G_SIZE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ANNOUNCEMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ANNOUNCEMENT.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ANNOUNCEMENT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ANNOUNCEMENT.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ANNOUNCEMENT.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ANNOUNCEMENT.UPLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "GENERAL_SETTINGS_TASK_CONSTRUCTION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_REF.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_REF.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_REF.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_REF.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_REF.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_SEGMENT.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "CUSTOMER_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "PAYMENT_SCHEMA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "RS_CUSTOMER_SECTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "RS_CUSTOMER_SECTOR.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "RS_CUSTOMER_SECTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "RS_CUSTOMER_SECTOR.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "RS_CUSTOMER_SECTOR.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_TYPE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_TYPE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_TYPE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_TYPE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_TYPE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EFFECTIVE_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_LOCATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_AREA.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EXPERTISE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_EXECUTOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_EFFECTIVE_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_EFFECTIVE_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FIELD_SUPERVISOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_EXECUTOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_LOCATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_LOCATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_LOCATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_LOCATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_AREA.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_AREA.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_AREA.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_AREA.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_EXPERTISE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_EXPERTISE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_EXPERTISE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION_SUPERVISOR_EXPERTISE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FCM_APPLICATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FCM_APPLICATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FCM_APPLICATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FCM_APPLICATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "FCM_APPLICATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK.SCHEDULE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK.ASSIGN")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK.REPLACE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_SUB_TASK.CGP_VERIFY_SUCCESS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_SUB_TASK.CGP_VERIFY_FAIL")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_SUB_TASK.CGP_EDIT_AFTER_CGP_VERIFY_SUCCESS")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_REPORT.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_REPORT_FILE_GROUP.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_REPORT_PICTURE.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_REPORT_PICTURE.DOWNLOAD")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "SUB_TASK_REPORT_PICTURE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK.SEARCH.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK.LIST.PER_CONSTRUCTION_AREA")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "TASK.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ROLE_PRIVILEGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ROLE_PRIVILEGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ROLE_PRIVILEGE.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ROLE_PRIVILEGE.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.ACTIVATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER.SUSPEND")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_MESSAGE.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_MESSAGE.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.READ")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.UPDATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.DELETE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.READ_BY_NAME")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "ORGANIZATION.READ_BY_UTAG")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_ROLE_MEMBERSHIP.LIST")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_ROLE_MEMBERSHIP.CREATE")
		_ = user_management.ModuleUserManagement.RolePrivilegeSWgMustInsert(wgMain, &log.Log, roleCGPAdminConstructionId, "USER_ROLE_MEMBERSHIP.DELETE")
	}
	err = general.ModuleGeneral.Property.SetAsString(&log.Log, "EMAIL-TEMPLATE-CREATE-USER", "Halo <name>,\n\nSelamat datang di PGN Partner! Kami sangat senang menyambut Anda. kami telah menyiapkan akun khusus untuk Anda.\n\nBerikut adalah detail login akun Anda:\n\nUsername: <username>\nPassword: <password>\n\nSilakan gunakan informasi ini untuk masuk ke aplikasi dan segera mengganti password Anda demi keamanan akun.\n\nJika Anda mengalami kesulitan atau membutuhkan bantuan, jangan ragu untuk menghubungi tim dukungan kami.\n\nTerima kasih telah bergabung dengan kami. Kami berharap Anda mendapatkan pengalaman yang luar biasa menggunakan PGN Partner.\n\nSalam hangat,\n\nTeam PGN Partner")
	if err != nil {
		return err
	}
	err = general.ModuleGeneral.Property.SetAsString(&log.Log, "EMAIL-TEMPLATE-RESET-PASSWORD", "Halo <name>,\n\nAkun Anda telah berhasil direset. Berikut adalah detail login terbaru Anda:\n\nUsername: <username>\nPassword: <password>\n\nKami sarankan untuk segera login dan mengganti password ini untuk keamanan akun Anda. Anda dapat melakukannya dengan mengikuti langkah-langkah di pengaturan akun setelah login.\n\nJika Anda tidak meminta reset ini, atau jika ada aktivitas yang mencurigakan, segera hubungi tim dukungan kami.\n\nTerima kasih,\n\nTeam PGN Partner")
	if err != nil {
		return err
	}

	log.Log.Info("Seed task-dispatcher completed")

	dbAuditLog := database.Manager.Databases["auditlog"]
	var dtx3 *database.DXDatabaseTx
	dtx3, err = dbAuditLog.TransactionBegin(sql.LevelReadCommitted)
	if err != nil {
		return err
	}
	{
		defer dtx3.Finish(&log.Log, err)
	}

	wgMain.Wait()
	return nil
}

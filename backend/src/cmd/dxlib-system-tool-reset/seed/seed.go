package seed

import (
	"database/sql"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/configuration_settings"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/general"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"sync"
)

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

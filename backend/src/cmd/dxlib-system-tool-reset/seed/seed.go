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

	_, err = general.ModuleGeneral.Property.GetAsInt64(&log.Log, "CONFIG.ORGANIZATION:OWNER.ID")
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

	err = general.ModuleGeneral.Property.SetAsString(&log.Log, "EMAIL-TEMPLATE-CREATE-USER", "Halo <name>,\n\nSelamat datang di DXLib System! Kami sangat senang menyambut Anda. kami telah menyiapkan akun khusus untuk Anda.\n\nBerikut adalah detail login akun Anda:\n\nUsername: <username>\nPassword: <password>\n\nSilakan gunakan informasi ini untuk masuk ke aplikasi dan segera mengganti password Anda demi keamanan akun.\n\nJika Anda mengalami kesulitan atau membutuhkan bantuan, jangan ragu untuk menghubungi tim dukungan kami.\n\nTerima kasih telah bergabung dengan kami. Kami berharap Anda mendapatkan pengalaman yang luar biasa menggunakan DXLib System.\n\nSalam hangat,\n\nTeam DXLib System")
	if err != nil {
		return err
	}
	err = general.ModuleGeneral.Property.SetAsString(&log.Log, "EMAIL-TEMPLATE-RESET-PASSWORD", "Halo <name>,\n\nAkun Anda telah berhasil direset. Berikut adalah detail login terbaru Anda:\n\nUsername: <username>\nPassword: <password>\n\nKami sarankan untuk segera login dan mengganti password ini untuk keamanan akun Anda. Anda dapat melakukannya dengan mengikuti langkah-langkah di pengaturan akun setelah login.\n\nJika Anda tidak meminta reset ini, atau jika ada aktivitas yang mencurigakan, segera hubungi tim dukungan kami.\n\nTerima kasih,\n\nTeam DXLib System")
	if err != nil {
		return err
	}

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
	log.Log.Info("Seed completed")

	return nil
}

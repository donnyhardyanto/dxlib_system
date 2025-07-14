package upload_data

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"strings"
)

func (am *UploadData) validateUserData(l *log.DXLog, data utils.JSON) utils.JSON {
	data["row_message"] = ""
	if orgCode, ok := data["organization_code"].(string); !ok {
		setErrorMessage(data, "organization_code unknown")
	} else {
		trimmed := strings.TrimSpace(strings.Trim(orgCode, "'"))
		data["organization_code"] = trimmed
	}

	if _, _, err := user_management.ModuleUserManagement.User.ShouldSelectOne(l, utils.JSON{
		"loginid": data["loginid"],
	}, nil, nil); err == nil {
		setErrorMessage(data, "loginid already exist")
	}
	if _, ok := data["identity_number"].(string); ok {
		if _, _, err := user_management.ModuleUserManagement.User.ShouldSelectOne(l, utils.JSON{
			"identity_number": data["identity_number"],
		}, nil, nil); err == nil {
			setErrorMessage(data, "identity_number already exist")
		}
	}
	if _, _, err := user_management.ModuleUserManagement.User.ShouldSelectOne(l, utils.JSON{
		"email": data["email"],
	}, nil, nil); err == nil {
		setErrorMessage(data, "email already exist")
	}
	if _, org, err := user_management.ModuleUserManagement.Organization.ShouldSelectOne(l, utils.JSON{
		"code": data["organization_code"],
	}, nil, nil); err != nil {
		setErrorMessage(data, fmt.Sprintf("organization_code '%v' not found", data["organization_code"]))
	} else {
		data["organization_id"] = org["id"]
	}

	// check role
	_, role, err := user_management.ModuleUserManagement.Role.ShouldSelectOne(l, utils.JSON{
		"nameid": data["role"],
	}, nil, nil)
	if err != nil {
		setErrorMessage(data, fmt.Sprintf("role '%v' not found", data["role"]))
	} else {
		if id, ok := role["id"].(int64); ok {
			data["role_id"] = id
		}
	}

	return data
}

func (am *UploadData) DoUploadUserCreate(l *log.DXLog, data utils.JSON) (id int64, err error) {
	data = am.validateUserData(l, data)
	fmt.Printf("%v", data)
	id, err = am.User.Insert(l, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (am *UploadData) DoUploadUserUpdate(l *log.DXLog, id int64, data utils.JSON) (err error) {
	_, err = am.User.UpdateOne(l, id, data)
	if err != nil {
		return err
	} else {
		data = am.validateUserData(l, data)

		return err
	}
	return nil
}

func (am *UploadData) DoUserCreate(l *log.DXLog, userData utils.JSON) (err error, userId int64) {
	fmt.Println("========= DoUserCreate ==========")
	userData = am.validateUserData(l, userData)
	if userData["row_status"] == "ERROR" {
		return errors.Errorf("can not process: %v", userData["row_message"]), userId
	}
	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]
	fmt.Println("========= start db transaction ==========")
	txErr := db.Tx(l, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err error) {
		userId, err = user_management.ModuleUserManagement.User.TxInsert(tx,
			utils.JSON{
				"loginid":              userData["loginid"],
				"fullname":             userData["fullname"],
				"email":                userData["email"],
				"phonenumber":          userData["phonenumber"],
				"identity_number":      userData["identity_number"],
				"identity_type":        userData["identity_type"],
				"must_change_password": true,
				"status":               "ACTIVE",
			})
		if err != nil {
			fmt.Println("ada error ==========")
			return err
		}
		// assign to organization
		data := utils.JSON{
			"user_id":         userId,
			"organization_id": userData["organization_id"],
		}
		if userData["membership_number"] != nil {
			data["membership_number"] = userData["membership_number"]
		}
		_, err = user_management.ModuleUserManagement.UserOrganizationMembership.TxInsert(
			tx, data)
		if err != nil {
			return err
		}
		// assign to role
		userRoleMembershipId, err := user_management.ModuleUserManagement.UserRoleMembership.TxInsert(
			tx, utils.JSON{
				"user_id":         userId,
				"organization_id": userData["organization_id"],
				"role_id":         userData["role_id"],
			})
		if err != nil {
			return err
		}

		//assign expertise
		if data1, ok := userData["expertise"].(string); ok {
			intResult := SplitAndParseInt64s(data1)
			if userData["role"] == "FIELD_EXECUTOR" {
				for _, item := range intResult {
					_, err := partner_management.ModulePartnerManagement.FieldExecutorExpertise.TxInsert(tx, utils.JSON{
						"user_role_membership_id": userRoleMembershipId,
						"sub_task_type_id":        item,
					})
					if err != nil {
						return err
					}
				}
			}
			if userData["role"] == "FIELD_SUPERVISOR" {
				for _, item := range intResult {
					_, err := partner_management.ModulePartnerManagement.FieldSupervisorExpertise.TxInsert(tx, utils.JSON{
						"user_role_membership_id": userRoleMembershipId,
						"sub_task_type_id":        item,
					})
					if err != nil {
						return err
					}
				}
			}

		}

		//assign area
		if data1, ok := userData["area_code"].(string); ok {
			result := SplitAndTrimStrings(data1)
			if userData["role"] == "FIELD_EXECUTOR" {
				for _, item := range result {
					_, err := partner_management.ModulePartnerManagement.FieldExecutorArea.TxInsert(tx, utils.JSON{
						"user_role_membership_id": userRoleMembershipId,
						"area_code":               item,
					})
					if err != nil {
						return err
					}
				}
			}
			if userData["role"] == "FIELD_SUPERVISOR" {
				for _, item := range result {
					_, err := partner_management.ModulePartnerManagement.FieldSupervisorArea.TxInsert(tx, utils.JSON{
						"user_role_membership_id": userRoleMembershipId,
						"area_code":               item,
					})
					if err != nil {
						return err
					}
				}
			}

		}

		//assign location
		if data1, ok := userData["location"].(string); ok {
			result := SplitAndTrimStrings(data1)
			if userData["role"] == "FIELD_EXECUTOR" {
				for _, item := range result {
					_, err := partner_management.ModulePartnerManagement.FieldExecutorLocation.TxInsert(tx, utils.JSON{
						"user_role_membership_id": userRoleMembershipId,
						"location_code":           item,
					})
					if err != nil {
						return err
					}
				}
			}
			if userData["role"] == "FIELD_SUPERVISOR" {
				for _, item := range result {
					_, err := partner_management.ModulePartnerManagement.FieldSupervisorLocation.TxInsert(tx, utils.JSON{
						"user_role_membership_id": userRoleMembershipId,
						"location_code":           item,
					})
					if err != nil {
						return err
					}
				}
			}

		}
		// Create password
		defaultPassword := "#pgN2025"
		err = user_management.ModuleUserManagement.TxUserPasswordCreate(tx, userId, defaultPassword)
		if err != nil {
			return err
		}

		// Call post-create hooks if they exist
		//if user_management.ModuleUserManagement.OnUserAfterCreate != nil {
		//	_, user, err := user_management.ModuleUserManagement.User.TxSelectOne(tx, utils.JSON{
		//		"id": userId,
		//	}, nil)
		//	if err != nil {
		//		return err
		//	}
		//	err = user_management.ModuleUserManagement.OnUserAfterCreate(nil, tx, user, defaultPassword) // Pass nil for aepr since we don't have it in this context
		//	if err != nil {
		//		return err
		//	}
		//}

		return nil
	})

	if txErr != nil {
		fmt.Println("====error===")
		fmt.Printf("%s", txErr.Error())
		// subTaskId is already set appropriately (either to a valid value or -1)
		return txErr, userId
	} else {
		fmt.Println("====no error===")
	}
	return nil, userId
}

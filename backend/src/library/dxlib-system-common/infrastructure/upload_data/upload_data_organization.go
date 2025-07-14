package upload_data

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib_module/module/user_management"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/master_data"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/partner_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"strconv"
	"strings"
)

func splitField(data map[string]interface{}, fieldName string) []string {
	if val, ok := data[fieldName].(string); ok {
		return strings.Split(val, ",")
	}
	return []string{} // or handle error as appropriate
}

func ValidateOrganizationData(l *log.DXLog, data utils.JSON) utils.JSON {
	data["row_message"] = nil
	_, ok := data["code"].(string)
	if !ok {
		setErrorMessage(data, "code is required")
	}

	if _, _, err := user_management.ModuleUserManagement.Organization.ShouldSelectOne(l, utils.JSON{
		"code": data["code"],
	}, nil, nil); err == nil {

		setErrorMessage(data, fmt.Sprintf("organization code '%v' already exist", data["code"]))
	}
	if _, _, err := user_management.ModuleUserManagement.Organization.ShouldSelectOne(l, utils.JSON{
		"name": data["name"],
	}, nil, nil); err == nil {
		setErrorMessage(data, fmt.Sprintf("organization name '%v' already exist", data["name"]))
	}

	if parentCode, exists := data["parent_code"]; exists && parentCode != nil && parentCode != "" {
		if _, row, err := user_management.ModuleUserManagement.Organization.ShouldSelectOne(l, utils.JSON{
			"code": parentCode,
		}, nil, nil); err == nil {
			if id, ok := row["id"]; ok {
				data["parent_id"] = id
			} else {
				setErrorMessage(data, "parent organization found but missing ID")
			}
		} else {
			setErrorMessage(data, fmt.Sprintf("parent_code '%v' not found", parentCode))
		}
	}

	//check executor area
	items := splitField(data, "field_executor_area")
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		_, _, err := master_data.ModuleMasterData.Area.ShouldSelectOne(l, utils.JSON{
			"code": trimmed,
		}, nil, nil)
		if err != nil {
			setErrorMessage(data, fmt.Sprintf("executor area %s not found", trimmed))
		}
	}

	//check supervisor area
	items = splitField(data, "field_supervisor_area")
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		_, _, err := master_data.ModuleMasterData.Area.ShouldSelectOne(l, utils.JSON{
			"code": trimmed,
		}, nil, nil)
		if err != nil {
			setErrorMessage(data, fmt.Sprintf("supervisor area %s not found", trimmed))
		}
	}

	//check executor location
	items = splitField(data, "field_executor_location")
	//Create a slice to store trimmed values
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		_, _, err := master_data.ModuleMasterData.Location.ShouldSelectOne(l, utils.JSON{
			"code": trimmed,
		}, nil, nil)
		if err != nil {
			setErrorMessage(data, fmt.Sprintf("executor location %s not found", trimmed))
		}
	}

	//check supervisor location
	items = splitField(data, "field_supervisor_location")
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		_, _, err := master_data.ModuleMasterData.Location.ShouldSelectOne(l, utils.JSON{
			"code": trimmed,
		}, nil, nil)
		if err != nil {
			setErrorMessage(data, fmt.Sprintf("supervisor location %s not found", trimmed))
			continue
		}
	}

	//check executor expertise
	items = splitField(data, "field_executor_expertise")

	//Create a slice to store trimmed values
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		subTaskTypeId, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			setErrorMessage(data, fmt.Sprintf("field_executor_expertise must be integer", trimmed))
			continue
		}
		_, _, err = task_management.ModuleTaskManagement.SubTaskType.ShouldSelectOne(l, utils.JSON{
			"id": subTaskTypeId,
		}, nil, nil)
		if err != nil {
			setErrorMessage(data, fmt.Sprintf("field_supervisor_expertise  %s not found", trimmed))
		}
	}

	return data
}

func (am *UploadData) DoUploadOrganizationCreate(l *log.DXLog, data utils.JSON) (id int64, err error) {
	data = ValidateOrganizationData(l, data)

	fmt.Printf("%v", data)
	id, err = am.Organization.Insert(l, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (am *UploadData) DoUploadOrganizationUpdate(l *log.DXLog, id int64, data utils.JSON) (err error) {
	data = ValidateOrganizationData(l, data)
	_, err = am.Organization.UpdateOne(l, id, data)
	if err != nil {
		return err
	}
	return nil
}

func (am *UploadData) DoOrganizationCreate(l *log.DXLog, data utils.JSON) (err error, organizationId int64) {
	data = ValidateOrganizationData(l, data)
	if data["row_status"] == "ERROR" {
		return errors.Errorf("cannot process: %v", data["row_message"]), organizationId
	}

	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]

	fmt.Println("========= start db transaction ==========")
	txErr := db.Tx(l, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err error) {
		organizationId, err := user_management.ModuleUserManagement.Organization.TxInsert(tx,
			utils.JSON{
				"code":        data["code"],
				"name":        data["name"],
				"email":       data["email"],
				"phonenumber": data["phonenumber"],
				"type":        data["type"],
				"address":     data["address"],
				"status":      "ACTIVE",
			})
		if err != nil {
			fmt.Println("ada error ==========")
			return err
		}
		fmt.Printf("==== %v ====", organizationId)
		// assign to role
		var result []string
		if rolesData, ok := data["roles"].(string); ok {
			result = SplitAndTrimStrings(rolesData)
			for _, item := range result {
				_, r, _ := user_management.ModuleUserManagement.Role.ShouldSelectOne(l, utils.JSON{
					"nameid": item,
				}, nil, nil)
				organizationRoleId, err := user_management.ModuleUserManagement.OrganizationRoles.TxInsert(tx, utils.JSON{
					"organization_id": organizationId,
					"role_id":         r["id"],
				})
				if err == nil {
					if item == "FIELD_EXECUTOR" {
						if data1, ok := data["field_executor_area"].(string); ok {
							rst1 := SplitAndTrimStrings(data1)
							for _, item1 := range rst1 {
								_, err1 := partner_management.ModulePartnerManagement.OrganizationExecutorArea.TxInsert(tx, utils.JSON{
									"organization_role_id": organizationRoleId,
									"area_code":            item1,
								})
								if err1 != nil {
									continue
								}
							}
						}
						if data1, ok := data["field_executor_location"].(string); ok {
							rst1 := SplitAndTrimStrings(data1)
							for _, item1 := range rst1 {
								_, err1 := partner_management.ModulePartnerManagement.OrganizationExecutorLocation.TxInsert(tx, utils.JSON{
									"organization_role_id": organizationRoleId,
									"location_code":        item1,
								})
								if err1 != nil {
									continue
								}
							}
						}
						if data1, ok := data["field_executor_expertise"].(string); ok {
							rst1 := SplitAndParseInt64s(data1)
							for _, item1 := range rst1 {
								_, err1 := partner_management.ModulePartnerManagement.OrganizationExecutorExpertise.TxInsert(tx, utils.JSON{
									"organization_role_id": organizationRoleId,
									"sub_task_type_id":     item1,
								})
								if err1 != nil {
									continue
								}
							}
						}
					}
					if item == "FIELD_SUPERVISOR" {
						if data1, ok := data["field_supervisor_area"].(string); ok {
							rst1 := SplitAndTrimStrings(data1)
							for _, item1 := range rst1 {
								_, err1 := partner_management.ModulePartnerManagement.OrganizationSupervisorArea.TxInsert(tx, utils.JSON{
									"organization_role_id": organizationRoleId,
									"area_code":            item1,
								})
								if err1 != nil {
									continue
								}
							}
						}
						if data1, ok := data["field_supervisor_location"].(string); ok {
							rst1 := SplitAndTrimStrings(data1)
							for _, item1 := range rst1 {
								_, err1 := partner_management.ModulePartnerManagement.OrganizationSupervisorLocation.TxInsert(tx, utils.JSON{
									"organization_role_id": organizationRoleId,
									"location_code":        item1,
								})
								if err1 != nil {
									continue
								}
							}
						}
						if data1, ok := data["field_supervisor_expertise"].(string); ok {
							rst1 := SplitAndParseInt64s(data1)
							for _, item1 := range rst1 {
								_, err1 := partner_management.ModulePartnerManagement.OrganizationSupervisorExpertise.TxInsert(tx, utils.JSON{
									"organization_role_id": organizationRoleId,
									"sub_task_type_id":     item1,
								})
								if err1 != nil {
									continue
								}
							}
						}
					}

				} else {
					continue
				}
			}
		}

		//userRoleMembershipId, err := user_management.ModuleUserManagement.UserRoleMembership.TxInsert(
		//	tx, utils.JSON{
		//		"user_id":         organizationId,
		//		"organization_id": organizationId,
		//		"role_id":         roleId,
		//	})
		//if err != nil {
		//	return err
		//}

		//assign expertise
		//items := strings.Split(organizationData["expertise"].(string), ",")
		//intResult := make([]int, 0, len(items))
		//for _, item := range items {
		//	trimmed := strings.TrimSpace(item)
		//	var val int
		//	_, err := fmt.Sscanf(trimmed, "%d", &val)
		//	if err != nil {
		//		continue
		//	}
		//	intResult = append(intResult, val)
		//}
		//
		//if organizationData["role"] == "FIELD_EXECUTOR" {
		//	for _, item := range intResult {
		//		_, err := partner_management.ModulePartnerManagement.FieldExecutorExpertise.TxInsert(tx, utils.JSON{
		//			"user_role_membership_id": userRoleMembershipId,
		//			"sub_task_type_id":        item,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
		//if organizationData["role"] == "FIELD_SUPERVISOR" {
		//	for _, item := range intResult {
		//		_, err := partner_management.ModulePartnerManagement.FieldSupervisorExpertise.TxInsert(tx, utils.JSON{
		//			"user_role_membership_id": userRoleMembershipId,
		//			"sub_task_type_id":        item,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
		//
		////assign area
		//items = strings.Split(organizationData["area_code"].(string), ",")
		//// Create a slice to store trimmed values
		//result := make([]string, 0, len(items))
		//for _, item := range items {
		//	trimmed := strings.TrimSpace(item)
		//	result = append(result, trimmed)
		//}
		//if organizationData["role"] == "FIELD_EXECUTOR" {
		//	for _, item := range result {
		//		_, err := partner_management.ModulePartnerManagement.FieldExecutorArea.TxInsert(tx, utils.JSON{
		//			"user_role_membership_id": userRoleMembershipId,
		//			"area_code":               item,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
		//if organizationData["role"] == "FIELD_SUPERVISOR" {
		//	for _, item := range result {
		//		_, err := partner_management.ModulePartnerManagement.FieldSupervisorArea.TxInsert(tx, utils.JSON{
		//			"user_role_membership_id": userRoleMembershipId,
		//			"area_code":               item,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}

		//assign location

		//items = strings.Split(organizationData["location"].(string), ",")
		//// Create a slice to store trimmed values
		//result = make([]string, 0, len(items))
		//for _, item := range items {
		//	trimmed := strings.TrimSpace(item)
		//	result = append(result, trimmed)
		//
		//}
		//if organizationData["role"] == "FIELD_EXECUTOR" {
		//	for _, item := range result {
		//		_, err := partner_management.ModulePartnerManagement.FieldExecutorLocation.TxInsert(tx, utils.JSON{
		//			"user_role_membership_id": userRoleMembershipId,
		//			"location_code":           item,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
		//if organizationData["role"] == "FIELD_SUPERVISOR" {
		//	for _, item := range result {
		//		_, err := partner_management.ModulePartnerManagement.FieldSupervisorLocation.TxInsert(tx, utils.JSON{
		//			"user_role_membership_id": userRoleMembershipId,
		//			"location_code":           item,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}

		return nil
	})

	if txErr != nil {
		fmt.Println("====error===")
		fmt.Printf("%s", txErr.Error())
		// subTaskId is already set appropriately (either to a valid value or -1)
		return txErr, organizationId
	} else {
		fmt.Println("====no error===")
	}
	return nil, organizationId
}

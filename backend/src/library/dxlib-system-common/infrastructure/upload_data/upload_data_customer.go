package upload_data

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/database"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"strings"
)

func (am *UploadData) validateCustomerData(l *log.DXLog, data utils.JSON) utils.JSON {
	newData := utils.JSON{}
	newData["row_message"] = ""

	if val, exists := data["row_status"]; exists && val != nil {
		newData["row_status"] = val
	}

	if val, exists := data["registration_number"]; exists && val != nil {
		newData["registration_number"] = val
		if _, _, err := task_management.ModuleTaskManagement.Customer.ShouldSelectOne(l, utils.JSON{
			"registration_number": newData["registration_number"],
		}, nil, nil); err == nil {
			setErrorMessage(newData, "registration_number already exist")
		}
	}

	if val, exists := data["customer_number"]; exists && val != nil {
		newData["customer_number"] = val
		if _, _, err := task_management.ModuleTaskManagement.Customer.ShouldSelectOne(l, utils.JSON{
			"customer_number": newData["customer_number"],
		}, nil, nil); err == nil {
			setErrorMessage(newData, "customer_number already exist")
		}
	}
	if val, exists := data["fullname"]; exists && val != nil {
		newData["fullname"] = val
	}
	if val, exists := data["status"]; exists && val != nil {
		newData["status"] = val
	}
	if val, exists := data["email"]; exists && val != nil {
		newData["email"] = val
	}
	if val, exists := data["phonenumber"]; exists && val != nil {
		newData["phonenumber"] = val
	}
	if val, exists := data["korespondensi_media"]; exists && val != nil {
		newData["korespondensi_media"] = val
	}
	if val, exists := data["identity_type"]; exists && val != nil {
		newData["identity_type"] = val
	}
	if val, exists := data["identity_number"]; exists && val != nil {
		newData["identity_number"] = val
	}
	if val, exists := data["npwp"]; exists && val != nil {
		newData["npwp"] = val
	}
	if val, exists := data["customer_segment_code"]; exists && val != nil {
		newData["customer_segment_code"] = val
	}
	if val, exists := data["customer_type_code"]; exists && val != nil {
		newData["customer_type_code"] = val
	}
	if val, exists := data["jenis_anggaran"]; exists && val != nil {
		newData["jenis_anggaran"] = val
	}
	if val, exists := data["rs_customer_sector_code"]; exists && val != nil {
		newData["rs_customer_sector_code"] = val
	}
	if val, exists := data["sales_area_code"]; exists && val != nil {
		newData["sales_area_code"] = val
	}
	if val, exists := data["latitude"]; exists && val != nil {
		newData["latitude"] = val
	}
	if val, exists := data["longitude"]; exists && val != nil {
		newData["longitude"] = val
	}
	if val, exists := data["geom"]; exists && val != nil {
		newData["geom"] = val
	}
	if val, exists := data["address_name"]; exists && val != nil {
		newData["address_name"] = val
	}
	if val, exists := data["address_street"]; exists && val != nil {
		newData["address_street"] = val
	}
	if val, exists := data["address_rt"]; exists && val != nil {
		newData["address_rt"] = val
	}
	if val, exists := data["address_rw"]; exists && val != nil {
		newData["address_rw"] = val
	}
	if val, exists := data["address_kelurahan_location_code"]; exists && val != nil {
		newData["address_kelurahan_location_code"] = val
	}
	if val, exists := data["address_kecamatan_location_code"]; exists && val != nil {
		newData["address_kecamatan_location_code"] = val
	}
	if val, exists := data["address_kabupaten_location_code"]; exists && val != nil {
		newData["address_kabupaten_location_code"] = val
	}
	if val, exists := data["address_province_location_code"]; exists && val != nil {
		newData["address_province_location_code"] = val
	}
	if val, exists := data["address_postal_code"]; exists && val != nil {
		newData["address_postal_code"] = val
	}
	if val, exists := data["register_at"]; exists && val != nil {
		newData["register_at"] = val
	} else {
		if val, exists := data["register_timestamp"]; exists && val != nil {
			newData["register_at"] = val
		}
	}
	if val, exists := data["jenis_bangunan"]; exists && val != nil {
		newData["jenis_bangunan"] = val
	}
	if val, exists := data["program_pelanggan"]; exists && val != nil {
		newData["program_pelanggan"] = val
	}
	if val, exists := data["payment_scheme_code"]; exists && val != nil {
		newData["payment_scheme_code"] = val
	}
	if val, exists := data["kategory_wilayah"]; exists && val != nil {
		newData["kategory_wilayah"] = val
	}
	if val, exists := data["cancellation_submission_status"]; exists && val != nil {
		newData["cancellation_submission_status"] = val
	}
	return newData
}

func (am *UploadData) DoUploadCustomerCreate(l *log.DXLog, data utils.JSON, option string) (id int64, err error) {
	fmt.Println("----- DoUploadCustomerCreate -----")
	data = am.validateCustomerData(l, data)

	fmt.Printf("%v", data)
	upperOption := strings.ToUpper(option)
	if upperOption == "CREATE_CONSTRUCTION" || upperOption == "CREATE_TASK_CONSTRUCTION" {
		data["is_create_construction"] = true
	}
	id, err = am.Customer.Insert(l, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (am *UploadData) DoUploadCustomerUpdate(l *log.DXLog, id int64, data utils.JSON) (err error) {
	_, err = am.Organization.UpdateOne(l, id, data)
	if err != nil {
		return err
	}
	return nil
}

func (am *UploadData) DoCustomerCreate(l *log.DXLog, customerData utils.JSON) (err error, customerId int64) {
	fmt.Println("========= DoCustomerCreate ==========")
	customerData = am.validateCustomerData(l, customerData)
	if customerData["row_status"] == "ERROR" {
		return errors.Errorf("cannot process: %v", customerData["row_message"]), customerId
	}

	db := database.Manager.Databases[task_management.ModuleTaskManagement.DatabaseNameId]

	fmt.Println("========= start db transaction ==========")
	txErr := db.Tx(l, database.LevelReadCommitted, func(tx *database.DXDatabaseTx) (err error) {
		insertData := utils.JSON{}
		fields := []string{
			"registration_number",
			"customer_number",
			"fullname",
			"email",
			"phonenumber",
			"identity_number",
			"identity_type",
			"npwp",
			"customer_segment_code",
			"customer_type_code",
			"jenis_anggaran",
			"rs_customer_sector_code",
			"sales_area_code",
			"latitude",
			"longitude",
			"address_name",
			"address_street",
			"address_rt",
			"address_rw",
			"address_kelurahan_location_code",
			"address_kecamatan_location_code",
			"address_kabupaten_location_code",
			"address_province_location_code",
			"address_postal_code",
			"register_at",
			"jenis_bangunan",
			"program_pelanggan",
			"payment_scheme_code",
			"kategory_wilayah",
			"cancellation_submission_status",
		}

		// Copy fields from customerData to insertData (skip null values)
		for _, field := range fields {
			if customerData[field] != nil {
				insertData[field] = customerData[field]
			}
		}

		// Set the status
		insertData["status"] = "ACTIVE"
		customerId, err = task_management.ModuleTaskManagement.Customer.TxInsert(tx, insertData)
		if err != nil {
			fmt.Println("ada error ==========")
			return err
		}
		// Check if is_create_construction exists and is a boolean
		isCreateConstruction, ok := customerData["is_create_construction"].(bool)
		if !ok {
			// Default to false if not present or not a boolean
			isCreateConstruction = false
		}

		if isCreateConstruction {
			//taskCode, _ := customerData["task_code"].(string) // Make task code optional
			taskCode, _ := customerData["registration_number"].(string) // Make task code optional
			if taskCode != "" {
				err = task_management.ModuleTaskManagement.ValidationTaskCodeShouldNotAlreadyExist(l, taskCode)
				if err != nil {
					return err
				}
			}
			_, err2 := task_management.ModuleTaskManagement.TaskTxCreateConstruction(tx, taskCode, customerId, "", "")
			if err2 != nil {
				return err2
			}
		}
		return nil
	})

	if txErr != nil {
		fmt.Println("====error===")
		fmt.Printf("%s", txErr.Error())
		// subTaskId is already set appropriately (either to a valid value or -1)
		return txErr, customerId
	} else {
		fmt.Println("====no error===")
	}
	return nil, customerId
}

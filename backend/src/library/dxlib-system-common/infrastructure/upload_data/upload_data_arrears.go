package upload_data

import (
	"fmt"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/pkg/errors"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/arrears_management"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/base"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	"strings"
	"time"
)

func (am *UploadData) validateArrearsData(l *log.DXLog, data utils.JSON) utils.JSON {
	data["row_message"] = ""
	if action, ok := data["action"].(string); !ok {
		setErrorMessage(data, "action column unknown")
	} else {
		trimmed := strings.TrimSpace(strings.Trim(action, "'"))
		data["action"] = trimmed
	}
	switch data["action"] {
	case "SPK SEGEL":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsStopGasFlow
	case "SPK CABUT":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsRemoveGasMeter
	case "SPK ALIR":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsOpenGasFlow
	case "SPK PASANG":
		data["sub_task_type_id"] = base.SubTaskTypeIdArrearsReinstallGasMeter
	default:
		setErrorMessage(data, fmt.Sprintf("action '%v' unknown", data["action"]))
	}
	if action, ok := data["id_pelanggan"].(string); !ok {
		setErrorMessage(data, "id_pelanggan column unknown")
	} else {
		trimmed := strings.TrimSpace(strings.Trim(action, "'"))
		data["id_pelanggan"] = trimmed
	}
	if _, row, err := task_management.ModuleTaskManagement.Customer.ShouldSelectOne(l, utils.JSON{
		"customer_number": data["id_pelanggan"],
	}, nil, nil); err == nil {
		data["customer_id"] = row["id"]
	} else {
		setErrorMessage(data, "id_pelanggan not found")
	}
	data["spk_no"] = data["nomor_surat"]

	if _, ok := data["tagihan_pemakaian_gas"].(int64); ok {
		data["amount_usage_bill"] = data["tagihan_pemakaian_gas"]
	}
	if _, ok := data["denda"].(int64); ok {
		data["amount_fine"] = data["denda"]
	}
	if _, ok := data["jaminan"].(int64); ok {
		data["amount_payment_guarantee"] = data["jaminan"]
	}
	if _, ok := data["biaya_pasang_kembali"].(int64); ok {
		data["amount_reinstallation_cost"] = data["biaya_pasang_kembali"]
	}
	if _, ok := data["biaya_alir_kembali"].(int64); ok {
		data["amount_reflow_cost"] = data["biaya_alir_kembali"]
	}
	if _, ok := data["jumlah_tagihan_rp"].(int64); ok {
		data["amount_bill_total"] = data["jumlah_tagihan_rp"]
	}

	layout := "2006-01-02"
	if dateStr, ok := data["periode_awal_tunggakan"].(string); ok {
		if date, err := time.Parse(layout, dateStr); err == nil {
			data["period_begin"] = date
		} else {
			setErrorMessage(data, "periode_awal_tunggakan format error")
		}
	}
	if dateStr, ok := data["periode_akhir_tunggakan"].(string); ok {
		if date, err := time.Parse(layout, dateStr); err == nil {
			data["period_end"] = date
		} else {
			setErrorMessage(data, "periode_akhir_tunggakan format error")
		}
	}
	if dateStr, ok := data["tgl_surat"].(string); ok {
		if date, err := time.Parse(layout, dateStr); err == nil {
			data["date_issued"] = date
		} else {
			setErrorMessage(data, "tgl_surat format error")
		}
	}
	return data
}

func (am *UploadData) DoUploadArrearsCreate(l *log.DXLog, data utils.JSON) (id int64, err error) {
	data = am.validateArrearsData(l, data)
	fmt.Printf("%v", data)
	id, err = am.Arrears.Insert(l, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (am *UploadData) DoArrearsCreate(l *log.DXLog, data utils.JSON) (err error, arrearsId int64) {
	fmt.Println("========= DoArrearsCreate ==========")
	data = am.validateArrearsData(l, data)
	if data["row_status"] == "ERROR" {
		return errors.Errorf("cannot process: %v", data["row_message"]), arrearsId
	}

	err, _ = arrears_management.ModuleArrearsManagement.DoTaskArrearsCreate(l, data)

	return err, arrearsId
}

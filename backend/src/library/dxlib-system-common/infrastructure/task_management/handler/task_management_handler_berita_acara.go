package handler

import (
	"bytes"
	"os"

	"github.com/pkg/errors"

	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/log"
	"github.com/donnyhardyanto/dxlib/object_storage"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func DownloadBeritaAcara(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskReportUid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetByUid(&aepr.Log, subTaskReportUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	filename := getData(subTaskReport, "berita_acara_link").(string)
	if filename == "" {
		return errors.New("berita acara belum digenerate")
	}

	err = task_management.ModuleTaskManagement.SubTaskReportBeritaAcara.DownloadSource(aepr, filename)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}

	return nil
}

func GenerateBeritaAcara(aepr *api.DXAPIEndPointRequest) (err error) {
	_, subTaskReportUid, err := aepr.GetParameterValueAsString("uid")
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.Log.Infof("Sub task report UID: %v", subTaskReportUid)
	_, subTaskReport, err := task_management.ModuleTaskManagement.SubTaskReport.ShouldGetByUid(&aepr.Log, subTaskReportUid)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if subTaskReport == nil {
		return errors.Errorf("SubTaskReport is not found:uid=%s", subTaskReportUid)
	}
	if len(subTaskReport) == 0 {
		return errors.Errorf("SubTaskReport is not found:uid=%s", subTaskReportUid)
	}

	if subTaskReport["sub_task_status"] != "WAITING_VERIFICATION" {
		return errors.New("SubTaskReport is not allowed to create berita acara")
	}

	subTaskTypeId := subTaskReport["sub_task_type_id"].(int64)
	subTaskReportId := subTaskReport["id"].(int64)

	customerId := subTaskReport["customer_id"].(int64)
	var gasMeter utils.JSON

	if subTaskTypeId == 4 {
		_, subTask, err := task_management.ModuleTaskManagement.SubTask.ShouldSelectOne(&aepr.Log, utils.JSON{
			"customer_id":      customerId,
			"sub_task_type_id": 3,
		}, nil, nil)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}

		if len(subTask) == 0 {
			return errors.New("subTask is empty")
		}

		lastSubTaskReportId := subTask["last_form_sub_task_report_id"].(int64)
		_, gasMeter, err = task_management.ModuleTaskManagement.SubTaskReport.ShouldGetById(&aepr.Log, lastSubTaskReportId)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}

		if len(gasMeter) == 0 {
			return errors.New("gasMeter is empty")
		}
	}

	objectStorage, exists := object_storage.Manager.ObjectStorages["berita-acara"]
	if !exists {
		return aepr.WriteResponseAndNewErrorf(http.StatusNotFound, "", "OBJECT_STORAGE_NAME_NOT_FOUND: berita-acara")
	}

	logo, err := os.ReadFile("./pertamina-logo.png")
	if err != nil {
		return errors.Wrap(err, "error reading logo file")
	}

	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithTopMargin(5).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(20).
		Build()
	doc := maroto.New(cfg)

	filesTTD, err := getTTDFile(&aepr.Log, subTaskReportId)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	if filesTTD == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "failed to get tanda tangan image")
	}

	var pdf []byte
	switch subTaskTypeId {
	case 1: // SK
		pdf, err = generateSK(doc, subTaskReport, logo, filesTTD, aepr.Log)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	case 4: // GAS IN
		pdf, err = generateGasIn(doc, subTaskReport, gasMeter, logo, filesTTD)
		if err != nil {
			return errors.Wrap(err, "error occured")
		}
	default:
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "UNSUPPORTED_SUB_TASK_TYPE_ID: %d", subTaskReport["sub_task_type_id"])
	}

	if pdf == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusBadRequest, "", "failed to generate PDF")
	}

	currentDate := time.Now()
	isoDate := currentDate.Format("2006-01-02")
	filename := fmt.Sprintf("%s_%d.pdf", isoDate, subTaskReportId)

	buf := bytes.NewBuffer(pdf)
	bufLen := int64(buf.Len())

	uploadInfo, err := objectStorage.UploadStream(buf, filename, filename, "application/octet-stream", false, bufLen)
	if err != nil {
		return errors.Wrap(err, "error occured")
	}
	aepr.Log.Infof("Upload info result: %v", uploadInfo)

	_, err = task_management.ModuleTaskManagement.SubTaskReport.Update(utils.JSON{
		"berita_acara_link": filename,
	}, utils.JSON{
		"id": subTaskReportId,
	})

	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"link": filename,
	})

	return nil
}

func generateGasIn(doc core.Maroto, gasIn utils.JSON, gasMeter utils.JSON, logo []byte, filesTTD [][]byte) ([]byte, error) {
	gasInReport, err := utils.GetJSONFromKV(gasIn, "report")
	if err != nil {
		return nil, err
	}

	gasMeterReport, err := utils.GetJSONFromKV(gasMeter, "report")
	if err != nil {
		return nil, err
	}

	// Add Pertamina logo
	doc.AddRow(20,
		image.NewFromBytesCol(12, logo, extension.Png,
			props.Rect{
				Left: 145,
			},
		),
	)

	currentTime := time.Now()
	day := currentTime.Day()
	month := getBulan(int(currentTime.Month()))
	year := currentTime.Year()
	weekday := currentTime.Weekday().String()
	weekday = getHari(weekday)

	bagt := fmt.Sprintf("%v%v%v-%v", currentTime.Day(), currentTime.Month(), currentTime.Year(), gasIn["id"])

	// Title
	titleProps := props.Text{
		Size:   16,
		Style:  fontstyle.Bold,
		Align:  align.Center,
		Family: fontfamily.Arial,
		Color:  getBlueColor(),
	}
	doc.AddRow(15, text.NewCol(12, "BERITA ACARA GAS IN", titleProps))

	doc.AddRow(7, text.NewCol(12, fmt.Sprintf("BAGt. %s", bagt), props.Text{
		Size:  12,
		Style: fontstyle.BoldItalic,
		Align: align.Left,
		Color: getRedColor(),
	}))

	r := doc.AddRow(12, text.NewCol(12, fmt.Sprintf(`Pada hari ini %s, %d %s %d, telah dilakukan penyaluran Gas pertama kali ("Tanggal dimulai") kepada`, weekday, day, month, year), props.Text{
		Top:   3,
		Left:  3,
		Size:  10,
		Color: &props.WhiteColor,
	}))
	r.WithStyle(&props.Cell{
		BackgroundColor: getBlueColor(),
	})

	doc.AddRow(3)
	r = doc.AddRow(12,
		text.NewCol(12, "Informasi Pelanggan", props.Text{
			Top:   3,
			Left:  3,
			Size:  12,
			Color: &props.WhiteColor,
			Style: fontstyle.Bold,
		}))
	r.WithStyle(&props.Cell{
		BackgroundColor: getBlueColor2(),
	})

	customerName := getData(gasIn, "customer_fullname")
	customerNo := getData(gasIn, "customer_number")
	doc.AddRow(8,
		text.NewCol(2, "No. Pelanggan", props.Text{
			Top:  5,
			Left: 3,
			Size: 10,
		}),
		text.NewCol(3, fmt.Sprintf("%v", customerNo), props.Text{
			Top:  5,
			Size: 10,
		}),
		text.NewCol(2, "Nama Lengkap", props.Text{
			Top:  5,
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", customerName), props.Text{
			Top:  5,
			Size: 10,
		}),
	)
	doc.AddRow(5,
		col.New(2),
		line.NewCol(3, props.Line{
			Thickness:     0.2,
			OffsetPercent: 50,
			SizePercent:   100,
		}),
		col.New(2),
		line.NewCol(5, props.Line{
			Thickness:     0.2,
			OffsetPercent: 50,
			SizePercent:   100,
		}),
	)

	customerAddress := getData(gasIn, "customer_address_street")
	doc.AddRow(5,
		text.NewCol(2, "Alamat", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(9, fmt.Sprintf("%v", customerAddress), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(1,
		col.New(2),
		line.NewCol(10, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	customerRT := getData(gasIn, "customer_address_rt")
	customerRW := getData(gasIn, "customer_address_rw")
	doc.AddRow(5,
		col.New(2),
		text.NewCol(8, "", props.Text{
			Size: 10,
		}),
		text.NewCol(1, fmt.Sprintf("RT %v", customerRT), props.Text{
			Size: 10,
		}),
		text.NewCol(1, fmt.Sprintf("RW %v", customerRW), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(1,
		col.New(2),
		line.NewCol(7, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(1),
		line.NewCol(1, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		line.NewCol(1, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	customerSubdistrict := getData(gasIn, "customer_address_kelurahan_location_name")
	customerDistrict := getData(gasIn, "customer_address_kecamatan_location_name")
	doc.AddRow(5,
		text.NewCol(2, "Kelurahan", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", customerSubdistrict), props.Text{
			Size: 10,
		}),
		text.NewCol(2, "Kecamatan", props.Text{
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", customerDistrict), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(1,
		col.New(2),
		line.NewCol(4, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(2),
		line.NewCol(4, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	customerCity := getData(gasIn, "customer_address_kabupaten_location_name")
	customerProvince := getData(gasIn, "customer_address_province_location_name")
	doc.AddRow(5,
		text.NewCol(2, "Kota Kabupaten", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(5, fmt.Sprintf("%v", customerCity), props.Text{
			Size: 10,
		}),
		text.NewCol(1, "Provinsi", props.Text{
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", customerProvince), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(1,
		col.New(2),
		line.NewCol(5, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(1),
		line.NewCol(4, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	customerPostalCode := getData(gasIn, "customer_address_postal_code")
	customerLatitude := getData(gasIn, "customer_latitude")
	customerLongitude := getData(gasIn, "customer_longitude")
	doc.AddRow(5,
		text.NewCol(3, "Latitude/Longitude", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(5, fmt.Sprintf("%v/%v", customerLatitude, customerLongitude), props.Text{
			Size: 10,
		}),
		text.NewCol(1, "Kode Pos", props.Text{
			Size: 10,
		}),
		text.NewCol(3, fmt.Sprintf("%v", customerPostalCode), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(5,
		col.New(3),
		line.NewCol(5, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(1),
		line.NewCol(3, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	).WithStyle(&props.Cell{
		BorderColor:     getBlueColor2(),
		BorderThickness: 0.5,
		BorderType:      border.Bottom,
	})

	doc.AddRow(10,
		text.NewCol(12, "Bersama ini telah dilakukan hal-hal sbb:", props.Text{
			Top:   5,
			Left:  5,
			Size:  12,
			Color: getBlueColor(),
			Style: fontstyle.Bold,
		}),
	)

	type Item struct {
		description string
		status      string
		keterangan  string
	}

	items := []Item{
		{"Berita Acara hasil PKS/Opname yang konstruksi", "OK", "Tersedia"},
		{"Terdapat prosedur Gas In ke konsumen dan telah penyerahan", "OK", "Terlaksana"},
		{"Tersedia perlengkapan K3PL yang memadai", "OK", "Tersedia"},
		{"Sosialisasi pengoperasian & pemeliharaan kepada pelanggan", "OK", "Terlaksana"},
		{"Meter terkalibrasi", "OK", "Tersedia"},
		{"Tes kebocoran gas", "OK", "Tidak ada kebocoran"},
	}

	doc.AddRow(2)
	for _, item := range items {
		doc.AddRow(7,
			text.NewCol(1, "•", props.Text{
				Size:  18,
				Left:  5,
				Style: fontstyle.Bold,
				Color: getBlueColor(),
			}),
			text.NewCol(7, item.description, props.Text{
				Top:  2,
				Size: 10,
			}),
			text.NewCol(1, item.status, props.Text{
				Top:  2,
				Size: 10,
			}),
			text.NewCol(3, item.keterangan, props.Text{
				Top:  2,
				Size: 10,
			}),
		)
		doc.AddRow(1,
			col.New(8),
			line.NewCol(4, props.Line{
				Thickness:     0.2,
				OffsetPercent: 0,
				SizePercent:   100,
			}),
		)
	}

	doc.AddRow(14,
		text.NewCol(12, "Data Meter Gas Terpasang:", props.Text{
			Top:   5,
			Left:  5,
			Size:  12,
			Color: getBlueColor(),
			Style: fontstyle.Bold,
		}),
	)

	meterBrand := getData(gasInReport, "meter_brand")
	snMeter := getData(gasInReport, "sn_meter")
	doc.AddRow(5,
		text.NewCol(2, "Jenis Meter:", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", meterBrand), props.Text{
			Size: 10,
		}),
		text.NewCol(2, "G Size/SN Meter:", props.Text{
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", snMeter), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(1,
		col.New(2),
		line.NewCol(2, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(4),
		line.NewCol(2, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	standMeterStart := getData(gasInReport, "stand_meter_start_number")
	qMin := getData(gasMeterReport, "qmin")
	qMax := getData(gasMeterReport, "qmax")
	doc.AddRow(5,
		text.NewCol(2, "Qmin/Qmax:", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v/%v m3/jam", qMin, qMax), props.Text{
			Size: 10,
		}),
		text.NewCol(2, "Stand Meter Awal:", props.Text{
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v", standMeterStart), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(1,
		col.New(2),
		line.NewCol(2, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(4),
		line.NewCol(2, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	temperature := getData(gasInReport, "temperature_start")
	pressure := getData(gasInReport, "pressure_start")
	startCalibrationMonth := int(getData(gasMeterReport, "start_calibration_month").(float64))
	startCalibrationYear := getData(gasMeterReport, "start_calibration_year")

	doc.AddRow(5,
		text.NewCol(2, "Awal Kalibrasi:", props.Text{
			Left: 3,
			Size: 10,
		}),
		text.NewCol(4, fmt.Sprintf("%v %v", getBulan(startCalibrationMonth), startCalibrationYear), props.Text{
			Size: 10,
		}),
		text.NewCol(1, "Tekanan:", props.Text{
			Size: 10,
		}),
		text.NewCol(2, fmt.Sprintf("%v mbar", pressure), props.Text{
			Size: 10,
		}),
		text.NewCol(1, "Suhu:", props.Text{
			Size: 10,
		}),
		text.NewCol(2, fmt.Sprintf("%v°C (Jika Ada)", temperature), props.Text{
			Size: 10,
		}),
	)
	doc.AddRow(5,
		col.New(2),
		line.NewCol(2, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(3),
		line.NewCol(1, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
		col.New(2),
		line.NewCol(2, props.Line{
			Thickness:     0.2,
			OffsetPercent: 0,
			SizePercent:   100,
		}),
	)

	doc.AddRow(15,
		text.NewCol(12, "Dengan dilakukannya Gas in ini, maka Pelanggan menyetujui untuk bertanggung jawab atas pengoperasian dan perawatan pipa instalasi beserta seluruh konsekuensinya. Demikian berita acara ini dibuat untuk dipergunakan sebagaiamana mestinya.", props.Text{
			Size:  10,
			Left:  1,
			Align: align.Justify,
		}),
	)

	doc.AddRow(5,
		text.NewCol(3, "Pelanggan", props.Text{
			Size: 10,
			Left: 15,
		}),
		text.NewCol(3, "Petugas", props.Text{
			Size: 10,
			Left: 15,
		}),
	)

	ttdPelanggan := filesTTD[0]
	ttdPetugas := filesTTD[1]
	doc.AddRow(15,
		image.NewFromBytesCol(3, ttdPelanggan, extension.Png, props.Rect{
			Top:     3,
			Left:    15,
			Percent: 75,
		}),
		image.NewFromBytesCol(3, ttdPetugas, extension.Png, props.Rect{
			Top:     3,
			Left:    15,
			Percent: 75,
		}),
	)

	officerName := getData(gasIn, "user_fullname")
	doc.AddRow(7,
		signature.NewCol(3, fmt.Sprintf("%v", customerName), props.Signature{
			FontSize:  10,
			FontStyle: fontstyle.Bold,
		}),
		signature.NewCol(3, fmt.Sprintf("%v", officerName), props.Signature{
			FontSize:  10,
			FontStyle: fontstyle.Bold,
		}),
	)

	doc.AddRow(10,
		text.NewCol(12, "Terima kasih atas kepercayaan Saudara kepada kami.", props.Text{
			Top:   5,
			Size:  10,
			Left:  1,
			Align: align.Justify,
		}),
	)
	doc.AddRow(2,
		text.NewCol(12, "PT Perusahaan Gas Negara Tbk.", props.Text{
			Size:  10,
			Left:  1,
			Align: align.Justify,
		}),
	)

	// Save the PDF
	document, err := doc.Generate()
	if err != nil {
		return nil, err
	}

	err = document.Save("GAS IN.pdf")
	if err != nil {
		return nil, err
	}

	pdf := document.GetBytes()

	return pdf, nil
}

func generateSK(doc core.Maroto, SK utils.JSON, logo []byte, filesTTD [][]byte, logger log.DXLog) ([]byte, error) {
	report, err := utils.GetJSONFromKV(SK, "report")
	if err != nil {
		return nil, err
	}

	doc.AddRow(18,
		image.NewFromBytesCol(12, logo, extension.Png,
			props.Rect{
				Left: 145,
			},
		),
	)

	doc.AddRow(7,
		text.NewCol(12, "BERITA ACARA", props.Text{
			Style: fontstyle.Bold,
			Size:  14,
			Align: align.Center,
		}),
	)
	doc.AddRow(15,
		text.NewCol(12, "SAMBUNGAN PIPA DAN PERALATAN GAS", props.Text{
			Style: fontstyle.Bold,
			Size:  14,
			Align: align.Center,
		}),
	)

	currentTime := time.Now()
	day := currentTime.Day()
	month := getBulan(int(currentTime.Month()))
	year := currentTime.Year()
	weekday := currentTime.Weekday().String()
	weekday = getHari(weekday)
	doc.AddRow(6,
		text.NewCol(12, fmt.Sprintf("Pada hari ini %s Tanggal %d Bulan %s Tahun %d, telah disepakati di lokasi:", weekday, day, month, year)),
	)

	customerName := getData(SK, "customer_fullname")
	doc.AddRow(6,
		text.NewCol(3, "Nama"),
		text.NewCol(9, fmt.Sprintf(": %v (Sesuai NPWP Jika Ada)", customerName)),
	)

	customerNIK := getData(SK, "customer_identity_number")
	doc.AddRow(6,
		text.NewCol(3, "NIK (No KTP)"),
		text.NewCol(9, fmt.Sprintf(": %v (Sesuai NPWP Jika Ada)", customerNIK)),
	)

	customerNPWP := getData(SK, "customer_npwp")
	doc.AddRow(6,
		text.NewCol(3, "NPWP (Jika ada)"),
		text.NewCol(9, fmt.Sprintf(": %v", customerNPWP)),
	)

	customerAddress := getData(SK, "customer_address")
	doc.AddAutoRow(
		text.NewCol(3, "Alamat (Sesuai KTP)"),
		text.NewCol(9, fmt.Sprintf(": %v", customerAddress)),
	)

	doc.AddRow(1)

	noHP := getData(SK, "customer_phonenumber")
	email := getData(SK, "customer_email")
	doc.AddRow(6,
		text.NewCol(3, "No. HP/WA"),
		text.NewCol(2, fmt.Sprintf(": %v", noHP)),
		text.NewCol(1, "Email"),
		text.NewCol(8, fmt.Sprintf(": %v", email)),
	)

	idPelanggan := getData(SK, "customer_registration_number")
	doc.AddRow(10,
		text.NewCol(3, "ID Pelanggan*)"),
		text.NewCol(11, fmt.Sprintf(": %v (Sesuai yang didaftarkan ke PGN)", idPelanggan)),
	)

	doc.AddRow(6,
		text.NewCol(12, "Kebutuhan Pipa Instalasi:"),
	)
	doc.AddRow(10,
		text.NewCol(1, "No", props.Text{
			Style: fontstyle.Bold,
			Left:  2,
			Align: align.Center,
		}),
		text.NewCol(3, "Pipa Instalasi", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Panjang (meter)", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Kelebihan Pipa (meter)", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Biaya Kelebihan Pipa (Exc. PPN)", props.Text{Style: fontstyle.Bold, Align: align.Center}),
	).WithStyle(
		&props.Cell{
			BorderThickness: 0.5,
			BorderType:      border.Full,
		},
	)

	pipeLength := getData(report, "pipe_length")
	extraPipeLength := getData(report, "calculated_extra_pipe_length")
	doc.AddRow(5,
		text.NewCol(1, fmt.Sprintf("%d", 1), props.Text{Left: 2}),
		text.NewCol(3, "Pipa Instalasi"),
		text.NewCol(2, fmt.Sprintf("%v", pipeLength)),
		text.NewCol(3, fmt.Sprintf("%v", extraPipeLength)),
		text.NewCol(3, "-"),
	).WithStyle(
		&props.Cell{
			BorderThickness: 0.5,
			BorderType:      border.Full,
			BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
		},
	)

	doc.AddRow(11,
		text.NewCol(12, "Kebutuhan Konversi Peralatan Gas:",
			props.Text{
				Top: 5,
			}),
	)

	gasAppliances := getData(report, "gas_appliances").([]interface{})
	listGasAppliances := []GasAppliance{}
	logger.Infof("gas appliances data: %v", gasAppliances)
	for _, gas := range gasAppliances {
		tmp := gas.(utils.JSON)
		gasAppliance := GasAppliance{
			name:   tmp["name"],
			jumlah: tmp["quantity"],
		}
		listGasAppliances = append(listGasAppliances, gasAppliance)
	}
	if len(listGasAppliances) == 0 {
		gasAppliance := GasAppliance{
			name:   "",
			jumlah: "",
		}
		listGasAppliances = append(listGasAppliances, gasAppliance)
	}
	rows, err := list.Build[GasAppliance](listGasAppliances)
	if err != nil {
		return nil, err
	}
	doc.AddRows(rows...)

	doc.AddRow(12,
		text.NewCol(4, "Total Tagihan: -", props.Text{
			Top:   5,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, ", termasuk PPN & materai (jika ada)", props.Text{
			Top: 5,
		}),
	)

	doc.AddRow(18,
		text.NewCol(12, "Calon Pelanggan bersedia melakukan pembayaran Biaya Berlangganan atas Kelebihan Sambungan Pipa dan Peralatan Gas sebelum pengaliran Gas sesuai dengan ketentuan PGN. Selanjutnya, pemeliharaan peralatan gas dan pipa instalasi menjadi tanggung jawab Pelanggan.", props.Text{
			Align:           align.Justify,
			VerticalPadding: 2,
		}),
	)

	doc.AddRow(10, text.NewCol(12, "Demikian Berita Acara ini dibuat dengan sebenarnya untuk dipergunakan sebagaimana mestinya.", props.Text{
		Align: align.Justify,
	}))

	doc.AddRow(5,
		text.NewCol(3, "Calon Pelanggan", props.Text{
			Size: 10,
			Left: 10,
		}),
		text.NewCol(3, "Petugas", props.Text{
			Size: 10,
			Left: 15,
		}),
	)

	ttdPelanggan := filesTTD[0]
	ttdPetugas := filesTTD[1]
	doc.AddRow(15,
		image.NewFromBytesCol(3, ttdPelanggan, extension.Png, props.Rect{
			Top:     3,
			Left:    15,
			Percent: 75,
		}),
		image.NewFromBytesCol(3, ttdPetugas, extension.Png, props.Rect{
			Top:     3,
			Left:    15,
			Percent: 75,
		}),
	)

	operatorName := getData(SK, "user_fullname")
	doc.AddRow(7,
		signature.NewCol(3, fmt.Sprintf("%v", customerName), props.Signature{
			FontSize:  10,
			FontStyle: fontstyle.Bold,
		}),
		signature.NewCol(3, fmt.Sprintf("%v", operatorName), props.Signature{
			FontSize:  10,
			FontStyle: fontstyle.Bold,
		}),
	)

	doc.AddRow(5)
	doc.AddRow(6, text.NewCol(12, "Keterangan (Khusus Pelanggan Rumah Tangga): "))
	doc.AddRow(6, text.NewCol(12, "1. Panjang maksimum sambungan kompor yang ditanggung PGN: RTI = 10 meter; GPR = 15 meter."))
	doc.AddRow(6, text.NewCol(12, "2. Konversi peralatan Gas yang ditanggung oleh PGN adalah 1 unit Kompor maks. 2 tungku."))
	doc.AddRow(6, text.NewCol(12, "3. Kelebihan atas sambungan kompor dan konversi peralatan gas dikenakan biaya sebesar:"))
	doc.AddRow(6, text.NewCol(12, "a. Pipa Instalasi: Rp. 75.000/meter (Exc. PPN)", props.Text{
		Left: 3,
	}))
	doc.AddRow(6, text.NewCol(12, "b. Konversi kompor: Rp. 25.000/Kompor (Exc. PPN)", props.Text{
		Left: 3,
	}))
	doc.AddRow(6, text.NewCol(12, "c. Konversi Water Heater Gas: RT1=Rp. 125.000/Unit (Exc. PPN), GPR=Rp. 75.000/Unit (Exc. PPN)", props.Text{
		Left: 3,
	}))
	doc.AddRow(6, text.NewCol(12, "4. Pelanggan/Calon Pelanggan agar mendokumentasikan (mengambil foto) berita acara ini."))
	doc.AddRow(6, text.NewCol(12, "5. Pembayaran dapat dilakukan mulai tanggal 6 sd akhir bulan melalui channel pembayaran tersedia (ATM, Mandiri, BRI, BNI, BSI, BCA, Indomaret, Alfamart, Gopay, Tokopedia, Shopee, LinjAja, Kantor Pos, PPOB).", props.Text{
		VerticalPadding: 2,
		Align:           align.Justify,
	}))

	document, err := doc.Generate()
	if err != nil {
		return nil, err
	}

	err = document.Save("SK.pdf")
	if err != nil {
		return nil, err
	}

	pdf := document.GetBytes()

	return pdf, nil
}

type GasAppliance struct {
	name           interface{}
	jumlah         interface{}
	kelebihan      interface{}
	biayaKelebihan interface{}
}

func (o GasAppliance) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(1, "No", props.Text{
			Style: fontstyle.Bold,
			Left:  2,
			Align: align.Center,
		}),
		text.NewCol(3, "Jenis Peralatan Gas", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Jumlah Spuyer (unit)", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Kelebihan Spuyer (unit)", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Biaya Kelebihan Spuyer (Exc. PPN)", props.Text{Style: fontstyle.Bold, Align: align.Center}),
	).WithStyle(
		&props.Cell{
			BorderThickness: 0.5,
			BorderType:      border.Full,
		},
	)
}

func (o GasAppliance) GetContent(i int) core.Row {
	r := row.New(5).Add(
		text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Left: 2}),
		text.NewCol(3, fmt.Sprintf("%v", o.name)),
		text.NewCol(2, fmt.Sprintf("%v", o.jumlah)),
		text.NewCol(3, "-"),
		text.NewCol(3, "-"),
	).WithStyle(
		&props.Cell{
			BorderThickness: 0.5,
			BorderType:      border.Full,
		},
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
			BorderThickness: 0.5,
			BorderType:      border.Full,
		})
	}

	return r
}

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   7,
		Green: 148,
		Blue:  210,
	}
}

func getRedColor() *props.Color {
	return &props.Color{
		Red:   255,
		Green: 0,
		Blue:  0,
	}
}

func getBlueColor2() *props.Color {
	return &props.Color{
		Red:   1,
		Green: 88,
		Blue:  158,
	}
}

func getHari(day string) string {
	dayInIndonesian := map[string]string{
		"Monday":    "Senin",
		"Tuesday":   "Selasa",
		"Wednesday": "Rabu",
		"Thursday":  "Kamis",
		"Friday":    "Jumat",
		"Saturday":  "Sabtu",
		"Sunday":    "Minggu",
	}

	return dayInIndonesian[day]
}

func getBulan(n int) string {
	indonesianMonths := map[int]string{
		1: "Januari", 2: "Februari", 3: "Maret", 4: "April", 5: "Mei", 6: "Juni",
		7: "Juli", 8: "Agustus", 9: "September", 10: "Oktober", 11: "November", 12: "Desember",
	}

	return indonesianMonths[n]
}

func getData(data map[string]interface{}, key string) interface{} {
	if value, ok := data[key]; ok {
		// If key exists, print the value
		return value
	}

	return ""
}

func getTTDFile(log *log.DXLog, subTaskReportId int64) ([][]byte, error) {
	idPelanggan := "TTD_PELANGGAN"
	idPetugas := "TTD_PETUGAS"

	// Function to download the TTD file
	downloadFile := func(nameId string) ([]byte, error) {
		// Perform database query to get the file details
		_, ttdData, err := task_management.ModuleTaskManagement.SubTaskReportFile.ShouldSelectOne(log, utils.JSON{
			"sub_task_report_id":                subTaskReportId,
			"sub_task_report_file_group_nameid": nameId,
		}, nil, nil)
		if err != nil {
			return nil, errors.Errorf("failed to get file info for ID %v: %v", nameId, err)
		}

		// Ensure type assertion success
		taskId, ok := ttdData["task_id"].(int64)
		if !ok {
			return nil, errors.New("task_id type assertion failed")
		}
		subTaskId, ok := ttdData["sub_task_id"].(int64)
		if !ok {
			return nil, errors.New("sub_task_id type assertion failed")
		}
		groupId, ok := ttdData["sub_task_report_file_group_id"].(int64)
		if !ok {
			return nil, errors.New("sub_task_report_file_group_id type assertion failed")
		}
		fileId, ok := ttdData["sub_task_file_id"].(int64)
		if !ok {
			return nil, errors.New("sub_task_file_id type assertion failed")
		}

		// Format file path
		filePath := fmt.Sprintf("%d/%d/%d/%d.png", taskId, subTaskId, groupId, fileId)

		// Fetch the file from object storage
		objectStorage, exists := object_storage.Manager.ObjectStorages["sub-task-report-picture-small"]
		if !exists {
			return nil, errors.New("OBJECT_STORAGE_NAME_NOT_FOUND: sub-task-report-picture-small")
		}

		object, err := objectStorage.DownloadStream(filePath)
		if err != nil {
			return nil, errors.Errorf("failed to download file for path %s: %v", filePath, err)
		}
		defer func() {
			_ = object.Close()
		}()

		file, err := io.ReadAll(object)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	// Download files for Pelanggan and Petugas
	filePelanggan, err := downloadFile(idPelanggan)
	if err != nil {
		return nil, errors.Errorf("error downloading file for Pelanggan: %v", err)
	}

	filePetugas, err := downloadFile(idPetugas)
	if err != nil {
		return nil, errors.Errorf("error downloading file for Petugas: %v", err)
	}

	// Return both files as a 2D byte slice
	return [][]byte{filePelanggan, filePetugas}, nil
}

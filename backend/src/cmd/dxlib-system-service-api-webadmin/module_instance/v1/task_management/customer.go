package task_management

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/table"
	"github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/donnyhardyanto/dxlib_module/module/self"
	"github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management"
	task_management_handler "github.com/donnyhardyanto/dxlib-system/common/infrastructure/task_management/handler"
)

func defineAPICustomer(anAPI *api.DXAPI) {

	anAPI.NewEndPoint("Customer.Upload.CMS",
		"Upload file csv or Excel to creates some new Customer in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Customer record with assigned unique identifier.",
		"/v1/customer/create_bulk", "POST", api.EndPointTypeHTTPUploadStream, http.ContentTypeApplicationOctetStream, []api.DXAPIEndPointParameter{},
		task_management_handler.CustomerCreateBulk, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.List.CMS",
		"Retrieves a paginated list of Customer  with filtering and sorting capabilities. "+
			"Returns a structured list of Customer  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.CustomerList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.List.Per.Construction.Area.CMS",
		"Retrieves a paginated list of Customer  with filtering and sorting capabilities. "+
			"Returns a structured list of Customer  with their basic information"+
			"filtering conditions and flexible result ordering.",
		"/v1/customer/list/per_construction_area", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.CustomerPerConstructionAreaList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.LIST.PER.CONSTRUCTION.AREA"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.Download.CMS",
		"Download CSV or Excel file based on given filters",
		"/v1/customer/list/download", "POST", api.EndPointTypeHTTPDownloadStream, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "format", Type: "protected-string", Description: "", IsMustExist: true},
		}, task_management_handler.CustomerListDownload, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.DOWNLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.Create.CMS",
		"Creates a new Customer in the system with validated information. "+
			"Handles the registration process with proper data validation and standardization. "+
			"Returns the created Customer record with assigned unique identifier.",
		"/v1/customer/create", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "registration_number", Type: "string", Description: "Customer Registration Number", IsMustExist: true},
			{NameId: "customer_number", Type: "string", Description: "Customer Number", IsMustExist: true},
			{NameId: "fullname", Type: "string", Description: "Customer Fullname", IsMustExist: true},
			{NameId: "status", Type: "string", Description: "Customer Status", IsMustExist: true},
			{NameId: "special_flag_sk_primer", Type: "bool", Description: "Special flag SK Primer", IsMustExist: false},
			{NameId: "email", Type: "email", Description: "Customer Email", IsMustExist: false},
			{NameId: "phonenumber", Type: "phonenumber", Description: "Customer Phonenumber", IsNullable: true, IsMustExist: false},
			{NameId: "korespondensi_media", Type: "string", Description: "Customer Korespondensi Media", IsMustExist: false},
			{NameId: "identity_type", Type: "string", Description: "Customer identity_type", IsMustExist: false},
			{NameId: "identity_number", Type: "string", Description: "Customer identity_number", IsMustExist: false},
			{NameId: "npwp", Type: "npwp", Description: "Customer npwp", IsNullable: true, IsMustExist: false},
			{NameId: "customer_segment_code", Type: "string", Description: "Customer customer_segment_code", IsMustExist: true},
			{NameId: "customer_type_code", Type: "string", Description: "Customer customer_type_code", IsMustExist: true},
			{NameId: "jenis_anggaran", Type: "string", Description: "Customer jenis anggaran", IsMustExist: true},
			{NameId: "rs_customer_sector_code", Type: "nullable-string", Description: "Customer rs_customer_sector_code", IsNullable: true, IsMustExist: false},
			{NameId: "sales_area_code", Type: "string", Description: "Customer sales_area_code", IsMustExist: true},
			{NameId: "latitude", Type: "float64", Description: "Customer latitude", IsMustExist: true},
			{NameId: "longitude", Type: "float64", Description: "Customer longitude", IsMustExist: true},
			{NameId: "address_name", Type: "string", Description: "Customer address_name", IsMustExist: false},
			{NameId: "address_street", Type: "string", Description: "Customer address_street", IsMustExist: true},
			{NameId: "address_rt", Type: "string", Description: "Customer address_rt", IsMustExist: true},
			{NameId: "address_rw", Type: "string", Description: "Customer address_rw", IsMustExist: true},
			{NameId: "address_kelurahan_location_code", Type: "string", Description: "Customer address_kelurahan_location_code", IsMustExist: true},
			{NameId: "address_kecamatan_location_code", Type: "string", Description: "Customer address_kecamatan_location_code", IsMustExist: true},
			{NameId: "address_kabupaten_location_code", Type: "string", Description: "Customer address_kabupaten_location_code", IsMustExist: true},
			{NameId: "address_province_location_code", Type: "string", Description: "Customer address_province_location_code", IsMustExist: true},
			{NameId: "address_postal_code", Type: "string", Description: "Customer address_postal_code", IsMustExist: true},
			{NameId: "register_at", Type: "iso8601", Description: "Customer registration_timestamp", IsMustExist: false},
			{NameId: "jenis_bangunan", Type: "string", Description: "Customer jenis_bangunan", IsMustExist: true},
			{NameId: "program_pelanggan", Type: "string", Description: "Customer program_pelanggan", IsMustExist: true},
			{NameId: "payment_scheme_code", Type: "nullable-string", Description: "Customer payment_scheme_code", IsNullable: true, IsMustExist: false},
			{NameId: "kategory_wilayah", Type: "string", Description: "Customer kategory_wilayah", IsMustExist: true},
		}, task_management_handler.CustomerCreate, nil, table.Manager.StandardOperationResponsePossibility["create"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.CREATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.Read.CMS",
		"Retrieves detailed information for a specific Customer by ID. "+
			"Returns comprehensive Customer data including specifications. "+
			"Essential for Announcement specification views and data verification.",
		"/v1/customer/read", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management.ModuleTaskManagement.Customer.RequestRead, nil, table.Manager.StandardOperationResponsePossibility["read"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.READ"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.Edit.CMS",
		"Updates Customer information with comprehensive data validation. "+
			"Allows modification of Customer information while maintaining data integrity. "+
			"Supports partial updates with selective field modifications.",
		"/v1/customer/edit", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "new", Type: "json", Description: "", IsMustExist: true, Children: []api.DXAPIEndPointParameter{
				{NameId: "registration_number", Type: "string", Description: "Customer Registration Number", IsMustExist: false},
				{NameId: "customer_number", Type: "string", Description: "Customer Number", IsMustExist: false},
				{NameId: "fullname", Type: "string", Description: "Customer Fullname", IsMustExist: false},
				{NameId: "status", Type: "string", Description: "Customer Status", IsMustExist: false},
				{NameId: "special_flag_sk_primer", Type: "bool", Description: "Special flag SK Primer", IsMustExist: false},
				{NameId: "email", Type: "string", Description: "Customer Email", IsMustExist: false},
				{NameId: "phonenumber", Type: "phonenumber", Description: "Customer Phonenumber", IsMustExist: false},
				{NameId: "korespondensi_media", Type: "string", Description: "Customer Korespondensi Media", IsMustExist: false},
				{NameId: "identity_type", Type: "string", Description: "Customer identity_type", IsMustExist: false},
				{NameId: "identity_number", Type: "string", Description: "Customer identity_number", IsMustExist: false},
				{NameId: "npwp", Type: "npwp", Description: "Customer npwp", IsMustExist: false},
				{NameId: "customer_segment_code", Type: "string", Description: "Customer customer_segment_code", IsMustExist: false},
				{NameId: "customer_type_code", Type: "string", Description: "Customer customer_type_code", IsMustExist: false},
				{NameId: "jenis_anggaran", Type: "string", Description: "Customer jenis_anggaran", IsMustExist: false},
				{NameId: "rs_customer_sector_code", Type: "string", Description: "Customer rs_customer_sector_code", IsMustExist: false},
				{NameId: "sales_area_code", Type: "string", Description: "Customer sales_area_code", IsMustExist: false},
				{NameId: "latitude", Type: "float64", Description: "Customer latitude", IsMustExist: false},
				{NameId: "longitude", Type: "float64", Description: "Customer longitude", IsMustExist: false},
				{NameId: "address_name", Type: "string", Description: "Customer address_name", IsMustExist: false},
				{NameId: "address_street", Type: "string", Description: "Customer address_street", IsMustExist: false},
				{NameId: "address_rt", Type: "string", Description: "Customer address_rt", IsMustExist: false},
				{NameId: "address_rw", Type: "string", Description: "Customer address_rw", IsMustExist: false},
				{NameId: "address_kelurahan_location_code", Type: "string", Description: "Customer address_kelurahan_location_code", IsMustExist: false},
				{NameId: "address_kecamatan_location_code", Type: "string", Description: "Customer address_kecamatan_location_code", IsMustExist: false},
				{NameId: "address_kabupaten_location_code", Type: "string", Description: "Customer address_kabupaten_location_code", IsMustExist: false},
				{NameId: "address_province_location_code", Type: "string", Description: "Customer address_province_location_code", IsMustExist: false},
				{NameId: "address_postal_code", Type: "string", Description: "Customer address_postal_code", IsMustExist: false},
				{NameId: "register_at", Type: "iso8601", Description: "Customer registration_timestamp", IsMustExist: false},
				{NameId: "jenis_bangunan", Type: "string", Description: "Customer jenis_bangunan", IsMustExist: false},
				{NameId: "program_pelanggan", Type: "string", Description: "Customer program_pelanggan", IsMustExist: false},
				{NameId: "payment_scheme_code", Type: "string", Description: "Customer payment_scheme_code", IsMustExist: false},
				{NameId: "kategory_wilayah", Type: "string", Description: "Customer kategory_wilayah", IsMustExist: false},
			}},
		}, task_management_handler.CustomerEdit, nil, table.Manager.StandardOperationResponsePossibility["edit"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPDATE"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.Delete.CMS",
		"Permanently removes a Customer record from the system. "+
			"Performs necessary validation checks before deletion. "+
			"Returns confirmation of successful deletion operation.",
		"/v1/customer/delete", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "id", Type: "int64", Description: "", IsMustExist: true},
		}, task_management_handler.CustomerSoftDelete, nil, table.Manager.StandardOperationResponsePossibility["delete"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.DELETE"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.Meter.List.CMS",
		"",
		"/v1/customer_meter/list", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filter_where", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_order_by", Type: "string", Description: "", IsMustExist: true},
			{NameId: "filter_key_values", Type: "json-passthrough", Description: "", IsMustExist: true},
			{NameId: "row_per_page", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "page_index", Type: "int64", Description: "", IsMustExist: true},
			{NameId: "is_deleted", Type: "bool", Description: "", IsMustExist: false},
		}, task_management_handler.CustomerMeterList, nil, table.Manager.StandardOperationResponsePossibility["list"], []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER_METER.LIST"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.CreateBulkBase64.CMS",
		"Upload base64 encoded file (csv or Excel) to create multiple customers in the system with validated information.",
		"/v1/customer/create_bulk/base64", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
			{NameId: "content_type", Type: "string", Description: "Content type of the file (csv or excel)", IsMustExist: true},
		}, task_management_handler.CustomerCreateBulkBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"CUSTOMER.UPLOAD"}, 0, "default",
	)

	anAPI.NewEndPoint("Customer.ConstructionTask.CreateBulkBase64.CMS",
		"Upload base64 encoded file (csv or Excel) to create multiple construction tasks in the system with validated information.",
		"/v1/task/construction/create_bulk/base64", "POST", api.EndPointTypeHTTPJSON, http.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "content_base64", Type: "string", Description: "File content in Base64", IsMustExist: true},
			{NameId: "content_type", Type: "string", Description: "Content type of the file (csv or excel)", IsMustExist: true},
		}, task_management_handler.CustomerConstructionTaskCreateBulkBase64, nil, nil, []api.DXAPIEndPointExecuteFunc{
			self.ModuleSelf.MiddlewareRequestRateLimitCheck,
			self.ModuleSelf.MiddlewareUserLoggedAndPrivilegeCheck,
		}, []string{"TASK.CONSTRUCTION.UPLOAD"}, 0, "default",
	)

}

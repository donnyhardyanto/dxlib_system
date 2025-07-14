package module_instance

import (
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/utils"
	"net/http"
)

// VersionHandler handles the /version endpoint request
func VersionHandler(aepr *api.DXAPIEndPointRequest) (err error) {
	// Return the version information as JSON
	aepr.WriteResponseAsJSON(http.StatusOK, nil, utils.JSON{
		"api_version":        1,
		"client_version":     1,
		"api_built":          1,
		"client_build":       1,
		"force_update_client": true,
	})
	return nil
}
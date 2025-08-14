package user_management

import (
	"github.com/donnyhardyanto/dxlib/api"
)

func DefineAPIEndPoints(anAPI *api.DXAPI) {
	defineAPITestUploadDownloadFile(anAPI)
	defineAPIRole(anAPI)
	defineAPIUserRoleMembership(anAPI)
	defineAPIPrivilege(anAPI)
	defineAPIRolePrivilege(anAPI)
	defineAPIOrganization(anAPI)
	defineAPIOrganizationRoles(anAPI)
	defineAPIUser(anAPI)
}

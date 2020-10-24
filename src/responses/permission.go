package responses

import (
	"worko.tech/iam/src/models"
)

type permissionResponse struct {
	AccessLevel	string `json:"accessLevel"`
}

func NewPermissionResponse(permission models.WorkspacePermission) *permissionResponse {
	response := new(permissionResponse)

	response.AccessLevel = permission.String()
	return response
}

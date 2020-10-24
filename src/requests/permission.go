package requests

import (
	"errors"

	"worko.tech/iam/src/models"

	"github.com/labstack/echo/v4"
)

type SetPermissionRequest struct {
	AccessLevel 	 	string `json:"accessLevel" validate:"required"`
	Permission 			models.WorkspacePermission
}
func (request *SetPermissionRequest) Bind(ctx echo.Context) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if err := ctx.Validate(request); err != nil {
		return err
	}

	request.Permission = models.GetWorkspacePermission(request.AccessLevel)
	if request.Permission == models.WsPermissionError {
		return errors.New("missing or wrong formatted access level")
	}
	return nil
}

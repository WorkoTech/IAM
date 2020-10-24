package requests;

import (
    "github.com/labstack/echo/v4"
)

type AccessRequest struct {
    Resource string `json:"resource" validate:"required"`
    Action string `json:"action" validate:"required"`
    WorkspaceId int `json:"workspaceId"`
    FileSize int `json:"fileSize"`
}

func (request *AccessRequest) Bind(ctx echo.Context) error {
    if err := ctx.Bind(request); err != nil {
        return err
    }
    if err := ctx.Validate(request); err != nil {
        return err
    }
    return nil
}

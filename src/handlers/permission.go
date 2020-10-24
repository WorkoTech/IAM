package handlers

import (
	"net/http"

	"worko.tech/iam/src/responses"
	"worko.tech/iam/src/utils"
	"worko.tech/iam/src/requests"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetUserPermissionOnWorkspace(ctx echo.Context) error {
	workspaceId := ctx.Param("workspaceId")
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)

	log.Info().Msgf("GetUserPermissionOnWorkspace user (%v) workspace (%v)", userId, workspaceId)
	permission, err := h.permissionDao.GetUserPermissionOnWorkspace(userId, workspaceId)
	if err != nil {
		log.Error().Msgf("Unable to retrieve user (%v) permissions on workspace (%v) : %v", userId, workspaceId, err.Error())
		return err
	}

	log.Info().Msgf("Successfully retrieved user (%v) permissions on workspace %v : %v", userId, workspaceId, permission)
	return ctx.JSON(http.StatusOK, responses.NewPermissionResponse(permission))
}

func (h *Handler) SetUserPermissionOnWorkspace(ctx echo.Context) error {
	workspaceId := ctx.Param("workspaceId")
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)

	log.Info().Msgf("SetUserPermissionOnWorkspace")
	req := &requests.SetPermissionRequest{}
	if err := req.Bind(ctx); err != nil {
		log.Error().Msgf("Unable to parse request : %v", err.Error())
        return ctx.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }
	log.Info().Msgf("Set user (%v) permission (%v) on workspace (%v)", userId, req.Permission, workspaceId)

	err := h.permissionDao.SetUserPermissionOnWorkspace(userId, workspaceId, req.Permission)
	if err != nil {
		log.Error().Msgf("Unable to set user (%v) permission (%v) on workspace (%v) : %v", userId, req.Permission, workspaceId, err.Error())
		return ctx.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	log.Info().Msgf("Successfuly set user (%v) permission (%v) on workspace (%v)", userId, req.Permission, workspaceId)
	return ctx.JSON(http.StatusCreated, responses.NewPermissionResponse(req.Permission))
}

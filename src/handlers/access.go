package handlers

import (
    "net/http"

    "worko.tech/iam/src/externals"
    "worko.tech/iam/src/requests"
    "worko.tech/iam/src/utils"

    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo/v4"
    "github.com/rs/zerolog/log"
)

func (h *Handler) GrantAccess(context echo.Context) error {
    // Get userId from JWT token
    user := context.Get("user").(*jwt.Token)

    // Bind request
    req := &requests.AccessRequest{}
    if err := req.Bind(context); err != nil {
        log.Error().Msgf("Unable to parse request : %v", err.Error())
        return context.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }

    var billing *externals.Billing
    usageErr := error(nil)
    // Retrieve Billing usage
    if req.WorkspaceId != 0 {
        billing, usageErr = externals.RetrieveWorkspaceUsage(user, req.WorkspaceId)
        if usageErr != nil {
            return context.JSON(http.StatusInternalServerError, utils.NewError(usageErr))
        }
    } else {
        billing, usageErr = externals.RetrieveUserUsage(user)
        if usageErr != nil {
            return context.JSON(http.StatusInternalServerError, utils.NewError(usageErr))
        }
    }

    // Get limit for the specific resource
    limit := 0
    for _, item := range billing.Offer.Items {
        if item.Resource == req.Resource {
            limit = item.Limit
        }
    }


    type EnforcerRequest struct {
        Limit int
        Usage externals.Usage
        FileSize int
    }
    m := EnforcerRequest{limit, billing.Usage, req.FileSize}
    log.Info().Msgf("m : %v", m)

    // Resolve PEP
    granted, err := h.enforcer.Enforce(m, req.Resource, req.Action)
    if err != nil {
        log.Error().Msgf("Unable to enforce request: %v", err.Error())
        return context.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    // Grant access or not
    if granted == true {
        return context.NoContent(http.StatusOK)
    }
    return context.JSON(http.StatusForbidden, utils.AccessForbidden())
}

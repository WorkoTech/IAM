package handlers

import (
    "net/http"
    "errors"

    "worko.tech/iam/src/utils"

    "github.com/labstack/echo/v4"
    "github.com/rs/zerolog/log"
)

func (h *Handler) Health(context echo.Context) error {
    goodHealth := h.healthDao.Health()

    log.Debug().Msgf("Health (%v)", goodHealth)
    if !goodHealth {
        return context.JSON(http.StatusInternalServerError, utils.NewError(errors.New("Unable to reach databases")))
    }
    return context.NoContent(http.StatusOK)
}

func (h *Handler) Ping(context echo.Context) error {
    return context.NoContent(http.StatusOK)
}

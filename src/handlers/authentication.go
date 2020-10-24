package handlers

import (
    "strconv"
    "net/http"
    "time"
    "os"

    "worko.tech/iam/src/requests"
    "worko.tech/iam/src/responses"
    "worko.tech/iam/src/models"
    "worko.tech/iam/src/externals"
    "worko.tech/iam/src/utils"

    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo/v4"
    "github.com/rs/zerolog/log"
)

func (h *Handler) Login(context echo.Context) error {
    req := &requests.LoginRequest{}
    if err := req.Bind(context); err != nil {
        log.Warn().Msgf("Unable to proccess login request (%v)", err.Error())
        return context.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }

    user, err := h.userDao.GetByEmail(req.Email)
    if err != nil {
        log.Warn().Msgf("Unable to retrieve user that tried to login (%v)", req.Email)
        return context.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    if user == nil {
        log.Warn().Msgf("Unable to retrieve user that tried to login (%v)", req.Email)
        return context.JSON(http.StatusForbidden, utils.AccessForbidden())
    }
    if !user.CheckPassword(req.Password) {
        log.Warn().Msgf("Unable to login user (%v): Not valid password", req.Email)
        return context.JSON(http.StatusForbidden, utils.AccessForbidden())
    }

    // Create token
    token := jwt.New(jwt.SigningMethodHS256)
    // Set claims
    claims := token.Claims.(jwt.MapClaims)
    claims["userId"] = strconv.FormatUint(uint64(user.ID), 10)
    claims["email"] = user.Email

    expirationTime := time.Now().Add(time.Hour * 72).Unix()
    claims["exp"] = expirationTime

    // Generate encoded token and send it as response.
    encodedToken, err := token.SignedString([]byte(os.Getenv("IAM_JWT_SIGNED_SECRET")))
    if err != nil {
        log.Error().Msgf("Unable to create token (%v)", err.Error())
        return err
    }

    // Log gamification event
    externals.LogLoginEvent("Bearer " + encodedToken)

    log.Info().Msgf("User (%v) successfuly logged in", req.Email)
    return context.JSON(http.StatusOK, responses.NewAccessUserResponse(user, "Bearer " + encodedToken))
}


func (h *Handler) Register(ctx echo.Context) error {
    var user models.User

    req := &requests.RegisterRequest{}
    if err := req.Bind(ctx, &user); err != nil {
        return ctx.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }
    log.Info().Msgf("Registering user (%v)", user.Email)
    if err := h.userDao.Create(&user); err != nil {
        log.Info().Msgf("Unable to registe user (%v)", user.Email)
        return ctx.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }


    // Create token
    token := jwt.New(jwt.SigningMethodHS256)
    // Set claims
    claims := token.Claims.(jwt.MapClaims)
    claims["userId"] = strconv.FormatUint(uint64(user.ID), 10)
    claims["email"] = user.Email

    expirationTime := time.Now().Add(time.Hour * 72).Unix()
    claims["exp"] = expirationTime

    // Generate encoded token and send it as response.
    encodedToken, err := token.SignedString([]byte(os.Getenv("IAM_JWT_SIGNED_SECRET")))
    if err != nil {
        log.Error().Msgf("Unable to create token (%v)", err.Error())
        return err
    }

    err = externals.CreateUserProfil("Bearer " + encodedToken)
    if err != nil {
        log.Error().Msgf("Unable to create user profil (%v)", err.Error())
    }

    log.Info().Msgf("User (%v) successfully registered", user.Email)
    return ctx.JSON(http.StatusCreated, responses.NewUserResponse(&user))
}

func (h *Handler) GetUserByEmail(ctx echo.Context) error {
    userEmail := ctx.Param("email")

    log.Info().Msgf("GetUserByEmail (%v)", userEmail)

    user, err := h.userDao.GetByEmail(userEmail)
    if err != nil {
        log.Warn().Msgf("Unable to retrieve user (%v) by email: %v", userEmail, err.Error())
        return ctx.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    if user == nil {
        log.Warn().Msgf("Unable to retrieve user (%v) by email: Not Found", userEmail)
        return ctx.JSON(http.StatusNotFound, utils.NotFound())
    }

    log.Info().Msgf("Successfully retrieved user (%v) by email", userEmail)
    return ctx.JSON(http.StatusOK, responses.NewUserResponse(user))
}

func (h *Handler) GetCurrentUser(ctx echo.Context) error {
    userToken := ctx.Get("user").(*jwt.Token)
    claims := userToken.Claims.(jwt.MapClaims)
    rawUserId := claims["userId"].(string)
    log.Info().Msgf("GetCurrentUser (%v)", rawUserId)

    userId, err := strconv.ParseInt(rawUserId, 10, 64)
    if err != nil {
        log.Warn().Msgf("Unable to parse userId (%v)", err.Error())
        return ctx.JSON(http.StatusBadRequest, utils.NewError(err))
    }

    user, err := h.userDao.GetById(userId)
    if err != nil {
        log.Error().Msgf("Unable to retrieve user (%v)", err.Error())
        return ctx.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    if user == nil {
        log.Warn().Msgf("User (%v) not found", userId)
        return ctx.JSON(http.StatusNotFound, utils.NotFound())
    }

    log.Info().Msgf("Successfully retrieve current user %v", userId)
    return ctx.JSON(http.StatusOK, responses.NewUserResponse(user))
}

func (h *Handler) SearchUsers(ctx echo.Context) error {
    log.Info().Msg("SearchUsers")
    req := &requests.SearchUsersRequest{}
    if err := req.Bind(ctx); err != nil {
        log.Warn().Msgf("Unable to process search user process (%v)", err.Error())
        return ctx.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }

    users, err := h.userDao.GetByIds(req.UserIds)
    if err != nil {
        log.Error().Msgf("Unabe to retrieve users (%v)", err.Error())
        return ctx.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    if users == nil {
        log.Warn().Msgf("Ids not found (%v)", req.UserIds)
        return ctx.JSON(http.StatusNotFound, utils.NotFound())
    }

    log.Info().Msgf("Sucessfuly searched user (%v)", req.UserIds)
    return ctx.JSON(http.StatusOK, responses.NewUsersResponse(users))
}

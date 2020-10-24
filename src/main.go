package main

import (
    "os"
    "time"

    "worko.tech/iam/src/router"
    "worko.tech/iam/src/handlers"
    "worko.tech/iam/src/utils"
    "worko.tech/iam/src/db"
    "worko.tech/iam/src/dao"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/jinzhu/gorm"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/casbin/casbin/v2"
)

func main() {
    e := newEcho()

    if os.Getenv("IAM_ENVIRONMENT") == "PRODUCTION" {
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    } else {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }
    log.Logger = log.Output(zerolog.ConsoleWriter{
        Out: os.Stdout,
        TimeFormat: time.RFC3339,
    }).With().Timestamp().Logger()

    logBanner()

    postgresConn := db.ConnectToPostgres()
    redisConn := db.ConnectToRedis()

    h := newHandler(postgresConn, redisConn)
    router.New(e, h)

    e.Start(":" + utils.GetEnv("IAM_PORT", "3002"))
}

func newEcho() (*echo.Echo){
    // Initialize echo
    e := echo.New()
    e.HideBanner = true

    // Midleware
    e.Pre(middleware.RemoveTrailingSlash())
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "[${time_rfc3339}] [${host} ${remote_ip}] ${method} ${uri} (${status})\n",
    }))
    e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{os.Getenv("IAM_CORS_DOMAIN")},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
        AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
    }))

    // Validator
    e.Validator = router.NewValidator()

    return e
}

func newHandler(postgres *gorm.DB, redis *db.RedisClient) *handlers.Handler {
    userDao := dao.NewUserDao(postgres)
    permissionDao := dao.NewPermissionDao(redis)
    healthDao := dao.NewHealthDao(postgres, redis)
    enforcer, err := casbin.NewEnforcer("./abac/model.conf", "./abac/policy.csv")
    if err != nil {
        log.Fatal().Err(err).Msgf("error: adapter")
    }

    return handlers.NewHandler(*userDao, *permissionDao, *healthDao, *enforcer)
}

func logBanner() {
    log.Info().Msg("                  ___           ___")
    log.Info().Msg("    ___          /  /\\         /__/\\")
    log.Info().Msg("   /  /\\        /  /::\\       |  |::\\")
    log.Info().Msg("  /  /:/       /  /:/\\:\\      |  |:|:\\")
    log.Info().Msg(" /__/::\\      /  /:/~/::\\   __|__|:|\\:\\")
    log.Info().Msg(" \\__\\/\\:\\__  /__/:/ /:/\\:\\ /__/::::| \\:\\")
    log.Info().Msg("    \\  \\:\\/\\ \\  \\:\\/:/__\\/ \\  \\:\\~~\\__\\/")
    log.Info().Msg("     \\__\\::/  \\  \\::/       \\  \\:\\")
    log.Info().Msg("     /__/:/    \\  \\:\\        \\  \\:\\")
    log.Info().Msg("     \\__\\/      \\  \\:\\        \\  \\:\\")
    log.Info().Msg("                 \\__\\/         \\__\\/")
    log.Info().Msgf("Starting Worko IAM on %v", utils.GetEnv("IAM_PORT", "3002"))
}

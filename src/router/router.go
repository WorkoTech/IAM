package router

import (
    "os"

    "worko.tech/iam/src/handlers"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func New(e *echo.Echo, h *handlers.Handler) {
    // Endpoints

    // Liviness and readiness
    e.GET("/ping", h.Ping)
    e.GET("/health", h.Health)

    // Authentication
    e.POST("/register", h.Register)
    e.POST("/login", h.Login)

    // User search API
    e.POST("/users/search", h.SearchUsers)
    e.GET("/users/:email", h.GetUserByEmail)

    // Get user from token
    userGroup := e.Group("/user");
    userGroup.Use(middleware.JWT([]byte(os.Getenv("IAM_JWT_SIGNED_SECRET"))))
    userGroup.GET("", h.GetCurrentUser)

    // Workspace Permissions
    protectedGroup := e.Group("/permission");
    protectedGroup.Use(middleware.JWT([]byte(os.Getenv("IAM_JWT_SIGNED_SECRET"))))
    protectedGroup.GET("/workspace/:workspaceId", h.GetUserPermissionOnWorkspace)
    protectedGroup.POST("/workspace/:workspaceId", h.SetUserPermissionOnWorkspace)

    // Access
    accessGroup := e.Group("/access");
    accessGroup.Use(middleware.JWT([]byte(os.Getenv("IAM_JWT_SIGNED_SECRET"))))
    accessGroup.POST("", h.GrantAccess)
}

package handlers

import (
    "worko.tech/iam/src/dao"

    "github.com/casbin/casbin/v2"
)

type Handler struct {
    userDao 	 	dao.UserDao
    permissionDao 	dao.PermissionDao
    healthDao       dao.HealthDao
    enforcer        casbin.Enforcer
}

func NewHandler(
    userDao dao.UserDao,
    permissionDao dao.PermissionDao,
    healthDao dao.HealthDao,
    enforcer casbin.Enforcer,
) *Handler {
    return &Handler{
        userDao: userDao,
        permissionDao: permissionDao,
        healthDao: healthDao,
        enforcer: enforcer,
    }
}

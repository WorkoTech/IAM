package models

import (
)

type WorkspacePermission int
const (
    WsPermissionNone WorkspacePermission = iota
    WsPermissionUser
    WsPermissionReferent
    WsPermissionCreator
    WsPermissionError
)

func (perm *WorkspacePermission) String() string {
    switch *perm {
    case WsPermissionCreator:
        return "CREATOR"
    case WsPermissionUser:
        return "USER"
    case WsPermissionReferent:
        return "REFERENT"
    case WsPermissionNone:
        return "NONE"
    }
    return ""
}

func GetWorkspacePermission(perm string) WorkspacePermission {
    switch perm {
    case "USER":
        return WsPermissionUser
    case "REFERENT":
        return WsPermissionReferent
    case "CREATOR":
        return WsPermissionCreator
    case "NONE":
        return WsPermissionNone
    }
    return WsPermissionError
}

package dao

import (
	"errors"

	"worko.tech/iam/src/db"
	"worko.tech/iam/src/models"
)

type PermissionDao struct {
	redis *db.RedisClient
}

func NewPermissionDao(redis *db.RedisClient) *PermissionDao {
	return &PermissionDao{
		redis: redis,
	}
}

func (dao *PermissionDao) SetUserPermissionOnWorkspace(userId string, workspaceId string, permission models.WorkspacePermission) error {
	if len(userId) == 0 && len(workspaceId) == 0 {
		return errors.New("userId and workspaceId cannot be empty")
	}

	key := userId + ":" + workspaceId
	perm := permission.String()

	status, err := dao.redis.SetValue(key, perm)
	if err != nil {
		return err
	}
	if !status {
		return errors.New("Error while querying Redis")
	}

	return nil
}

func (dao *PermissionDao) GetUserPermissionOnWorkspace(userId string, workspaceId string) (models.WorkspacePermission, error) {
	key := userId + ":" + workspaceId

	value, err := dao.redis.GetValue(key)
	if err != nil {
		return models.WsPermissionNone, err
	}

	if len(value) == 0 {
		return models.WsPermissionNone, err
	}

	return models.GetWorkspacePermission(value), nil
}

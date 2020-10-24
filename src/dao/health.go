package dao

import (
    "worko.tech/iam/src/db"

    "github.com/jinzhu/gorm"
)

type HealthDao struct {
    postgres *gorm.DB
    redis *db.RedisClient
}

func NewHealthDao(postgres *gorm.DB, redis *db.RedisClient) *HealthDao {
    return &HealthDao {
        postgres: postgres,
        redis: redis,
    }
}

type PostgresPingResult struct {
    reachable int
}
func (dao *HealthDao) Health() (bool) {
    var postgresStatus PostgresPingResult

    dao.postgres.Raw("SELECT 1 + 1 as reachable",).Scan(&postgresStatus)
    redisConnected, _ := dao.redis.Ping()

    if (postgresStatus.reachable > 0 && redisConnected) {
        return true
    }
    return false;
}

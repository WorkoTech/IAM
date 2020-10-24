package db

import (
	"time"

	"worko.tech/iam/src/utils"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
    conn 	 		*redis.Client
}
func ConnectToRedis() *RedisClient {
	var redisClient RedisClient

	redisClient.conn = redis.NewClient(&redis.Options{
		Addr: utils.GetEnv("REDIS_HOST", "localhost") + ":" + utils.GetEnv("REDIS_PORT", "5432"),
		PoolSize: 100,
		MaxRetries: 2,
		Password: "",
		DB: 0,
	})

	var err error
	redisConnected := false
	for redisConnected != true {
		redisConnected, err = redisClient.Ping()
		if !redisConnected {
			log.Warn().Msgf("Unable to connect to redis (%v). Retrying...", err.Error())
			time.Sleep(2000 * time.Millisecond)
		}
	}
	return &redisClient
}

func (client *RedisClient) GetValue(key string) (string, error) {
	log.Info().Msgf("REDIS Get %v", key)
	return client.conn.Get(key).Result()
}

func (client *RedisClient) SetValue(key string, value string) (bool, error) {
	log.Info().Msgf("REDIS Set %v:%v", key, value)
	err := client.conn.Set(key, value, 0).Err()
	return true, err
}

func (client *RedisClient) Ping() (bool, error){
	ping, err := client.conn.Ping().Result()
	if err == nil && len(ping) > 0 {
		log.Info().Msg("Sucessfuly pinged redis")
		return true, err;
	} else {
		log.Info().Msgf("Error while pinging redis: %v", err.Error())
		return false, err;
	}
}

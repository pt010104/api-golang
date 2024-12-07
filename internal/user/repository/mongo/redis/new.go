package redis

import (
	"github.com/pt010104/api-golang/pkg/redis"

	"github.com/pt010104/api-golang/pkg/log"
)

type implRedis struct {
	l     log.Logger
	redis redis.Client
}

func New(l log.Logger, redisClient redis.Client) implRedis {
	return implRedis{
		l:     l,
		redis: redisClient,
	}

}

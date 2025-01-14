package common

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var RedisClient redis.Client

type RedisModel struct {
	Host     string
	Port     string
	Password string
}

type IRedisConfig interface {
	Open() *error
}

func NewRedisConfig(model RedisModel) IRedisConfig {
	return RedisModel{
		Host:     model.Host,
		Port:     model.Port,
		Password: model.Password,
	}
}

func (r RedisModel) Open() *error {

	client, err := open(r)
	if err != nil {
		return err
	}
	RedisClient = *client

	return nil
}

func open(model RedisModel) (*redis.Client, *error) {

	client := redis.NewClient(&redis.Options{
		Addr:     model.Host + ":" + model.Port,
		Password: model.Password,
		DB:       0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, &err
	}

	return client, nil
}

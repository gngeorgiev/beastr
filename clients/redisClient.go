package clients

import (
	"log"

	"github.com/spf13/viper"
	"gopkg.in/redis.v4"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		log.Println(viper.GetString("redis_address"))
		client = redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_address"),
			Password: "",
			DB:       0,
		})

		pong, err := client.Ping().Result()
		log.Println(pong, err)
	}

	return client
}

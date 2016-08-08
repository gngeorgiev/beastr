package clients

import (
	"log"

	"github.com/cenk/backoff"
	"github.com/spf13/viper"
	"gopkg.in/redis.v4"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	return client
}

func StartRedisConnection() {
	connectOperation := func() error {
		newClient := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_address"),
			Password: "",
			DB:       0,
		})

		_, err := newClient.Ping().Result()
		if err != nil {
			return err
		}

		client = newClient
		return nil
	}

	b := backoff.NewExponentialBackOff()
	ticker := backoff.NewTicker(b)

	var err error
	for range ticker.C {
		if err = connectOperation(); err != nil {
			log.Println(err, "Retrying redis connection...")
			continue
		}

		ticker.Stop()
		break
	}

	if err != nil {
		log.Println(err, "Redis connection error %s")
	}
}

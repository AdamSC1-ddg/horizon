package horizon

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"net/url"
	"time"
)

func initRedis(app *App) {
	if app.config.RedisUrl == "" {
		return
	}

	redisUrl, err := url.Parse(app.config.RedisUrl)

	if err != nil {
		log.Panic(err)
	}

	app.redis = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         dialRedis(redisUrl),
		TestOnBorrow: pingRedis,
	}

	// test the connection
	c := app.redis.Get()
	defer c.Close()

	_, err = c.Do("PING")

	if err != nil {
		log.Panic(err)
	}
}

func dialRedis(redisUrl *url.URL) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisUrl.Host)
		if err != nil {
			return nil, err
		}

		if redisUrl.User == nil {
			return c, err
		}

		if pass, ok := redisUrl.User.Password(); ok {
			if _, err := c.Do("AUTH", pass); err != nil {
				c.Close()
				return nil, err
			}
		}

		return c, err
	}
}

func pingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("PING")
	return err
}

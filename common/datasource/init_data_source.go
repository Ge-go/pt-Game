package datasource

import (
	goredis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"ptc-Game/common/pkg/config"
	"ptc-Game/common/pkg/mysql"
	"ptc-Game/common/pkg/redis"
)

var DefaultDataSource *DataSources

type DataSources struct {
	DB          *gorm.DB
	RedisClient *goredis.Client
}

func NewDataSources(config *config.Config) (*DataSources, error) {
	db, err := mysql.New(config.MySQL.Default)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(config.Redis.Default)
	if err != nil {
		return nil, err
	}

	ds := &DataSources{
		DB:          db,
		RedisClient: redisClient,
	}

	DefaultDataSource = ds

	return ds, nil
}

func GetDataSource() *DataSources {
	return DefaultDataSource
}

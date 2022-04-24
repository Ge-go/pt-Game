package repositories

import (
	"github.com/go-redis/redis/v8"
	"gobasic/ptc-Game/common/datasource"
	"gorm.io/gorm"
)

type DocusignRepository interface {
}

//初始化
func NewDocusignRepository(ds *datasource.DataSources) DocusignRepository {
	return &docusignRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

//struct
type docusignRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

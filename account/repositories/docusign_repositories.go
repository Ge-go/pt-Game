package repositories

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"ptc-Game/common/datasource"
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

package repositories

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"ptc-Game/common/datasource"
)

type MaterialPermissionRepository interface {
}

func NewMaterialPermissionRepository(ds *datasource.DataMaterialSources) MaterialPermissionRepository {
	return &mterialPermissionRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type mterialPermissionRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

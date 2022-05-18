package repositories

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ptc-Game/account/datamodels"
	"ptc-Game/common/datasource"
)

type DocusignRepository interface {
	UpdateSignInfo(ctx context.Context, sign datamodels.StreamerSign) (int64, error)
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

func (d *docusignRepository) UpdateSignInfo(ctx context.Context, data datamodels.StreamerSign) (int64, error) {
	rs := d.db.WithContext(ctx).Where(&datamodels.StreamerSign{UserId: data.UserId, EnvelopeId: data.EnvelopeId}).Updates(&data)
	return rs.RowsAffected, errors.Wrapf(rs.Error, "update sign info failed")
}

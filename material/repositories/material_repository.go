package repositories

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ptc-Game/common/datasource"
	"ptc-Game/material/datamodels"
	"ptc-Game/material/web/viewmodels"
	UserInfo "ptc-Game/userinfo/datamodels"
)

type MaterialRepository interface {
	//MateHome 素材库首页
	MateHome(ctx context.Context, uid int64, req viewmodels.MateHomeReq) ([]datamodels.ConfirmUpload, int64, error)
}

func NewMaterialRepository(dm *datasource.DataMaterialSources, ds *datasource.DataSources) MaterialRepository {
	return &materialRepository{
		dbm:   dm.DB,
		dbs:   ds.DB,
		redis: dm.RedisClient,
	}
}

type materialRepository struct {
	dbm   *gorm.DB
	dbs   *gorm.DB
	redis *redis.Client
}

func (m *materialRepository) MateHome(ctx context.Context, uid int64, req viewmodels.MateHomeReq) ([]datamodels.ConfirmUpload, int64, error) {
	// 查询该用户的等级和语种
	var (
		mateType   string
		mateModule string
		user       UserInfo.User
		total      int64
		level      int64
	)

	if err := m.dbs.WithContext(ctx).Model(UserInfo.User{}).Select("level", "language", "email").Where("Id = ?", uid).
		Find(&user).Error; err != nil {
		errors.Wrapf(err, "FindUserLevelAndLanguage failed")
	}

	//增加白名单
	m.IsWhiteUser(ctx, user.Email, uid)

}

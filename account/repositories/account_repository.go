package repositories

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ptc-Game/account/datamodels"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/conf"
	"ptc-Game/common/datasource"
	"ptc-Game/common/response"
	"strconv"
	"time"
)

const (
	//cache前缀： 业务 + 模块 + key
	AccountEmailPrefix    = conf.KeyPrefix + "-streamer-system:account:email:%v" //邮箱
	AccountCaptcha        = conf.KeyPrefix + "-streamer-system:account:captcha"  //验证码
	AccountTokenPrefix    = conf.KeyPrefix + "-streamer-system:account:token:%v" //token
	AccountRegisterPrefix = conf.KeyPrefix + "-streamer-system:account:registerinfo:%v"
	AuthorLevel           = conf.KeyPrefix + "-streamer-system:authorlevel:level" //作者等级粉丝区间zset
)

type AccountRepository interface {
	FindByName(ctx context.Context, name string) (*datamodels.User, error)
	StoreEmailCode(ctx context.Context, email string, code string) error
	GetTagList(ctx context.Context) ([]viewmodels.UserTag, error)
	Register(ctx context.Context, data viewmodels.RegisterInfo) error
	IsEmailExist(ctx context.Context, email string) (bool, error)
	VerifyEmailCode(ctx context.Context, email string, emailCode string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*datamodels.User, error)
	ResetPassword(ctx context.Context, email string, password string) error
	ClearEmailCode(ctx context.Context, email string) error
	FindById(ctx context.Context, uid uint) (*datamodels.User, error)
	UpdateYoutubeAccount(ctx context.Context, uid string, userName string, channelId string) (int, error)
	SetCaptcha(ctx context.Context, id, content string) error
	GetCaptcha(ctx context.Context, id string) (string, error)
	FindBySub(ctx context.Context, sub string) (*datamodels.User, error)
	FindMinorCertByName(ctx context.Context, name string) (*datamodels.MinorCertConfig, error)
	SetToken(ctx context.Context, sub string, token string) error
	GetToken(ctx context.Context, sub string) (string, error)
	SaveRegisterInfo(ctx context.Context, sub string, data string) error
	GetRegisterInfo(ctx context.Context, sub string) (string, error)
	GetAuthorLevel(ctx context.Context, fans int) ([]string, error)
	UserLoginRecord(userId int)
}

func NewAccountRepository(ds *datasource.DataSources) AccountRepository {
	return &accountRepository{
		db:    ds.DB,
		redis: ds.RedisClient,
	}
}

type accountRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func (a *accountRepository) FindByName(ctx context.Context, name string) (*datamodels.User, error) {
	user := &datamodels.User{}
	err := a.db.WithContext(ctx).Where(&datamodels.User{UserName: name}).First(user).Error
	return user, err
}

func (a *accountRepository) StoreEmailCode(ctx context.Context, email string, code string) error {
	cacheKey := fmt.Sprintf(AccountEmailPrefix, email)

	val := a.redis.Exists(ctx, cacheKey).Val()
	if val == 1 {
		return response.ErrEmailIsVerifying
	}

	err := a.redis.Set(ctx, cacheKey, code, 60*time.Second).Err()
	return errors.Wrapf(err, "StoreEmailCode failed")
}

func (a *accountRepository) GetTagList(ctx context.Context) ([]viewmodels.UserTag, error) {
	tags := []viewmodels.UserTag{}

	err := a.db.WithContext(ctx).Table("streamer_tag").Select("id,tag_name").Find(&tags).Error
	return tags, errors.Wrapf(err, "GetTagList failed")
}

func (a *accountRepository) Register(ctx context.Context, data viewmodels.RegisterInfo) error {
	user := &datamodels.User{
		UserName:  data.UserName,
		Password:  data.Password,
		Email:     data.Email,
		Kind:      data.Kind,
		Gender:    data.Gender,
		Country:   data.Country,
		Region:    data.Region,
		Language:  data.Language,
		IsAdult:   data.IsAdult,
		Birthday:  data.Birthday,
		GoogleSub: data.Sub,
		IsLocked:  0,
		Level:     data.Level,
	}
	if data.IsEuropean == true { //保存签名
		user.IsSigned = 1
	}

	result := a.db.WithContext(ctx).Create(&user)

	if result.RowsAffected > 0 { //注册成功则删除缓存
		cacheKey := fmt.Sprintf(AccountRegisterPrefix, data.Sub)
		a.redis.Del(ctx, cacheKey)
	}

	return errors.Wrapf(result.Error, "Register failed")
}

func (a *accountRepository) IsEmailExist(ctx context.Context, email string) (bool, error) {
	var user datamodels.User
	var total int64

	result := a.db.WithContext(ctx).Model(&user).Where(datamodels.User{Email: email}).Count(&total)
	return total >= 1, errors.Wrapf(result.Error, "GetWonderfulWorkList failed")
}

func (a *accountRepository) VerifyEmailCode(ctx context.Context, email string, emailCode string) (bool, error) {
	cacheKey := fmt.Sprintf(AccountEmailPrefix, email)
	// check whether email code was in redis
	isExist := a.redis.Exists(ctx, cacheKey).Val()
	if isExist == 0 {
		return false, nil
	}

	//存在则去获取值
	cmd := a.redis.Get(ctx, cacheKey)
	if cmd.Err() != nil {
		return false, errors.Wrapf(cmd.Err(), "VerifyEmailCode failed")
	}
	if cmd.Val() != emailCode {
		return false, nil
	}

	return true, nil
}

func (a *accountRepository) FindByEmail(ctx context.Context, email string) (*datamodels.User, error) {
	user := &datamodels.User{}
	err := a.db.WithContext(ctx).Where(&datamodels.User{Email: email}).First(user).Error
	return user, errors.Wrapf(err, "FindByEmail failed")
}

func (a *accountRepository) ResetPassword(ctx context.Context, email string, password string) error {
	err := a.db.WithContext(ctx).Model(&datamodels.User{}).Where("email = ?", email).Updates(&datamodels.User{
		Password: password,
	}).Error

	return errors.Wrapf(err, "ResetPassword failed")
}

func (a *accountRepository) ClearEmailCode(ctx context.Context, email string) error {
	cacheKey := fmt.Sprintf(AccountEmailPrefix, email)
	cmd := a.redis.Del(ctx, cacheKey)
	return cmd.Err()
}

func (a *accountRepository) FindById(ctx context.Context, uid uint) (*datamodels.User, error) {
	user := &datamodels.User{}
	err := a.db.WithContext(ctx).Model(&datamodels.User{}).Where("id = ?", uid).First(user).Error
	return user, errors.Wrapf(err, "FindById failed")
}

func (a *accountRepository) UpdateYoutubeAccount(ctx context.Context, uid string, userName string, channelId string) (int, error) {
	sql := a.db.WithContext(ctx).Model(&datamodels.User{}).Where("id = ? AND user_name = ?", uid, userName).Update("youtube_account", channelId)
	return int(sql.RowsAffected), errors.Wrapf(sql.Error, "UpdateYoutubeAccount failed")
}

func (a *accountRepository) SetCaptcha(ctx context.Context, id, content string) error {
	err := a.redis.HSetNX(ctx, AccountCaptcha, id, content).Err()
	return errors.Wrapf(err, "SetCaptcha failed")
}

func (a *accountRepository) GetCaptcha(ctx context.Context, id string) (string, error) {
	res := a.redis.HGet(ctx, AccountCaptcha, id)
	if res.Err() != nil {
		return "", errors.Wrapf(res.Err(), "GetCaptcha failed")
	}
	err := a.redis.HDel(ctx, AccountCaptcha, id).Err()
	if err != nil {
		return "", err
	}
	return res.Val(), nil
}

func (a *accountRepository) FindBySub(ctx context.Context, sub string) (*datamodels.User, error) {
	user := &datamodels.User{}
	result := a.db.WithContext(ctx).Where(&datamodels.User{GoogleSub: sub}).First(user)

	return user, result.Error
}

func (a *accountRepository) FindMinorCertByName(ctx context.Context, name string) (*datamodels.MinorCertConfig, error) {
	minorCert := &datamodels.MinorCertConfig{}
	err := a.db.WithContext(ctx).Where(&datamodels.MinorCertConfig{
		ShortName: name,
	}).First(&minorCert).Error
	return minorCert, errors.Wrapf(err, "FindMinorCertByName failed")
}

func (a *accountRepository) SetToken(ctx context.Context, sub string, token string) error {
	cacheKey := fmt.Sprintf(AccountTokenPrefix, sub)

	err := a.redis.Set(ctx, cacheKey, token, 30*time.Minute).Err()
	return errors.Wrapf(err, "redis Set google token failed")
}

func (a *accountRepository) GetToken(ctx context.Context, sub string) (string, error) {
	cacheKey := fmt.Sprintf(AccountTokenPrefix, sub)

	res := a.redis.Get(ctx, cacheKey)

	if res.Err() != nil {
		return "", res.Err()
	}

	return res.Val(), nil
}

func (a *accountRepository) SaveRegisterInfo(ctx context.Context, sub string, data string) error {
	cacheKey := fmt.Sprintf(AccountRegisterPrefix, sub)

	err := a.redis.Set(ctx, cacheKey, data, 30*time.Minute).Err()
	return errors.Wrapf(err, "redis Set redis DATA failed")
}

func (a *accountRepository) GetRegisterInfo(ctx context.Context, sub string) (string, error) {
	cacheKey := fmt.Sprintf(AccountRegisterPrefix, sub)

	cmd := a.redis.Get(ctx, cacheKey)

	if cmd.Err() != nil {
		return "", errors.Wrapf(cmd.Err(), "get register info failed")
	}

	return cmd.Val(), nil
}

func (a *accountRepository) GetAuthorLevel(ctx context.Context, fans int) ([]string, error) {
	cacheKey := fmt.Sprintf(AuthorLevel)
	fansString := strconv.Itoa(fans) //大于等于
	cmd := a.redis.ZRevRangeByScore(ctx, cacheKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: fansString,
	})

	//错误不为空
	if cmd.Err() != nil {
		return nil, errors.Wrapf(cmd.Err(), "get auther level  failed")
	}
	return cmd.Val(), nil
}

func (a *accountRepository) UserLoginRecord(userId int) {
	a.db.Exec("INSERT INTO streamer_user_operation_record (user_id,type,time) VALUES (?,?,?)", userId, 1, time.Now().Unix())
}

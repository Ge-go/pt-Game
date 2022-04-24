package services

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gobasic/ptc-Game/account/repositories"
	"gobasic/ptc-Game/account/web/viewmodels"
	"gobasic/ptc-Game/common/pkg/captcha"
	"time"
)

type AccountService interface {
	GetRegisterInfo(ctx context.Context, sub string) (*viewmodels.RegisterInfo, error)
	Register(ctx context.Context, data viewmodels.RegisterInfo) error
	MinorCert(ctx context.Context, req viewmodels.MinorLimitReq) (*viewmodels.MinorLimitRsp, error)
	SaveRegisterInfo(ctx context.Context, info viewmodels.RegisterInfo) error
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{
		repo:       repo,
		CaptchaSvc: captcha.New(captcha.Digit),
	}
}

type accountService struct {
	repo       repositories.AccountRepository
	CaptchaSvc captcha.Captcha
}

func (a *accountService) SaveRegisterInfo(ctx context.Context, info viewmodels.RegisterInfo) error {
	jsonData, err := json.Marshal(info)
	if err != nil {
		return errors.Wrapf(err, "parse json failed")
	}

	if err = a.repo.SaveRegisterInfo(ctx, info.Sub, string(jsonData)); err != nil {
		return errors.Wrapf(err, "set redis data(RegisterInfo)failed")
	}

	return nil
}

//未成年验证
func (a *accountService) MinorCert(ctx context.Context, req viewmodels.MinorLimitReq) (*viewmodels.MinorLimitRsp, error) {
	var result viewmodels.MinorLimitRsp
	certInfo, err := a.repo.FindMinorCertByName(ctx, req.ShortName)

	if err != nil {
		return nil, err
	}

	if certInfo.AdultAge == 0 { //空置或者为0则默认为18
		certInfo.AdultAge = 18
	}
	//根据生日计算实际年龄
	birthday, timeErr := time.Parse("2006-01-02", req.Birthday)
	if timeErr != nil {
		return nil, timeErr
	}
	t := time.Now()
	age := t.Year() - birthday.Year()
	if birthday.Month() > t.Month() {
		age = age - 1
	}

	if birthday.Month() == t.Month() {
		if birthday.Day() > t.Day() {
			age = age - 1
		}
	}
	result.Age = age
	//判断是否成年,比设置的年龄大则成年
	if age >= certInfo.AdultAge {
		result.IsAdult = true
	} else {
		result.IsAdult = false
	}
	//判断是否欧盟标准
	if certInfo.IsEuropean == 1 || certInfo.IsEuropean == 2 {
		result.IsEuropean = true
	} else {
		result.IsEuropean = false
	}

	return &result, nil
}

func (a *accountService) Register(ctx context.Context, data viewmodels.RegisterInfo) error {
	data.Level = "1"

	return a.repo.Register(ctx, data)
}

func (a *accountService) GetRegisterInfo(ctx context.Context, sub string) (*viewmodels.RegisterInfo, error) {
	registerString, err := a.repo.GetRegisterInfo(ctx, sub)
	if err != nil {
		return nil, errors.Wrapf(err, "get redis data(RegisterInfo) failed")
	}

	var registerInfo viewmodels.RegisterInfo
	json.Unmarshal([]byte(registerString), &registerInfo)

	return &registerInfo, nil
}

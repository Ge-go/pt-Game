package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"ptc-Game/account/repositories"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/pkg/captcha"
	"ptc-Game/common/pkg/email"
	"ptc-Game/common/response"
	"ptc-Game/common/util"
	"time"
)

type AccountService interface {
	GetRegisterInfo(ctx context.Context, sub string) (*viewmodels.RegisterInfo, error)
	Register(ctx context.Context, data viewmodels.RegisterInfo) error
	MinorCert(ctx context.Context, req viewmodels.MinorLimitReq) (*viewmodels.MinorLimitRsp, error)
	SaveRegisterInfo(ctx context.Context, info viewmodels.RegisterInfo) error
	VerifyEmailCode(ctx context.Context, email, emailCode string) (bool, error)
	SendEmail(ctx context.Context, req viewmodels.VerifyEmailReq) (*viewmodels.VerifyEmailRsp, error)
	Login(ctx context.Context, data viewmodels.LoginReq) (*viewmodels.JwtToken, error)
	IsEmailExist(ctx context.Context, email string) (bool, error)
	ResetPassword(ctx context.Context, req viewmodels.ResetPasswordReq) error
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

func (a *accountService) ResetPassword(ctx context.Context, req viewmodels.ResetPasswordReq) error {
	hashPassword, err := util.GenHashPassword(req.Password)
	if err != nil {
		return err
	}
	return a.repo.ResetPassword(ctx, req.Email, hashPassword)
}

func (a *accountService) IsEmailExist(ctx context.Context, email string) (bool, error) {
	return a.repo.IsEmailExist(ctx, email)
}

func (a *accountService) Login(ctx context.Context, data viewmodels.LoginReq) (*viewmodels.JwtToken, error) {
	// verify captcha 存在redis,这样可以支持集群
	val, err := a.repo.GetCaptcha(ctx, data.CaptchaId)

	if err != nil || val != data.CaptchaValue {
		return nil, response.ErrInvalidCaptcha
	}

	user, err := a.repo.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, response.ErrInvalidPassword
	}

	// compare with user submit
	if !util.MatchPassword(data.Password, user.Password) {
		return nil, response.ErrInvalidPassword
	}

	// account has been locked
	if user.IsLocked == 1 {
		return &viewmodels.JwtToken{IsLocked: user.IsLocked}, nil
	}

	// put necessar user info to jwt payload
	payload := map[string]interface{}{"uid": user.Id, "region": user.Region}

	// sign jwt token
	token, err := util.SignJwtToken(payload)
	if err != nil {
		return nil, err
	}

	a.repo.UserLoginRecord(user.Id)

	return &viewmodels.JwtToken{Token: token, UserName: user.UserName}, nil
}

func (a *accountService) SendEmail(ctx context.Context, req viewmodels.VerifyEmailReq) (*viewmodels.VerifyEmailRsp, error) {
	randomDigit := util.RandomDigit()
	fmt.Println(randomDigit)
	err := a.repo.StoreEmailCode(ctx, req.Email, randomDigit)
	if err != nil {
		return nil, err
	}

	// send email to user
	subject := "Please verify your email"
	emailTemplate := fmt.Sprintf(`Your code verificationcode is: %v`, randomDigit)
	err = email.SendEmail(ctx, req.Email, subject, emailTemplate)
	if err != nil {
		return nil, err
	}

	return &viewmodels.VerifyEmailRsp{EmailCode: randomDigit}, nil
}

func (a *accountService) VerifyEmailCode(ctx context.Context, email, emailCode string) (bool, error) {
	return a.repo.VerifyEmailCode(ctx, email, emailCode)
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

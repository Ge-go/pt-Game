package controllers

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"google.golang.org/api/youtube/v3"
	"gorm.io/gorm"
	"ptc-Game/account/services"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/conf"
	"ptc-Game/common/pkg/logiclog"
	"ptc-Game/common/response"
	"ptc-Game/common/util"
)

type AccountController struct {
	Service        services.AccountService
	YoutubeService services.YoutubeService
}

// @Tags 账号管理 相关接口
// @Summary 账号注册
// @Description 账号注册
// @Router /api/v1/user/register [post]
// @Accept json
// @Produce json
// @Param param body viewmodels.RegisterReq true "请求参数"
// @Success 200 {object} response.Res "请求响应"
func (a *AccountController) Register(c iris.Context) {
	// bind param
	var req viewmodels.RegisterReq
	if err := c.ReadJSON(&req); err != nil {
		logiclog.CtxLogger(c).Warnf("bind err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}
	// validate param
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("param validate err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	//1 取回registerInfo
	registerInfo, err := a.Service.GetRegisterInfo(c.Request().Context(), req.Sub)
	if registerInfo == nil || err == redis.Nil {
		response.Send(c, response.GetMessage(response.NoYoutubeAuth, conf.EN), nil)
		return
	}
	if err != nil {
		logiclog.CtxLogger(c).Warnf("get registerInfo err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	//保存
	registerInfo.Gender = req.Gender
	registerInfo.Language = req.Language
	registerInfo.Kind = req.Kind

	//2 判断是否签署欧盟国家协议
	if registerInfo.IsEuropean == true {
		if req.IsEuropean == 0 {
			response.Send(c, response.GetMessage(response.NoAgreementSigned, conf.EN), nil)
			return
		}
	}

	//注册register
	err = a.Service.Register(c.Request().Context(), *registerInfo)
	if err != nil {
		logiclog.CtxLogger(c).Errorf("service err: (%+v)", err)
		response.Send(c, err, nil)
		return
	}

	response.Send(c, nil, nil)
}

// MinorCert
// @Tags 账号管理 相关接口
// @Summary 未成年人年龄验证和youtube频道粉丝数验证
// @Description 未成年人年龄验证和youtube频道粉丝数验证
// @Router /api/v1/user/minorCert [Post]
// @Accept json
// @Produce json
// @Param param query viewmodels.MinorLimitReq true "请求参数"
// @Success 200 {object} viewmodels.ResponseMinorCert "请求响应"
func (a *AccountController) MinorCert(c iris.Context) {
	//绑定数据
	var minorLimitReq viewmodels.MinorLimitReq
	if err := c.ReadJSON(&minorLimitReq); err != nil {
		logiclog.CtxLogger(c).Warnf("[MinorCert] required code parameter: %+v", err)
		response.Send(c, response.GetMessage(response.ErrInvalidParam, conf.ZH), nil)
		return
	}

	//验证数据
	_, err := govalidator.ValidateStruct(minorLimitReq)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("[MinorCert] param validate err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	//5 判断是否成年
	data, minorErr := a.Service.MinorCert(c.Request().Context(), minorLimitReq)
	if minorErr != nil {
		logiclog.CtxLogger(c).Warnf("[MinorCert] get is adult failed: %+v", minorErr)
		response.Send(c, minorErr, nil)
		return
	}

	if data.IsAdult == false {
		response.Send(c, response.GetMessage(response.NoAdult, conf.ZH), nil)
		return
	}

	var sub = uuid.NewString()

	if err = a.Service.SaveRegisterInfo(c.Request().Context(), viewmodels.RegisterInfo{
		Sub:        sub,
		Country:    minorLimitReq.ShortName,
		Birthday:   minorLimitReq.Birthday,
		IsAdult:    1,
		IsEuropean: data.IsEuropean,
	}); err != nil {
		logiclog.CtxLogger(c).Warnf("[MinorCert] save register info: %+v", err)
		response.Send(c, err, nil)
		return
	}

	data.Sub = sub

	response.Send(c, nil, data)

	return
}

// @Tags 账号管理 相关接口
// @Summary 邮箱校验
// @Description 邮箱校验（传emailCode，则校验邮箱验证码是否正确,不传emailCode 则发送验证码到指定邮箱）
// @Router /api/v1/user/verifyEmail [post]
// @Accept json
// @Produce json
// @Param param body viewmodels.VerifyEmailReq true "请求参数"
// @Success 200 {object} response.Res "请求响应"
func (a *AccountController) VerifyEmail(c iris.Context) {
	// bind param
	var req viewmodels.VerifyEmailReq
	if err := c.ReadJSON(&req); err != nil {
		logiclog.CtxLogger(c).Warnf("invalid param err:(%+v)", err)
		response.Send(c, response.GetMessage(response.ErrInvalidParam, conf.ZH), nil)
		return
	}
	// validate param
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("param validate err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	//verify email code if not empty
	if req.EmailCode != "" {
		isOk, verifyErr := a.Service.VerifyEmailCode(c.Request().Context(), req.Email, req.EmailCode)
		if verifyErr != nil {
			logiclog.CtxLogger(c).Warnf("param validate err:(%+v)", verifyErr)
			response.Send(c, verifyErr, nil)
			return
		}
		if !isOk {
			response.Send(c, response.GetMessage(response.ErrEmailCode, conf.ZH), nil)
			return
		}

		response.Send(c, nil, nil)
		return
	}

	_, err = a.Service.SendEmail(c.Request().Context(), req)
	if err != nil {
		logiclog.CtxLogger(c).Errorf("service err: (%+v)", err)
		response.Send(c, err, nil)
	}
	response.Send(c, nil, nil)
}

// @Tags 账号管理 相关接口
// @Summary 账号登陆
// @Description 账号登陆
// @Router /api/v1/user/login [post]
// @Accept json
// @Produce json
// @Param param body viewmodels.LoginReq true "请求参数"
// @Success 200 {object} viewmodels.JwtToken "请求响应"
func (a *AccountController) Login(context iris.Context) {
	// bind param
	var req viewmodels.LoginReq
	if err := context.ReadJSON(&req); err != nil {
		logiclog.CtxLogger(context).Warnf("bind err:(%+v)", err)
		response.Send(context, err, nil)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logiclog.CtxLogger(context).Warnf("param validate err:(%+v)", err)
		response.Send(context, err, nil)
		return
	}

	data, err := a.Service.Login(context.Request().Context(), req)
	if err != nil {
		logiclog.CtxLogger(context).Errorf("service err: (%+v)", err)
		response.Send(context, err, nil)
		return
	}

	if data.IsLocked == 1 {
		response.Send(context, response.GetMessage(response.AccountHasBeenFrozen, conf.ZH), nil)
		return
	}

	response.Send(context, nil, data)
}

// @Tags 账号管理 相关接口
// @Summary 密码找回
// @Description 密码找回
// @Router /api/v1/user/resetPassword [post]
// @Accept json
// @Produce json
// @Param param body viewmodels.ResetPasswordReq true "请求参数"
// @Success 200 {object} response.Res "请求响应"
func (a *AccountController) ResetPassword(context iris.Context) {
	//bind para
	var req viewmodels.ResetPasswordReq
	if err := context.ReadJSON(&req); err != nil {
		logiclog.CtxLogger(context).Warnf("bind err:(%+v)", err)
		response.Send(context, err, nil)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logiclog.CtxLogger(context).Warnf("param validate err:(%+v)", err)
		response.Send(context, err, nil)
		return
	}

	//confirm that twice password was the same
	if req.Password != req.ConfirmPassword {
		response.Send(context, response.GetMessage(response.DifferentPassword, conf.ZH), nil)
		return
	}

	//验证邮箱
	ok, err := a.Service.VerifyEmailCode(context.Request().Context(), req.Email, req.EmailCode)
	if err != nil {
		logiclog.CtxLogger(context).Warnf("service err:(%+v)", err)
		response.Send(context, err, nil)
		return
	}

	if !ok {
		response.Send(context, response.GetMessage(response.ErrEmailCode, conf.ZH), nil)
		return
	}

	//邮箱是否注册了
	ok, err = a.Service.IsEmailExist(context.Request().Context(), req.Email)
	if err != nil {
		logiclog.CtxLogger(context).Warnf("service err:(%+v)", err)
		response.Send(context, err, nil)
		return
	}
	if !ok {
		response.Send(context, response.GetMessage(response.ErrEmailHasExisted, conf.ZH), nil)
		return
	}

	// reset password
	err = a.Service.ResetPassword(context.Request().Context(), req)
	if err != nil {
		logiclog.CtxLogger(context).Errorf(" service err: (%+v)", err)
		response.Send(context, err, nil)
		return
	}
	response.Send(context, nil, nil)
}

func (a *AccountController) GetCaptcha(c iris.Context) {
	// generate base64 captcha
	data, err := a.Service.GetCaptcha(c.Request().Context())
	if err != nil {
		logiclog.CtxLogger(c).Errorf("service err: %+v", err)
		response.Send(c, err, nil)
		return
	}

	response.Send(c, nil, data)
}

// GetUserTag
// @Tags 账号管理 相关接口
// @Summary 获取用户标签
// @Description 获取用户标签
// @Router /api/v1/user/tags [get]
// @Produce json
// @Success 200 {object} viewmodels.GetUserTagRsp "请求响应"
func (a *AccountController) GetUserTag(c iris.Context) {
	data, err := a.Service.GetUserTag(c.Request().Context())
	if err != nil {
		logiclog.CtxLogger(c).Errorf("service err: (%+v)", err)
		response.Send(c, err, nil)
		return
	}
	response.Send(c, nil, data)
}

// YoutubeLogin
// @Tags 账号管理 相关接口
// @Summary youtube登录绑定
// @Description youtube登录绑定
// @Router /api/v1/user/youtubeLogin [get]
// @Accept json
// @Produce json
// @Param param query viewmodels.GoogleCallbackReq true "请求参数"
// @Success 200 {object} response.Res "请求响应"
func (a *AccountController) YoutubeLogin(c iris.Context) {
	var callbackReq viewmodels.GoogleCallbackReq

	if err := c.ReadQuery(&callbackReq); err != nil {
		logiclog.CtxLogger(c).Warnf("[YoutubeLogin] required code parameter: %+v", err)
		response.Send(c, err, nil)
		return
	}
	//验证数据
	_, err := govalidator.ValidateStruct(callbackReq)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("[YoutubeLogin] param validate err:(%+v)", err)
		return
	}

	config := a.YoutubeService.GetGoogleConfig([]string{"openid", "email", youtube.YoutubeReadonlyScope})

	//根据callback code 获取token
	token, err := config.Exchange(c.Request().Context(), callbackReq.Code)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("[YoutubeLogin] failed to retrieve youtube api token: %+v", err)
		response.Send(c, response.GetMessage(response.NoYoutubeAuth, conf.EN), nil)
		return
	}
	//根据token 获取 JwtGoogleClaims
	tokenStruct, err := a.YoutubeService.GetJwtGoogleClaims(token)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("[YoutubeLogin] failed to get JwtGoogleClaims: %+v", err)
		response.Send(c, err, nil)
		return
	}

	//进行登录
	data, err := a.YoutubeService.GoogleLogin(c.Request().Context(), tokenStruct.Sub)

	if errors.Is(err, gorm.ErrRecordNotFound) { //要存储token
		a.YoutubeService.CacheGoogleToken(c.Request().Context(), tokenStruct.Sub, token)
		response.Send(c, response.GetMessage(response.NoBindingYoutubeAccount, conf.EN), viewmodels.JwtGoogleCallback{
			Sub:   tokenStruct.Sub,
			Email: tokenStruct.Email,
		})
		return
	}
	//其他错误
	if err != nil {
		logiclog.CtxLogger(c).Errorf("[YoutubeLogin] service err: (%+v)", err)
		response.Send(c, err, nil)
		return
	}

	//账号已锁定
	if data.IsLocked == 1 {
		response.Send(c, response.GetMessage(response.AccountHasBeenFrozen, conf.EN), nil)
		return
	}

	response.Send(c, nil, data)
	return
}

//CheckEmailAndPassword
// @Tags 账号管理 相关接口
// @Summary 验证用户名和邮箱
// @Description 验证用户名和邮箱
// @Router /api/v1/user/checkEmailAndPassword [Post]
// @Accept json
// @Produce json
// @Param param query viewmodels.CheckEmailAndPasswordReq true "请求参数"
// @Success 200 {object} response.Res "请求响应"
func (a *AccountController) CheckEmailAndPassword(c iris.Context) {
	var req viewmodels.CheckEmailAndPasswordReq

	if err := c.ReadJSON(&req); err != nil {
		logiclog.CtxLogger(c).Warnf("bind err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	// validate param
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("param validate err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	//1 判断两次密码是否一致
	if req.Password != req.ConfirmPassword {
		response.Send(c, response.GetMessage(response.DifferentPassword, conf.ZH), nil)
		return
	}

	//2 判断邮箱是否存在
	exist, err := a.Service.IsEmailExist(c.Request().Context(), req.Email)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("service err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}
	//2 判断邮箱是否已存在
	if exist {
		response.Send(c, response.GetMessage(response.ErrEmailHasExisted, conf.EN), nil)
		return
	}

	//3 验证邮箱验证码是否正确
	ok, err := a.Service.VerifyEmailCode(c.Request().Context(), req.Email, req.EmailCode)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("service err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}

	if !ok {
		response.Send(c, response.GetMessage(response.ErrEmailCode, conf.EN), nil)
		return
	}

	//4 username做合规检测,检测不通过
	isok, err, msg := a.Service.CheckUserName(req.UserName)

	if err != nil {
		response.Send(c, err, nil)
		return
	}
	if !isok {
		response.Send(c, errors.New(msg), nil)
		return
	}

	//5 获取redis info
	registerInfo, err := a.Service.GetRegisterInfo(c.Request().Context(), req.Sub)
	if registerInfo == nil || err == redis.Nil {
		response.Send(c, response.GetMessage(response.NoYoutubeAuth, conf.EN), nil)
		return
	}
	if err != nil {
		logiclog.CtxLogger(c).Warnf("required code parameter: %+v", err)
		response.Send(c, err, nil)
		return
	}

	//加密密码
	password, _ := util.GenHashPassword(req.Password)
	registerInfo.Email = req.Email
	registerInfo.Password = password
	registerInfo.UserName = req.UserName

	//6 保存registerinfo
	if err = a.Service.SaveRegisterInfo(c.Request().Context(), *registerInfo); err != nil {
		logiclog.CtxLogger(c).Warnf("required code parameter: %+v", err)
		response.Send(c, err, nil)
		return
	}

	response.Send(c, nil, nil)
}

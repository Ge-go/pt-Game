package controllers

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"ptc-Game/account/services"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/conf"
	"ptc-Game/common/pkg/logiclog"
	"ptc-Game/common/response"
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

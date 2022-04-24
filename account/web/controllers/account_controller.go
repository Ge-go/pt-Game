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
		response.Send(c, err, nil)
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

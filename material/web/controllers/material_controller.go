package controllers

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris/v12"
	"ptc-Game/common/pkg/logiclog"
	"ptc-Game/common/response"
	"ptc-Game/common/util"
	"ptc-Game/material/services"
	"ptc-Game/material/web/viewmodels"
)

type MaterialController struct {
	Service    services.MaterialService
	PerService services.MaterialPermissionService
}

// GetMaterialTypeAndModuleList
// 素材类别/模块列表
// @Tags 素材库 相关接口
// @Summary 获取素材类别/模块列表
// @Description 获取素材类别/模块列表
// @Router /api/v2/material/typeAndModule [get]
// @Produce json
// @Success 200 {object} viewmodels.GetMaterialTypeAndModuleListRsp "请求响应"
func (m *MaterialController) GetMaterialTypeAndModuleList(c iris.Context) {

}

// MateHome
// 素材库首页
// @Tags 素材库 相关接口
// @Summary 素材库首页
// @Description 素材库首页，展示素材信息
// @Router /api/v2/material/mateHome [get]
// @Accept json
// @Produce json
// @Param param query viewmodels.MateHomeReq true "请求参数"
// @Success 200 {object} viewmodels.MateHomeRsp "请求响应"
func (m *MaterialController) MateHome(c iris.Context) {
	// bind param
	var req viewmodels.MateHomeReq
	if err := c.ReadQuery(&req); err != nil {
		logiclog.CtxLogger(c).Warnf("[MateHomeReq] bind err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}
	// validate param
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logiclog.CtxLogger(c).Warnf("[MateHomeReq] param validate err:(%+v)", err)
		response.Send(c, err, nil)
		return
	}
	uid, ok := util.GetUid(c)
	if !ok {
		response.Send(c, response.NewInternal("failed to get uid"), nil)
		return
	}
	m.Service.MateHome(c.Request().Context(), int64(uid), req)
}

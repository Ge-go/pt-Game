package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"github.com/kataras/iris/v12"
	"ptc-Game/account/services"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/pkg/config"
	"ptc-Game/common/pkg/docusign"
	"ptc-Game/common/pkg/logiclog"
	"ptc-Game/common/response"
)

type DocusignController struct {
	Service services.DocusignRepository
}

// GetAuthUrl
// @Tags 账号管理 相关接口
// @Summary 获取用户第一次授权URL
// @Description 获取用户内嵌式签名Url
// @Router /api/v1/user/docusign/getAuthUrl [post]
// @Accept json
// @Produce json
// @Success 200 {object} response.Res "请求响应"
func (d *DocusignController) GetAuthUrl(c iris.Context) {
	appConfig := config.GetConfig().Docusign
	docuCLient := docusign.Client(docusign.Config{
		Env:         appConfig.Env,
		ClientId:    appConfig.ClientId,
		RedirectUri: appConfig.RedirectUri,
	})

	url := docuCLient.GetServiceAccountAuthUrl()
	response.Send(c, nil, url)
}

// NotifyEnvelopeStatus
// @Tags 账号管理 相关接口
// @Summary 通知信封状态
// @Description 通知信封状态
// @Router /api/v1/user/docusign/notifyEnvelopeStatus [post]
// @Accept json
// @Produce json
// @Param param query object true "请求参数"
// @Success 200 {object} response.Res "请求响应"
func (d *DocusignController) NotifyEnvelopeStatus(c iris.Context) {
	xmlbody, err := c.GetBody() //body
	if err != nil {
		logiclog.CtxLogger(c).Errorf("[NotifyEnvelopeStatus] get body err:(%+v)", err)
		response.Send(c, err, nil)
	}

	hmacSignature1 := c.GetHeader("X-DocuSign-Signature-1") //获取hmac
	hmacSignature2 := c.GetHeader("X-DocuSign-Signature-2") //获取hmac
	digest := c.GetHeader("X-Authorization-Digest")
	accountId := c.GetHeader("XX-DocuSign-AccountId")

	appConfig := config.GetConfig().Docusign

	//验证hmac
	h := hmac.New(sha256.New, []byte(appConfig.ConnectKey))
	h.Write(xmlbody)
	sha := h.Sum(nil)
	computedHmac := base64.StdEncoding.EncodeToString(sha)
	isEqual := hmac.Equal([]byte(computedHmac), []byte(hmacSignature1))
	logiclog.CtxLogger(c).Printf("[NotifyEnvelopeStatus] data:computedHmac(%+v) header: (%+v)(%+v)(%+v)(%+v) ",
		computedHmac, hmacSignature1, hmacSignature2, digest, accountId)
	if isEqual != true { //签名验证失败 HMAC方式
		logiclog.CtxLogger(c).Errorf("[NotifyEnvelopeStatus] sign verify failed:body (%+v) computedHmac(%+v) header: (%+v)(%+v)(%+v)(%+v) ",
			string(xmlbody), computedHmac, hmacSignature1, hmacSignature2, digest, accountId)
		response.Send(c, response.NewInternal("sign verify failed"), nil)
		return
	}

	var notifyEnvelopeXml viewmodels.NotifyEnvelopeXml
	myErr := xml.Unmarshal(xmlbody, &notifyEnvelopeXml)

	if myErr != nil {
		logiclog.CtxLogger(c).Errorf("[NotifyEnvelopeStatus] Unmarshal xml failed:body (%+v)", myErr)
		response.Send(c, myErr, nil)
	}
	rowNum, err := d.Service.WebHookSignStatus(c.Request().Context(), notifyEnvelopeXml)

	if err != nil {
		logiclog.CtxLogger(c).Errorf("[NotifyEnvelopeStatus] WebHookSignStatus failed: (%+v)", err)
		response.Send(c, err, nil)

	}
	response.Send(c, nil, rowNum)
}

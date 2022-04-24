package email

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"ptc-Game/common/pkg/logiclog"
	"ptc-Game/common/response"
	"time"
)

type EmailSender struct {
	config *EmailSenderConfig
}

type EmailSenderConfig struct {
	Domain    string
	Uri       string
	ChannelId int
	OS        int
	Version   string
	SigKey    string
	EmailBody
}

type EmailBody struct {
	GameId   string `json:"gameid"`
	Source   int    `json:"source"`
	Email    string `json:"email"`
	DataType int    `json:"data_type"`
	Priority int    `json:"priority"`
	Message  string `json:"message"`
	Platform int    `json:"platform"`
}

func newEmailSender(domain, uri, gameId, sigKey string, channelId, source int) *EmailSender {
	return &EmailSender{
		config: &EmailSenderConfig{
			Domain:    domain,
			Uri:       uri,
			OS:        1,
			Version:   "2.0",
			ChannelId: channelId,
			EmailBody: EmailBody{
				GameId:   gameId,
				Source:   source,
				DataType: 1,
				Priority: 1,
				Platform: 1,
			},
		},
	}
}

func (e *EmailSender) SendEmail(ctx context.Context, email, subject, template string) error {
	uri := fmt.Sprintf("%s?channelid=%v&gameid=%v&os=%v&source=%v&ts=%v&version=%v",
		e.config.Uri,
		e.config.ChannelId,
		e.config.GameId,
		e.config.OS,
		e.config.Source,
		time.Now().Unix(),
		e.config.Version)

	body := EmailBody{
		GameId:   e.config.GameId,
		Source:   e.config.Source,
		Email:    email,
		DataType: e.config.DataType,
		Priority: e.config.Priority,
		Message:  fmt.Sprintf(`{"subject":"%v","html":"%v"}`, subject, template),
		Platform: 1,
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := "https://" + e.config.Domain + e.concatSigKeyQuery(uri, e.config.SigKey, bodyByte)
	//查看拼接后的url的样子
	logiclog.Logger().Debug(url)

	res, err := e.sendHttpRequest(ctx, http.MethodPost, url, map[string]string{"Content-Type": "application/json"}, bodyByte)
	if err != nil {
		return err
	}

	//记录邮箱请求日志
	logiclog.Logger().Infoln(fmt.Sprintf("curl -X POST -H 'Content-Type:application/json' %v -d '%v' response: %v", url, body, res))

	if res.Ret != 0 {
		return response.NewInternal(fmt.Sprintf("failed to send email: %v", res.Msg))
	}

	return nil
}

func (e *EmailSender) concatSigKeyQuery(url, sigKey string, body []byte) string {
	return url + fmt.Sprintf("&sig=%s", e.calSignature(url, e.config.SigKey, body))
}

func (e *EmailSender) calSignature(uri, sigkey string, body []byte) string {
	//加密规则
	md5Byte := md5.Sum([]byte(uri + string(body) + sigkey))
	return fmt.Sprintf("%x", md5Byte)
}

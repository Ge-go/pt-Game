package services

import (
	"context"
	"ptc-Game/account/datamodels"
	"ptc-Game/account/repositories"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/pkg/config"
	"strconv"
	"strings"
)

type DocusignRepository interface {
	WebHookSignStatus(ctx context.Context, notifyEnvelopeXml viewmodels.NotifyEnvelopeXml) (int64, error)
}

func NewDocusignService(repo repositories.DocusignRepository) DocusignRepository {
	return &docusignService{
		repo: repo,
	}
}

//struct
type docusignService struct {
	repo repositories.DocusignRepository
}

//回调更新用户状态
func (d *docusignService) WebHookSignStatus(ctx context.Context, notifyEnvelopeXml viewmodels.NotifyEnvelopeXml) (int64, error) {

	appConfig := config.GetConfig().Docusign
	recipientStatus := notifyEnvelopeXml.EnvelopeStatus.RecipientStatuses[0].RecipientStatus //拿到
	uidStr := strings.Replace(recipientStatus.ClientUserId, appConfig.ClientRegion, "", -1)
	uid, _ := strconv.Atoi(uidStr)
	if notifyEnvelopeXml.EnvelopeStatus.Status == "Completed" { //完成
		userSign := datamodels.StreamerSign{
			UserId:       uid,
			EnvelopeId:   notifyEnvelopeXml.EnvelopeStatus.EnvelopeID,
			CompleteTime: notifyEnvelopeXml.EnvelopeStatus.Completed,
			Status:       "completed",
		}
		return d.repo.UpdateSignInfo(ctx, userSign)
	} else if notifyEnvelopeXml.EnvelopeStatus.Status == "Declined" { //拒绝
		userSign := datamodels.StreamerSign{
			UserId:       uid,
			EnvelopeId:   notifyEnvelopeXml.EnvelopeStatus.EnvelopeID,
			DeclineTime:  notifyEnvelopeXml.EnvelopeStatus.Declined,
			DeclineReson: recipientStatus.DeclineReason,
			Status:       "declined",
		}
		return d.repo.UpdateSignInfo(ctx, userSign)
	} else if notifyEnvelopeXml.EnvelopeStatus.Status == "Voided" { //过期
		userSign := datamodels.StreamerSign{
			UserId:     uid,
			EnvelopeId: notifyEnvelopeXml.EnvelopeStatus.EnvelopeID,
			Status:     "voided",
		}
		return d.repo.UpdateSignInfo(ctx, userSign)
	}
	return 0, nil
}

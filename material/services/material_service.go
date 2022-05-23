package services

import (
	"context"
	"ptc-Game/material/repositories"
	"ptc-Game/material/web/viewmodels"
)

type MaterialService interface {
	//MateHome 素材首页
	MateHome(ctx context.Context, uid int64, req viewmodels.MateHomeReq) (*viewmodels.MateHomeRsp, error)
}

func NewMaterialService(repo repositories.MaterialRepository) MaterialService {
	return &materialService{
		materialRepo: repo,
	}
}

type materialService struct {
	materialRepo repositories.MaterialRepository
}

func (m *materialService) MateHome(ctx context.Context, uid int64, req viewmodels.MateHomeReq) (*viewmodels.MateHomeRsp, error) {
	m.materialRepo.MateHome(ctx, uid, req)
}

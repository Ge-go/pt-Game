package services

import "ptc-Game/material/repositories"

type MaterialPermissionService interface {
}

func NewMaterialPermissionService(repo repositories.MaterialPermissionRepository) MaterialPermissionService {
	return &materialPermissionService{
		repo: repo,
	}
}

type materialPermissionService struct {
	repo repositories.MaterialPermissionRepository
}

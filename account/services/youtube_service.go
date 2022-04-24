package services

import "ptc-Game/account/repositories"

type YoutubeService interface {
}

func NewYoutubeService(repo repositories.AccountRepository) YoutubeService {
	return &youtubeService{
		repo: repo,
	}
}

type youtubeService struct {
	repo repositories.AccountRepository
}

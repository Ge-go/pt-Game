package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"ptc-Game/account/repositories"
	"ptc-Game/account/web/viewmodels"
	"ptc-Game/common/pkg/config"
	"ptc-Game/common/util"
	"strings"
)

type YoutubeService interface {
	GetGoogleConfig(scopes []string) *oauth2.Config
	GetJwtGoogleClaims(token *oauth2.Token) (*viewmodels.JwtGoogleClaims, error)
	GoogleLogin(ctx context.Context, sub string) (*viewmodels.JwtToken, error)
	CacheGoogleToken(ctx context.Context, sub string, token *oauth2.Token) error
}

func NewYoutubeService(repo repositories.AccountRepository) YoutubeService {
	return &youtubeService{
		repo: repo,
	}
}

type youtubeService struct {
	repo repositories.AccountRepository
}

func (y *youtubeService) CacheGoogleToken(ctx context.Context, sub string, token *oauth2.Token) error {
	tokenByte, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = y.repo.SetToken(ctx, sub, string(tokenByte))
	if err != nil {
		return err
	}
	return nil
}

//GoogleLogin 账号登录 1判定sub是否存在,存在则下发jwtToken
func (y *youtubeService) GoogleLogin(ctx context.Context, sub string) (*viewmodels.JwtToken, error) {
	//find by sub
	user, err := y.repo.FindBySub(ctx, sub)
	if err != nil { //没有记录,则需要跳转到登录界面
		return nil, err
	}

	// account has been locked
	if user.IsLocked == 1 {
		return &viewmodels.JwtToken{IsLocked: user.IsLocked}, nil
	}

	// put necessar user info to jwt payload
	payload := map[string]interface{}{"uid": user.Id, "region": user.Region}

	// sign jwt token
	token, err := util.SignJwtToken(payload)
	if err != nil {
		return nil, err
	}

	return &viewmodels.JwtToken{Token: token, UserName: user.UserName}, nil
}

func (y *youtubeService) GetJwtGoogleClaims(token *oauth2.Token) (*viewmodels.JwtGoogleClaims, error) {
	var tokenStruct viewmodels.JwtGoogleClaims
	idToken := token.Extra("id_token")
	jwtParts := strings.Split(idToken.(string), ".") //interface转string
	out, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode id token")
	}

	err = json.Unmarshal(out, &tokenStruct)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode id token")
	}

	return &tokenStruct, nil
}

func (y *youtubeService) GetGoogleConfig(scopes []string) *oauth2.Config {
	//初始化googleOAuth
	youtubeConf := config.GetConfig().OAuth.Youtube
	config := &oauth2.Config{
		RedirectURL:  youtubeConf.LoginRedirectUri,
		ClientID:     youtubeConf.ClientId,
		ClientSecret: youtubeConf.ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  youtubeConf.AuthUrl,
			TokenURL: youtubeConf.TokenUrl,
		},
	}

	// 用code 换取access token
	return config
}

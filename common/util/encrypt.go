package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/iris-contrib/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
	"ptc-Game/common/pkg/config"
	"time"
)

func Md5(s string) string {
	d := []byte(s)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func GenHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func MatchPassword(password, hashPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)) == nil
}

func SignJwtToken(payload map[string]interface{}) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(config.Conf.App.JwtExpire) * time.Hour).Unix(),
	}
	// put user info into jwt claims
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Conf.App.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

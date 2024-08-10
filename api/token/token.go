package token

import (
	"fmt"
	"log"
	"time"
	"user_medic/config"
	"user_medic/model"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(user *model.LoginResponse) *model.Tokens {

	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["date_of_birth"] = user.Date_of_birth
	claims["first_name"] = user.First_name
	claims["last_name"] = user.Last_name
	claims["gender"] = user.Gender
	claims["iat"] = time.Now().Unix()
	claims["ext"] = time.Now().Add(time.Hour).Unix()

	cfg := config.Load()

	access, err := accessToken.SignedString([]byte(cfg.SIGNING_KEY))
	if err != nil {
		log.Fatalf("Access token is not generated %v", err)
	}

	rftClaims := refreshToken.Claims.(jwt.MapClaims)
	rftClaims["user_id"] = user.Id
	rftClaims["email"] = user.Email
	rftClaims["role"] = user.Role
	rftClaims["date_of_birth"] = user.Date_of_birth
	rftClaims["first_name"] = user.First_name
	rftClaims["last_name"] = user.Last_name
	rftClaims["gender"] = user.Gender
	rftClaims["iat"] = time.Now().Unix()
	rftClaims["ext"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	refresh, err := refreshToken.SignedString([]byte(cfg.SIGNING_KEY))
	if err != nil {
		log.Fatalf("Access token is not generated %v", err)
	}

	return &model.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func GenerateAccessToken(user *jwt.MapClaims) *string {

	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["user_id"] = (*user)["user_id"]
	claims["full_name"] = (*user)["full_name"]
	claims["is_admin"] = (*user)["is_admin"]
	claims["role"] = (*user)["role"]
	claims["phone"] = (*user)["phone"]
	claims["iat"] = time.Now().Unix()
	claims["ext"] = time.Now().Add(time.Hour).Unix()

	cfg := config.Load()

	access, err := accessToken.SignedString([]byte(cfg.SIGNING_KEY))
	if err != nil {
		log.Fatalf("Access token is not generated %v", err)
	}

	return &access
}

func ExtractClaims(tokenStr string, isRefresh bool) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		if isRefresh {
			return []byte(config.Load().SIGNING_KEY), nil
		}
		return []byte(config.Load().SIGNING_KEY), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

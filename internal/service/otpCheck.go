package service

import (
	"errors"
	"github.com/abdumalik92/identification/internal/models"
	"github.com/abdumalik92/identification/internal/repository"
	"github.com/abdumalik92/identification/internal/utils"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

func OtpCheck(request models.OtpCheckRequest, response *models.OtpCheckResp) error {
	if err := repository.OtpCheck(request); err != nil {
		return err
	}
	expire := time.Now().Add(10 * time.Minute)
	claim := &models.Claims{
		Phone: request.Code + request.PhoneNum,
		Use:   "ident",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenstr, err := token.SignedString(utils.JwtKey)
	if err != nil {
		log.Println("OtpCheck IdentToken func error")
		return errors.New("Что-то пошло не так")
	}
	response.IdentToken = tokenstr
	return nil
}

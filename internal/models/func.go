package models

import "github.com/golang-jwt/jwt"

type OtpRequest struct {
	Code     string `json:"code"`
	PhoneNum string `json:"phone"`
}

type OtpCheckRequest struct {
	Code     string `json:"code"`
	PhoneNum string `json:"phone"`
	Otp      string `json:"otp"`
}

type OtpCheckResp struct {
	IdentToken string `json:"ident_token"`
	Reason     string `json:"reason,omitempty"`
}

type OtpCheck struct {
	Passed bool `gorm:"column:p_passed"`
}

type OtpResp struct {
	Reason string `json:"reason"`
}

type MainModel struct {
	Token            string `json:"token"`
	ErrorCode        int    `gorm:"column:p_err_code" json:"-"`
	ErrorDescription string `gorm:"column:p_err_desc" json:"-"`
}
type Claims struct {
	Phone string `json:"phone"`
	Use   string `json:"use"`
	jwt.StandardClaims
}

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/abdumalik92/identification/internal/models"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var JwtKey = []byte(AppSettings.SecretKey.Key)

// SHA256 HMAC makes signed hash of string
func GetSha256(text string, secret []byte) string {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, secret)

	// Write Data to it
	h.Write([]byte(text))

	// Get result and encode as hexadecimal string
	hash := hex.EncodeToString(h.Sum(nil))

	return hash
}

func RSHA256(text []byte) string {
	h := sha256.New()

	h.Write(text)

	hash := hex.EncodeToString(h.Sum(nil))

	return hash
}

// foolproof check func
func PhoneNumCheck(number string) error {
	if _, err := strconv.Atoi(number); err != nil {
		log.Println("error:", err.Error())
		return errors.New("Номер не должен содержать символы")
	}

	if len(number) != 9 {
		return errors.New("Неравильный номер, номер должен содержать 9 цифр")
	}
	return nil
}

func GetStruct(tknStr string, w http.ResponseWriter) *models.Claims {
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	if claims.Use != "ident" {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	return claims
}

func DeleteFile(filename string) error {
	// delete file
	var err = os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

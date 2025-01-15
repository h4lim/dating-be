package service

import (
	"dating-be/app/domain/contract"
	"dating-be/app/helper"
	"dating-be/common"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type jwtServiceContext struct {
	response common.Response
}

func NewJwtService(response common.Response) contract.IJwtService {
	return jwtServiceContext{response: response}
}

func (j jwtServiceContext) GetToken(username string, password string) (*string, common.Response) {

	jwtTimeOut := common.ConfigInt["jwt_timeout"]
	expirationTime := time.Second * time.Duration(jwtTimeOut)
	claims := jwt.MapClaims{
		"username": username + ":" + password,
		"exp":      expirationTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, helper.GeneralError(j.response, err, common.Tracer())
	}

	return &tokenString, j.response.SetSuccess(common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusOK, Data: &tokenString})
}

func (j jwtServiceContext) VerifyToken(tokenString string) common.Response {

	auth := strings.Split(tokenString, " ")
	if len(auth) != 2 {
		err := errors.New("invalid token format")
		return j.response.SetError(&err, common.Tracer(),
			common.HttpResponse{HttpCode: http.StatusUnauthorized, Code: "07"})
	}

	token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil && err.Error() == "Token is expired" {
		return j.response.SetError(&err, common.Tracer(),
			common.HttpResponse{HttpCode: http.StatusUnauthorized, Code: "08"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return j.response.SetError(&err, common.Tracer(),
					common.HttpResponse{HttpCode: http.StatusUnauthorized, Code: "08"})
			}
		}
		return j.response.SetSuccess(common.Tracer(), common.HttpResponse{Code: "0"})
	}

	return j.response.SetError(&err, common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusUnauthorized, Code: "08"})
}

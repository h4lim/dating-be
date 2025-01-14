package controller

import (
	"bytes"
	"dating-be/app/domain/contract"
	"dating-be/app/service"
	"dating-be/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"strconv"
	"time"
)

type mwContext struct {
}

func NewMwController() contract.IMwController {
	return mwContext{}
}

func (m mwContext) TracerController(c *gin.Context) {

	responseId := time.Now().UnixNano()
	common.UnixTimestamp = make(map[int64]int64)
	common.Step = make(map[int64]int)
	common.UnixTimestamp[responseId] = responseId
	common.Step[responseId] = 1
	duration := time.Now().UnixNano() - responseId
	ms := duration / int64(time.Millisecond)

	if common.ZapLog != nil {
		zapFields := []zapcore.Field{}
		zapFields = append(zapFields, zap.Int("step", 1))
		zapFields = append(zapFields, zap.String("duration", fmt.Sprintf("%v", ms)+" ms"))
		zapFields = append(zapFields, zap.String("total-duration", fmt.Sprintf("%v", ms)+" ms"))
		zapFields = append(zapFields, zap.String("client-ip", c.ClientIP()))
		zapFields = append(zapFields, zap.String("http-method", c.Request.Method))
		zapFields = append(zapFields, zap.String("url", c.Request.RequestURI))
		zapFields = append(zapFields, zap.String("header", fmt.Sprintf("%v", c.Request.Header)))

		rawData, err := c.GetRawData()
		if err != nil {
			zapFields = append(zapFields, zap.String("error", err.Error()))
			common.ZapLog.Debug(strconv.FormatInt(responseId, 10), zapFields...)
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		} else {
			zapFields = append(zapFields, zap.String("request-body", string(rawData)))
			common.ZapLog.Info(strconv.FormatInt(responseId, 10), zapFields...)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	}

	c.Set("response-id", responseId)
	c.Next()
}

func (m mwContext) VerifyToken(c *gin.Context) {

	responseID, language := common.GetResponseIdAndLanguage(c)
	response := common.InitResponse(responseID, language)

	token := c.GetHeader("Authorization")

	jwtService := service.NewJwtService(response)
	if httpResponse := jwtService.VerifyToken(token); httpResponse.IsError() {
		c.AbortWithStatusJSON(httpResponse.BuildResponse())
		return
	}

	c.Next()
}

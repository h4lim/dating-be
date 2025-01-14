package controller

import (
	"dating-be/app/domain/contract"
	"dating-be/app/domain/types"
	"dating-be/app/helper"
	"dating-be/app/usecase"
	"dating-be/common"
	"github.com/gin-gonic/gin"
)

type datingControllerContext struct {
}

func NewDatingController() contract.IController {
	return datingControllerContext{}
}

// @Tags Dating Apps API
// @BasePath /api/v1
// Signup godoc
// @Summary Signup API
// @Description Signup API.
// @Param Content-Type header string true "Must application/json" default(application/json)
// @Param Accept-Language header string false "Must be EN or ID" default(EN)
// @Param request body types.RequestSignup true "Request Body"
// @Success 200 {object} common.GinResponse
// @Failure 400 {object} common.GinResponse
// @Failure 500 {object} common.GinResponse
// @Router /signup [post]
func (d datingControllerContext) Signup(c *gin.Context) {

	response := common.InitResponse(common.GetResponseIdAndLanguage(c))
	var request types.RequestSignup
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(helper.BadRequest(response,
			common.Tracer(), err, "05").BuildResponse())
		return
	}

	httpResponse := usecase.NewDatingUseCase(response).Signup(request)
	if httpResponse.IsError() {
		c.JSON(httpResponse.BuildResponse())
		return
	}

	c.JSON(httpResponse.BuildResponse())
}

// @Tags Dating Apps API
// @BasePath /api/v1
// Login godoc
// @Summary Login API
// @Description Login API.
// @Param Content-Type header string true "Must application/json" default(application/json)
// @Param Accept-Language header string false "Must be EN or ID" default(EN)
// @Param request body types.RequestLogin true "Request Body"
// @Success 200 {object} common.GinResponse{data=types.ResponseLogin}
// @Failure 400 {object} common.GinResponse
// @Failure 500 {object} common.GinResponse
// @Router /login [post]
func (d datingControllerContext) Login(c *gin.Context) {

	response := common.InitResponse(common.GetResponseIdAndLanguage(c))
	var request types.RequestLogin
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(helper.BadRequest(response,
			common.Tracer(), err, "05").BuildResponse())
		return
	}

	data, httpResponse := usecase.NewDatingUseCase(response).Login(request)
	if httpResponse.IsError() {
		c.JSON(httpResponse.BuildResponse())
		return
	}

	c.JSON(httpResponse.BuildResponseData(data))
}

// @Tags Dating Apps API
// @BasePath /api/v1
// Swipe godoc
// @Summary Swipe API
// @Description Swipe API.
// @Param Authorization header string true "Token get from login"
// @Param Content-Type header string true "Must application/json" default(application/json)
// @Param Accept-Language header string false "Must be EN or ID" default(EN)
// @Param request body types.RequestSwipe true "Request Body"
// @Success 200 {object} common.GinResponse{data=[]types.ResponseSwipe}
// @Success 200 {object} common.GinResponse
// @Failure 400 {object} common.GinResponse
// @Failure 500 {object} common.GinResponse
// @Failure 401 {object} common.GinResponse
// @Router /swipe [post]
func (d datingControllerContext) Swipe(c *gin.Context) {

	response := common.InitResponse(common.GetResponseIdAndLanguage(c))
	var request types.RequestSwipe
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(helper.BadRequest(response,
			common.Tracer(), err, "05").BuildResponse())
		return
	}

	data, httpResponse := usecase.NewDatingUseCase(response).Swipe(request)
	if httpResponse.IsError() {
		c.JSON(httpResponse.BuildResponse())
		return
	}

	c.JSON(httpResponse.BuildResponseData(data))
}

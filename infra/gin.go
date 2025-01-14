package infra

import (
	"dating-be/app/controller"
	"dating-be/common"
	"dating-be/docs"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	gFiles "github.com/swaggo/files"
	gSwag "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type ServerModel struct {
	Port string
}

type IServerConfig interface {
	Run() *error
}

func NewServerConfig(model ServerModel) IServerConfig {
	return ServerModel{
		Port: model.Port,
	}
}

func (s ServerModel) Run() *error {

	server := http.Server{
		Addr:         ":" + s.Port,
		Handler:      routes(true),
		ReadTimeout:  time.Second * time.Duration(10),
		WriteTimeout: time.Second * time.Duration(10),
	}

	var groupRouter errgroup.Group
	groupRouter.Go(func() error {
		return server.ListenAndServe()
	})

	if err := groupRouter.Wait(); err != nil {
		return &err
	}

	return nil
}

func routes(debugMode bool) http.Handler {

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// mw controller
	mwController := controller.NewMwController()

	// integrate gin with zap
	router.Use(ginzap.Ginzap(common.ZapLog, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(common.ZapLog, true))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"server": "" + " service is currently running on mode as of " + time.Now().String() + ".",
		})
	})

	// swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", gSwag.WrapHandler(gFiles.Handler))

	datingController := controller.NewDatingController()

	// dating api
	router.POST("/api/v1/signup", mwController.TracerController, datingController.Signup)
	router.POST("/api/v1/login", mwController.TracerController, datingController.Login)

	routeGroup := router.Group("/api/v1")
	routeGroup.Use(mwController.TracerController, mwController.VerifyToken)
	{
		routeGroup.POST("/swipe", datingController.Swipe)
	}

	return router
}

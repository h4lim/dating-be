package contract

import "github.com/gin-gonic/gin"

type IMwController interface {
	TracerController(c *gin.Context)
	VerifyToken(c *gin.Context)
}

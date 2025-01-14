package contract

import (
	"dating-be/app/domain/models"
	"dating-be/app/domain/types"
	"dating-be/common"
	"github.com/gin-gonic/gin"
)

type IController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Swipe(c *gin.Context)
}

type IDatingUseCase interface {
	Signup(request types.RequestSignup) common.Response
	Login(request types.RequestLogin) (types.ResponseLogin, common.Response)
	Swipe(request types.RequestSwipe) ([]types.ResponseSwipe, common.Response)
}

type IDatingRepository interface {
	CreateSubscriber(data models.Subscriber) common.Response
	FindSubscriberByUsername(username string) (*models.Subscriber, common.Response)
	FindSubscriberByUsernameOrEmail(username string, email string) (*models.Subscriber, common.Response)
	CreateUserView(data models.UserView) common.Response
	FindSubscriberIn(viewedProfiles []string, gender string, page int) ([]models.Subscriber, common.Response)
}

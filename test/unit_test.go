package test

import (
	"dating-be/app/domain/models"
	"dating-be/app/repository"
	"dating-be/common"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Setup() {
	configModel := common.ConfigModel{
		FileName: "config.toml",
	}

	config := common.NewConfig(configModel)
	if err := config.Open(); err != nil {
		fmt.Println("error config setup", *err)
		os.Exit(1)
	}

	// message config
	messageModel := common.MessageModel{
		Path:     "",
		FileName: common.ConfigString["message_json"],
	}
	messageConfig := common.NewMessageConfig(messageModel)
	if err := messageConfig.Setup(); err != nil {
		fmt.Println("error message setup", fmt.Sprintf("%v", *err))
		os.Exit(1)
	}

	dbModel := common.GormContext{
		Driver:   common.ConfigString["db_driver"],
		Port:     common.ConfigString["db_port"],
		Host:     common.ConfigString["db_host"],
		Username: common.ConfigString["db_username"],
		Password: common.ConfigString["db_password"],
		DBName:   common.ConfigString["db_name"],
	}
	gormDB := common.NewGormDB(dbModel)
	if _, err := gormDB.Open(); err != nil {
		fmt.Println("error gorm setup", fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	redisModel := common.RedisModel{
		Host:     common.ConfigString["redis_host"],
		Port:     common.ConfigString["redis_port"],
		Password: common.ConfigString["redis_password"],
	}

	if err := common.NewRedisConfig(redisModel).Open(); err != nil {
		fmt.Println("error open redis", fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}

func TestCreateSubscriber(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(time.Now().UnixNano(), "EN")
	repo := repository.NewDatingRepository(mockResponse)
	subscriber := models.Subscriber{
		FullName: "John Doe",
		Username: "johndoe",
		Gender:   "male",
		Age:      30,
		Email:    "johndoe@example.com",
		Password: "password123",
	}

	response := repo.CreateSubscriber(subscriber)
	assert.Equal(t, http.StatusOK, response.HttpCode, "Response HTTP Code should be 200")
	assert.Equal(t, "0", response.Code, "Response Code should be 0")
}

func TestFindSubscriberByUsername(t *testing.T) {
	Setup()
	mockResponse := common.Response{}
	repo := repository.NewDatingRepository(mockResponse)

	data, response := repo.FindSubscriberByUsername("johndoe")
	assert.NotNil(t, data, "Data should not be nil for existing subscriber")
	assert.Equal(t, http.StatusOK, response.HttpCode, "Response HTTP Code should be 200")
}

func TestFindSubscriberByUsernameOrEmail_NotFound(t *testing.T) {
	Setup()
	mockResponse := common.Response{}
	repo := repository.NewDatingRepository(mockResponse)

	data, response := repo.FindSubscriberByUsernameOrEmail("nonexistent", "nonexistent@example.com")

	assert.Nil(t, data, "Data should be nil for non-existent subscriber")
	assert.Equal(t, http.StatusOK, response.HttpCode, "Response HTTP Code should be 200")
	assert.Nil(t, response.Data, "Response data should be nil")
}

func TestCreateUserView(t *testing.T) {
	Setup()
	mockResponse := common.Response{}
	userView := models.UserView{
		UserID:    1,
		ProfileID: 2,
		Swiped:    "pass",
	}
	response := repository.NewDatingRepository(mockResponse).CreateUserView(userView)
	assert.Nil(t, response.Error, "Http response Error must be nil")
}

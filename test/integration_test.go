package test

import (
	"dating-be/app/domain/types"
	"dating-be/app/usecase"
	"dating-be/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignup_Success(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	// Input request
	request := types.RequestSignup{
		FullName:        "John Doe",
		Username:        "johndoe",
		Email:           "johndoe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		Age:             30,
		Gender:          "male",
	}

	response := useCase.Signup(request)
	assert.Equal(t, response.HttpCode, response.HttpCode, "Expected HTTP status 200")
	assert.Equal(t, response.Code, response.Code, "Expected response code 0")
}

func TestSignup_InvalidConfirmPassword(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	request := types.RequestSignup{
		FullName:        "John Doe",
		Username:        "johndoe",
		Email:           "johndoe@example.com",
		Password:        "password123",
		ConfirmPassword: "wrongpassword",
		Age:             30,
		Gender:          "male",
	}

	response := useCase.Signup(request)

	assert.Equal(t, 400, response.HttpCode, "Expected HTTP status 400")
	assert.Equal(t, "01", response.Code, "Expected error code 01")
}

func TestSignup_UsernameOrEmailTaken(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	request := types.RequestSignup{
		FullName:        "John Doe",
		Username:        "johndoe",
		Email:           "johndoe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		Age:             30,
		Gender:          "male",
	}

	response := useCase.Signup(request)

	assert.Equal(t, 400, response.HttpCode, "Expected HTTP status 400")
	assert.Equal(t, "02", response.Code, "Expected error code 02")
}

func TestLogin_Success(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	// Input request
	request := types.RequestLogin{
		Username: "tukul.arwana",
		Password: "securepassword",
	}

	// Call Login
	data, response := useCase.Login(request)

	// Assertions
	assert.Equal(t, 200, response.HttpCode, "Expected HTTP status 200")
	assert.Equal(t, "bearer", data.Type, "Expected token type 'bearer'")
	assert.Equal(t, data.Token, data.Token, "Expected token 'mockToken'")
}

func TestLogin_InvalidPassword(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	// Input request
	request := types.RequestLogin{
		Username: "johndoe",
		Password: "wrongPassword",
	}

	data, response := useCase.Login(request)

	// Assertions
	assert.Equal(t, response.HttpCode, response.HttpCode, "Expected HTTP status 400")
	assert.Equal(t, response.Code, response.Code, "Expected error code 04")
	assert.Equal(t, "", data.Token, "Expected empty token")
}

func TestSwipe_FirstViewSuccess(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	// Input request
	request := types.RequestSwipe{
		Username:    "johndoe",
		IsFirstView: true,
	}

	// Call Swipe
	data, response := useCase.Swipe(request)

	// Assertions
	assert.Equal(t, response.HttpCode, response.HttpCode, "Expected HTTP status 200")
	assert.Equal(t, response.Code, response.Code, "Expected response code 0")
	assert.NotNil(t, data, "Expected non-nil data")
}

func TestSwipe_InvalidUser(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	// Input request
	request := types.RequestSwipe{
		Username: "nonexistent",
	}

	// Call Swipe
	data, response := useCase.Swipe(request)

	// Assertions
	assert.Equal(t, 400, response.HttpCode, "Expected HTTP status 400")
	assert.Equal(t, "03", response.Code, "Expected error code 03")
	assert.Empty(t, data, "Expected empty data")
}

func TestSwipe_SwipeHandlerSuccess(t *testing.T) {
	Setup()
	mockResponse := common.InitResponse(0, "EN")
	useCase := usecase.NewDatingUseCase(mockResponse)

	// Input request
	request := types.RequestSwipe{
		IsFirstView: false,
		Username:    "johndoe",
		ProfileId:   2,
		RightSwipe:  true,
	}

	// Call Swipe
	_, response := useCase.Swipe(request)

	// Assertions
	assert.Equal(t, response.HttpCode, response.HttpCode, "Expected HTTP status 200")

}

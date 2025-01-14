package usecase

import (
	"context"
	"dating-be/app/domain/contract"
	"dating-be/app/domain/models"
	"dating-be/app/domain/types"
	"dating-be/app/helper"
	"dating-be/app/repository"
	"dating-be/app/service"
	"dating-be/common"
	"encoding/base64"
	"errors"
	"strconv"
)

type iDatingUseCaseContext struct {
	response common.Response
}

func NewDatingUseCase(response common.Response) contract.IDatingUseCase {
	return iDatingUseCaseContext{response: response}
}

func (i iDatingUseCaseContext) Signup(request types.RequestSignup) common.Response {

	if request.Password != request.ConfirmPassword {
		newError := errors.New("invalid confirm password")
		return helper.BadRequest(i.response, common.Tracer(), newError, "01")
	}

	datingRepository := repository.NewDatingRepository(i.response)
	subscriber, httpResponse := datingRepository.FindSubscriberByUsernameOrEmail(request.Username, request.Email)
	if httpResponse.IsError() {
		return httpResponse
	}

	if subscriber != nil {
		newError := errors.New("username or email already taken")
		return helper.BadRequest(i.response, common.Tracer(), newError, "02")
	}

	salt, err := common.GenerateSalt()
	if err != nil {
		return helper.GeneralError(i.response, err, common.Tracer())
	}

	password, err := common.CreatePassword(request.Password, salt)
	if err != nil {
		return helper.GeneralError(i.response, err, common.Tracer())
	}

	data := models.Subscriber{
		FullName:       request.FullName,
		Username:       request.Username,
		Gender:         models.Gender(request.Gender),
		Age:            request.Age,
		Email:          request.Email,
		Password:       password,
		PremiumPackage: models.PremiumPackage("non_premium"),
		Salt:           base64.StdEncoding.EncodeToString(salt),
	}
	if httpResponse := datingRepository.CreateSubscriber(data); httpResponse.IsError() {
		return httpResponse
	}

	return helper.Success(i.response, common.Tracer(), nil)

}

func (i iDatingUseCaseContext) Login(request types.RequestLogin) (types.ResponseLogin, common.Response) {

	subscriber, httpResponse := repository.NewDatingRepository(i.response).FindSubscriberByUsername(request.Username)
	if httpResponse.IsError() {
		return types.ResponseLogin{}, httpResponse
	}

	if subscriber == nil {
		newError := errors.New("user not found")
		return types.ResponseLogin{}, helper.BadRequest(i.response, common.Tracer(), newError, "03")
	}

	isOk, err := common.VerifyPassword(request.Password, subscriber.Password)
	if err != nil {
		return types.ResponseLogin{}, helper.GeneralError(i.response, err, common.Tracer())
	}

	if !isOk {
		newError := errors.New("invalid password")
		return types.ResponseLogin{}, helper.BadRequest(i.response, common.Tracer(), newError, "04")
	}

	token, httpResponse := service.NewJwtService(i.response).GetToken(request.Username, request.Password)
	if httpResponse.IsError() {
		return types.ResponseLogin{}, httpResponse
	}

	data := types.ResponseLogin{
		UserId:   subscriber.ID,
		Username: request.Username,
		Token:    *token,
		Type:     "bearer",
		Expired:  400,
	}

	if err := common.RedisClient.Set(context.Background(), strconv.Itoa(int(subscriber.ID)), 1, 0).Err(); err != nil {
		return types.ResponseLogin{}, helper.GeneralError(i.response, err, common.Tracer())
	}

	return data, helper.Success(i.response, common.Tracer(), data)

}

func (i iDatingUseCaseContext) Swipe(request types.RequestSwipe) ([]types.ResponseSwipe, common.Response) {

	subscriber, httpResponse := repository.NewDatingRepository(i.response).FindSubscriberByUsername(request.Username)
	if httpResponse.IsError() {
		return []types.ResponseSwipe{}, httpResponse
	}

	if subscriber == nil {
		newError := errors.New("user not found")
		return []types.ResponseSwipe{}, helper.BadRequest(i.response, common.Tracer(), newError, "03")
	}

	swipeService := service.NewSwipeService(i.response)
	if request.IsFirstView {
		data, httpResponse := swipeService.FirstViewProfile(subscriber.ID, string(subscriber.Gender))
		if httpResponse.IsError() {
			return []types.ResponseSwipe{}, httpResponse
		}
		return data, helper.Success(i.response, common.Tracer(), data)
	}

	var action models.Swiped = "like"
	if !request.RightSwipe {
		action = "pass"
	}

	if httpResponse := swipeService.SwipeHandler(subscriber.ID, request.ProfileId, subscriber.PremiumPackage); httpResponse.IsError() {
		return []types.ResponseSwipe{}, httpResponse
	}

	data, httpResponse := swipeService.FindProfiles(subscriber.ID, request.ProfileId, string(subscriber.Gender), action)
	if httpResponse.IsError() {
		return []types.ResponseSwipe{}, httpResponse
	}

	return data, helper.Success(i.response, common.Tracer(), data)
}

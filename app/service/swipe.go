package service

import (
	"context"
	"dating-be/app/domain/contract"
	"dating-be/app/domain/models"
	"dating-be/app/domain/types"
	"dating-be/app/helper"
	"dating-be/app/repository"
	"dating-be/common"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type swipeServiceContext struct {
	response common.Response
}

func NewSwipeService(response common.Response) contract.ISwipeService {
	return swipeServiceContext{response: response}
}

func (s swipeServiceContext) SwipeHandler(userId, profileId uint, premiumPackage models.PremiumPackage) common.Response {

	bgContext := context.Background()
	quotaKey := fmt.Sprintf("user:%d:swipe_quota", userId)
	if common.RedisClient.Exists(bgContext, quotaKey).Val() == 0 {
		common.RedisClient.Set(bgContext, quotaKey, 10, 24*time.Hour)
	}

	if premiumPackage != models.SwipeQuota {
		quota, err := common.RedisClient.Get(bgContext, quotaKey).Int()
		if err != nil {
			return helper.GeneralError(s.response, err, common.Tracer())
		}
		if quota <= 0 {
			newError := errors.New("limit quota exceeded")
			return helper.BadRequest(s.response, common.Tracer(), newError, "06")
		}
		common.RedisClient.Decr(bgContext, quotaKey)
	}

	return helper.Success(s.response, common.Tracer(), nil)
}

func (s swipeServiceContext) FindProfiles(userId, profileId uint, gender string, action models.Swiped) ([]types.ResponseSwipe, common.Response) {

	dateKey := time.Now().Format("2006-01-02")
	viewedProfilesKey := fmt.Sprintf("user:%d:viewed_profiles:%s", userId, dateKey)
	bgContext := context.Background()
	strUserId := strconv.Itoa(int(userId))

	viewedProfiles := common.RedisClient.SMembers(bgContext, viewedProfilesKey).Val()
	if viewedProfiles == nil {
		viewedProfiles = append(viewedProfiles, strUserId)
		common.RedisClient.SAdd(bgContext, viewedProfilesKey, profileId)
	}

	page, err := common.RedisClient.Get(bgContext, strUserId).Int()
	if err != nil {
		return []types.ResponseSwipe{}, helper.GeneralError(s.response, err, common.Tracer())
	}

	var subscribers []models.Subscriber
	subscribers, httpResponse := repository.NewDatingRepository(s.response).FindSubscriberIn(viewedProfiles, gender, page)
	if httpResponse.IsError() {
		return []types.ResponseSwipe{}, httpResponse
	}
	common.RedisClient.SAdd(bgContext, viewedProfilesKey, profileId)

	if len(subscribers) == 0 {
		if err := common.RedisClient.Set(context.Background(), strUserId, 1, 0).Err(); err != nil {
			return []types.ResponseSwipe{}, helper.GeneralError(s.response, err, common.Tracer())
		}
		subscribers, httpResponse = repository.NewDatingRepository(s.response).FindSubscriberIn(viewedProfiles, gender, page)
		if httpResponse.IsError() {
			return []types.ResponseSwipe{}, httpResponse
		}
	} else {
		common.RedisClient.Incr(bgContext, strUserId)
	}

	data := models.UserView{
		UserID:    userId,
		ProfileID: profileId,
		Swiped:    string(action),
	}
	if httpResponse := repository.NewDatingRepository(s.response).CreateUserView(data); httpResponse.IsError() {
		return []types.ResponseSwipe{}, httpResponse
	}

	var responses []types.ResponseSwipe
	for _, subscriber := range subscribers {
		response := types.ResponseSwipe{
			ProfileId:      subscriber.ID,
			FullName:       subscriber.FullName,
			Username:       subscriber.Username,
			Gender:         string(subscriber.Gender),
			Age:            subscriber.Age,
			Email:          subscriber.Email,
			PremiumPackage: string(subscriber.PremiumPackage),
		}
		common.RedisClient.SAdd(bgContext, viewedProfilesKey, response.ProfileId)
		responses = append(responses, response)
	}

	return responses, helper.Success(s.response, common.Tracer(), responses)
}

func (s swipeServiceContext) FirstViewProfile(userId uint, gender string) ([]types.ResponseSwipe, common.Response) {

	strUserId := strconv.Itoa(int(userId))
	bgContext := context.Background()
	page, err := common.RedisClient.Get(bgContext, strUserId).Int()
	if err != nil {
		return []types.ResponseSwipe{}, helper.GeneralError(s.response, err, common.Tracer())
	}

	viewedProfiles := []string{strUserId}
	subscribers, httpResponse := repository.NewDatingRepository(s.response).FindSubscriberIn(viewedProfiles, gender, page)
	if httpResponse.IsError() {
		return []types.ResponseSwipe{}, httpResponse
	}
	common.RedisClient.Incr(bgContext, strUserId)

	var responses []types.ResponseSwipe
	for _, subscriber := range subscribers {
		response := types.ResponseSwipe{
			ProfileId:      subscriber.ID,
			FullName:       subscriber.FullName,
			Username:       subscriber.Username,
			Gender:         string(subscriber.Gender),
			Age:            subscriber.Age,
			Email:          subscriber.Email,
			PremiumPackage: string(subscriber.PremiumPackage),
		}
		responses = append(responses, response)
	}

	viewedProfilesKey := fmt.Sprintf("user:%d:viewed_profiles:%s", userId, time.Now().Format("2006-01-02"))
	common.RedisClient.SAdd(bgContext, viewedProfilesKey, userId)
	common.RedisClient.Expire(bgContext, viewedProfilesKey, 24*time.Hour)

	return responses, helper.Success(s.response, common.Tracer(), responses)
}

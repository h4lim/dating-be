package repository

import (
	"dating-be/app/domain/contract"
	"dating-be/app/domain/models"
	"dating-be/app/helper"
	"dating-be/common"
	"fmt"
	"net/http"
)

const GORM_DATA_NOT_FOUND = "record not found"

type datingRepositoryContext struct {
	response common.Response
}

func NewDatingRepository(response common.Response) contract.IDatingRepository {
	return datingRepositoryContext{response: response}
}

func (d datingRepositoryContext) CreateSubscriber(data models.Subscriber) common.Response {

	if err := common.GormDB.Create(&data).Error; err != nil {
		return helper.GeneralError(d.response, err, common.Tracer())
	}

	return d.response.SetSuccess(common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusOK, Code: "0", Data: data})
}

func (d datingRepositoryContext) FindSubscriberByUsername(username string) (*models.Subscriber, common.Response) {

	var data models.Subscriber
	if err := common.GormDB.Where("username = ?",
		username).First(&data).Error; err != nil {
		if fmt.Sprintf("%v", err) == GORM_DATA_NOT_FOUND {
			return nil, d.response.SetSuccess(common.Tracer(),
				common.HttpResponse{Code: "0", Data: nil})
		}
		return nil, helper.GeneralError(d.response, err, common.Tracer())
	}

	return &data, d.response.SetSuccess(common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusOK, Code: "0", Data: data})
}

func (d datingRepositoryContext) FindSubscriberByUsernameOrEmail(username string, email string) (*models.Subscriber, common.Response) {

	var data models.Subscriber
	if err := common.GormDB.Where("username = ? OR email = ? ",
		username, email).First(&data).Error; err != nil {
		if fmt.Sprintf("%v", err) == GORM_DATA_NOT_FOUND {
			return nil, d.response.SetSuccess(common.Tracer(),
				common.HttpResponse{Code: "0", Data: nil})
		}
		return nil, helper.GeneralError(d.response, err, common.Tracer())
	}

	return &data, d.response.SetSuccess(common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusOK, Code: "0", Data: data})
}

func (d datingRepositoryContext) CreateUserView(data models.UserView) common.Response {

	if err := common.GormDB.Create(&data).Error; err != nil {
		return helper.GeneralError(d.response, err, common.Tracer())
	}

	return d.response.SetSuccess(common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusOK, Code: "0", Data: data})
}

func (d datingRepositoryContext) FindSubscriberIn(viewedProfiles []string, gender string, page int) ([]models.Subscriber, common.Response) {
	var data []models.Subscriber
	offset := (page - 1) * 1
	if err := common.GormDB.Debug().Offset(offset).Limit(1).Where("id NOT IN ? AND gender != ? ",
		viewedProfiles, gender).Find(&data).Error; err != nil {
		if fmt.Sprintf("%v", err) == GORM_DATA_NOT_FOUND {
			return nil, d.response.SetSuccess(common.Tracer(),
				common.HttpResponse{Code: "0", Data: nil})
		}
		return nil, helper.GeneralError(d.response, err, common.Tracer())
	}

	return data, d.response.SetSuccess(common.Tracer(),
		common.HttpResponse{HttpCode: http.StatusOK, Code: "0", Data: data})
}

package helper

import (
	"dating-be/common"
	"net/http"
)

func Success(existingResponse common.Response, tracer common.TracerModel, data any) common.Response {
	if data == nil {
		return existingResponse.SetSuccess(tracer, common.HttpResponse{HttpCode: http.StatusOK, Code: "0"})
	}
	return existingResponse.SetSuccess(tracer, common.HttpResponse{HttpCode: http.StatusOK, Code: "0", Data: data})
}

func BadRequest(existingResponse common.Response, tracer common.TracerModel, err error, code string) common.Response {
	return existingResponse.SetError(&err, tracer,
		common.HttpResponse{HttpCode: http.StatusBadRequest, Code: code})
}

func GeneralError(existingResponse common.Response, err error, tracer common.TracerModel) common.Response {
	return existingResponse.SetError(&err, tracer,
		common.HttpResponse{HttpCode: http.StatusInternalServerError, Code: "99", Message: "General Error"})
}

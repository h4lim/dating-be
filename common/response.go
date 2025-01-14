package common

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"strings"
	"time"
)

type GinResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"traceId"`
}

type GinResponseData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"traceId"`
	Data    any    `json:"data"`
}

type TracerModel struct {
	FunctionName  string `json:"function_name"`
	FileName      string `json:"file_name"`
	Line          int    `json:"line"`
	Step          int    `json:"step"`
	Duration      string `json:"duration"`
	TotalDuration string `json:"total_duration"`
}

type Response struct {
	HttpCode         int
	Code             string
	Message          string
	Data             any
	Error            *error
	Tracer           TracerModel
	TraceId          int64
	Language         string
	AdditionalTracer []string
}

type HttpResponse struct {
	HttpCode int
	Code     string
	Message  string
	Data     any
}

func InitResponse(traceId int64, language string) Response {
	return Response{
		TraceId:  traceId,
		Language: language,
	}
}

func (r *Response) SetSuccess(Tracer TracerModel, response ...HttpResponse) Response {
	r.Tracer = Tracer

	var responseData HttpResponse
	if len(response) > 0 {
		responseData = response[len(response)-1]
	}

	if responseData.HttpCode == 0 {
		responseData.HttpCode = 200
	}

	if responseData.Code == "" {
		responseData.Code = "0"
	}

	r.HttpCode = responseData.HttpCode
	r.Code = responseData.Code
	r.Message = responseData.Message
	r.Data = responseData.Data

	if r.Message == "" {
		getMessage(r)
	}

	r.debug(true)

	return *r
}

func (r *Response) SetError(Error *error, Tracer TracerModel, httpResponse ...HttpResponse) Response {

	r.Error = Error
	r.Tracer = Tracer

	var responseData HttpResponse
	if len(httpResponse) > 0 {
		responseData = httpResponse[len(httpResponse)-1]
	}

	if responseData.HttpCode == 0 {
		responseData.HttpCode = 400
	}

	if responseData.Code == "" {
		responseData.Code = "99"
	}

	r.HttpCode = responseData.HttpCode
	r.Code = responseData.Code
	r.Message = responseData.Message
	r.Data = responseData.Data

	if r.Message == "" {
		getMessage(r)
	}

	r.debug(true)

	return *r
}

func (r *Response) SetAdditionalTracer(additionalTracer string) Response {
	r.AdditionalTracer = append(r.AdditionalTracer, additionalTracer)
	return *r
}

func (r Response) BuildResponse() (int, any) {

	r.debug(true)

	delete(UnixTimestamp, r.TraceId)
	delete(Step, r.TraceId)
	delete(RequestId, r.TraceId)

	return r.HttpCode, GinResponse{
		Code:    r.Code,
		Message: r.Message,
		TraceId: strconv.FormatInt(r.TraceId, 10),
	}
}

func (r Response) BuildResponseData(data any) (int, any) {

	r.Data = data
	r.debug(true)

	delete(UnixTimestamp, r.TraceId)
	delete(Step, r.TraceId)
	delete(RequestId, r.TraceId)

	return r.HttpCode, GinResponseData{
		Code:    r.Code,
		Message: r.Message,
		TraceId: strconv.FormatInt(r.TraceId, 10),
		Data:    r.Data,
	}
}

func (r *Response) IsError() bool {
	return r.Error != nil
}

func (r *Response) debug(nextStep bool) {

	if ZapLog != nil {

		duration := time.Now().UnixNano() - r.TraceId
		ms := float64(duration) / float64(time.Millisecond)
		zapFields := []zapcore.Field{}

		if nextStep {
			zapFields = append(zapFields, zap.String("step", GetNextStep(r.TraceId)))
		} else {
			zapFields = append(zapFields, zap.String("step", GetStep(r.TraceId)))
		}

		zapFields = append(zapFields, zap.String("duration", GetDuration(r.TraceId)+" ms"))
		zapFields = append(zapFields, zap.String("total-duration", fmt.Sprintf("%v", ms)+" ms"))
		zapFields = append(zapFields, zap.String("additional-tracer", strings.Join(r.AdditionalTracer, " ")))
		zapFields = append(zapFields, zap.String("filename", r.Tracer.FileName))
		zapFields = append(zapFields, zap.String("function-name", r.Tracer.FunctionName))
		zapFields = append(zapFields, zap.Int("line", r.Tracer.Line))
		zapFields = append(zapFields, zap.String("trace", r.Tracer.FileName+":"+strconv.Itoa(r.Tracer.Line)))
		zapFields = append(zapFields, zap.String("message ", getRawResponse(r)))

		if r.Error != nil {
			zapFields = append(zapFields, zap.String("error", fmt.Sprintf("%v", *r.Error)))
			ZapLog.Debug(strconv.FormatInt(r.TraceId, 10), zapFields...)
		} else {
			ZapLog.Info(strconv.FormatInt(r.TraceId, 10), zapFields...)
		}
	}
}

func getRawResponse(r *Response) string {

	rawData := ""
	if r.Data != nil {
		response := GinResponseData{
			Code:    r.Code,
			Message: r.Message,
			TraceId: strconv.FormatInt(r.TraceId, 10),
			Data:    r.Data,
		}
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			rawData = err.Error()
		}
		rawData = string(jsonBytes)
	} else {
		response := GinResponse{
			Code:    r.Code,
			Message: r.Message,
			TraceId: strconv.FormatInt(r.TraceId, 10),
		}
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			rawData = err.Error()
		}
		rawData = string(jsonBytes)
	}

	return rawData

}

func getMessage(r *Response) {

	message := MessageMap[strings.ToLower(r.Language)][r.Code]
	if message == "" {
		message = "unknown message"
	}
	r.Message = message

}

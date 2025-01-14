package types

type HttpResponse struct {
	HttpCode int
	Response Response
}
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

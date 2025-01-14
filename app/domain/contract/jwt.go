package contract

import "dating-be/common"

type IJwtService interface {
	GetToken(username string, password string) (*string, common.Response)
	VerifyToken(token string) common.Response
}

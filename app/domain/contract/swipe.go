package contract

import (
	"dating-be/app/domain/models"
	"dating-be/app/domain/types"
	"dating-be/common"
)

type ISwipeService interface {
	SwipeHandler(userId, profileId uint, premiumPackage models.PremiumPackage) common.Response
	FindProfiles(userId, profileId uint, gender string, action models.Swiped) ([]types.ResponseSwipe, common.Response)
	FirstViewProfile(userId uint, gender string) ([]types.ResponseSwipe, common.Response)
}

package models

import "gorm.io/gorm"

type Gender string
type PremiumPackage string
type Swiped string

const (
	Male          Gender         = "male"
	Female        Gender         = "female"
	NonPremium    PremiumPackage = "non_premium"
	SwipeQuota    PremiumPackage = "swipe_quota"
	VerifiedLabel PremiumPackage = "verified_label"
	Like          Swiped         = "like"
	Pass          Swiped         = "pass"
)

type Subscriber struct {
	gorm.Model
	FullName       string         `db:"full_name"`
	Username       string         `db:"username"`
	Gender         Gender         `db:"gender"`
	Age            int8           `db:"age"`
	Email          string         `db:"email"`
	Password       string         `db:"password"`
	PremiumPackage PremiumPackage `db:"premium_package"`
	Salt           string         `db:"salt"`
	ProfilesViewed []UserView     `gorm:"foreignKey:UserID"`
}

type UserView struct {
	gorm.Model
	UserID    uint   `gorm:"not null"` // The user viewing the profile
	ProfileID uint   `gorm:"not null"` // The profile being viewed
	Swiped    string `db:"swiped"`
}

package types

type RequestSignup struct {
	FullName        string `json:"full_name"`
	Gender          string `json:"gender"`
	Age             int8   `json:"age"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Type     string `json:"type"`
	Expired  int    `json:"expired"`
}

type RequestSwipe struct {
	Username    string `json:"username"`
	IsFirstView bool   `json:"is_first_view"`
	ProfileId   uint   `json:"profile_id"`
	RightSwipe  bool   `json:"right_swipe"`
}

type ResponseSwipe struct {
	ProfileId      uint   `json:"profile_id"`
	FullName       string `json:"full_name"`
	Username       string `json:"username"`
	Gender         string `json:"gender"`
	Age            int8   `json:"age"`
	Email          string `json:"email"`
	PremiumPackage string `json:"premium_package"`
}

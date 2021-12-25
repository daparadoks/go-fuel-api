package responses

import "time"

type LoginResponse struct {
	Token    string
	Username string
	MemberId uint
}

type RegisterResponse struct {
	Id       uint
	Username string
}

type MemberResponse struct {
	Id          uint
	Username    string
	Mail        string
	Token       string
	DeviceToken string
	IsGues      bool
	LastLogin   time.Time
	Password    string
}

func InitGuest(deviceToken string) MemberResponse {
	var response MemberResponse
	response.DeviceToken = deviceToken
	response.IsGues = true

	return response
}

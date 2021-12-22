package responses

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
	Id       uint
	Username string
	Mail     string
}

package requests

type LoginRequest struct {
	Username    string
	Password    string
	DeviceToken string
	RememberMe  bool
}

type RegisterRequest struct {
	Username        string
	Password        string
	ConfirmPassword string
	Mail            string
	DeviceToken     string
}

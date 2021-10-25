package requests

type LoginRequest struct {
	Username    string
	Password    string
	DeviceToken string
	RememberMe  bool
}

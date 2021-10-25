package http

import (
	"encoding/json"
	"net/http"

	"github.com/daparadoks/go-fuel-api/internal/requests"
	"github.com/daparadoks/go-fuel-api/internal/responses"
)

// Login -
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var request requests.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		GetErrorResponse(w, "Failed to decode json from body")
		return
	}
	member, err := h.MemberService.GetMember(request.Username)
	if err != nil {
		GetErrorResponse(w, "Failed to login")
		return
	}
	if member.Password != request.Password {
		GetErrorResponse(w, "Username or password is invalid")
		return
	}

	token, err := h.MemberService.GetTokenByMemberId(request.DeviceToken, member.ID)
	if err != nil {
		GetErrorResponse(w, "Failed to login")
		return
	}

	var loginResponse = responses.LoginResponse{
		MemberId: member.ID,
		Username: member.Username,
		Token:    token.Token,
	}
	SendOkResponse(w, loginResponse)
}

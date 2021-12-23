package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/daparadoks/go-fuel-api/internal/member"
	"github.com/daparadoks/go-fuel-api/internal/requests"
	"github.com/daparadoks/go-fuel-api/internal/responses"
)

// Login -
func (h *Handler) Login(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SetHeaders(w)
	var request requests.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		GetErrorResponseWithError(w, err)
		return
	}
	var validationErrorMessages []string
	if request.Password != request.ConfirmPassword || len(request.Password) == 0 {
		validationErrorMessages = append(validationErrorMessages, "Geçersiz şifre")
	}
	if len(request.Username) == 0 {
		validationErrorMessages = append(validationErrorMessages, "Geçersiz kullanıcı adı")
	}
	if len(request.Mail) == 0 {
		validationErrorMessages = append(validationErrorMessages, "Geçersiz email")
	}

	if len(validationErrorMessages) > 0 {
		GetErrorResponse(w, strings.Join(validationErrorMessages[:], ","))
		return
	}
	existsMember, _ := h.MemberService.GetMember(request.Username)
	if existsMember.ID > 0 {
		GetErrorResponse(w, request.Username+" kullanıcı adıyla kayıtlı bir üyelik mevcuttur.")
		return
	}
	existsMember, _ = h.MemberService.GetMemberByMail(request.Mail)
	if existsMember.ID > 0 {
		GetErrorResponse(w, " Bu mail adresiyle kayıtlı bir üyelik mevcuttur: "+request.Mail)
		return
	}
	var member member.Member
	member.Username = request.Username
	member.Password = request.Password
	member.Mail = request.Mail
	member.CreatedAt = time.Now()

	member, err := h.MemberService.Register(member)
	var response responses.RegisterResponse

	if err != nil {
		GetErrorResponseWithError(w, err)
		return
	}
	if member.ID == 0 {
		GetErrorResponse(w, "Üye kayıt sırasında bir hata oluştu")
		return
	}
	response.Id = member.ID
	response.Username = member.Username

	SendOkResponse(w, response)
}

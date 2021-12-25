package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/daparadoks/go-fuel-api/internal/member"
	redisService "github.com/daparadoks/go-fuel-api/internal/redis"
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
	passwordHash := md5.Sum([]byte(request.Password))
	if member.Password != hex.EncodeToString(passwordHash[:]) {
		GetErrorResponse(w, "Username or password is invalid")
		return
	}

	now := time.Now()
	token, err := h.MemberService.GetTokenByMemberId(m.Token, member.ID)
	if err != nil || token.ExpireDate.Before(now) {
		GetErrorResponse(w, "Failed to login")
		return
	}

	m.Id = member.ID
	m.Username = member.Username
	m.Mail = member.Mail
	m.LastLogin = time.Now()
	m.Token = token.Token
	m.IsGues = false
	m.Password = member.Password

	jsonData, _ := json.Marshal(m)
	redisService.Set(context.Background(), "LoginInfo_"+m.Token, string(jsonData))
	m.Password = ""
	SendOkResponse(w, m)
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
	passwordHash := md5.Sum([]byte(request.Password))
	var member member.Member
	member.Username = request.Username
	member.Password = hex.EncodeToString(passwordHash[:])
	member.Mail = request.Mail
	member.CreatedAt = time.Now()

	member, err := h.MemberService.AddOrUpdateMember(member)
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

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	member, err := h.MemberService.GetMemberById(m.Id)
	if err != nil {
		GetErrorResponseWithError(w, err)
		return
	}

	newPassword := md5.Sum([]byte("1234"))
	member.Password = hex.EncodeToString(newPassword[:])
	member, err = h.MemberService.AddOrUpdateMember(member)
	if err != nil {
		GetErrorResponseWithError(w, err)
		return
	}

	SendOkResponse(w, true)
}

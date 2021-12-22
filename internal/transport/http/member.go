package http

import (
	"net/http"

	"github.com/daparadoks/go-fuel-api/internal/responses"
)

func (h *Handler) GetMember(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header["Token"]
	token := tokenHeader[0]
	member, err := h.MemberService.GetMember(token)
	if err != nil {
		GetErrorResponseWithError(w, err)
		return
	}
	var response responses.MemberResponse
	response.Id = member.ID
	response.Username = member.Username
	response.Mail = member.Mail
}

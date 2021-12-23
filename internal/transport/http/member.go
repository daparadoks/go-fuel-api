package http

import (
	"net/http"

	"github.com/daparadoks/go-fuel-api/internal/responses"
)

func (h *Handler) GetMember(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SendOkResponse(w, m)
}

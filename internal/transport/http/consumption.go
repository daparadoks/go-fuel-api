package http

import (
	"net/http"

	"github.com/daparadoks/go-fuel-api/internal/responses"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {

	SendOkResponse(w, nil)
}

package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/daparadoks/go-fuel-api/internal/consumption"
	"github.com/daparadoks/go-fuel-api/internal/responses"
	"github.com/gorilla/mux"
)

func (h *Handler) Consumptions(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SetHeaders(w)
	consumptions, err := h.ConsumptionService.GetList(m.Id)
	if err != nil {
		GetErrorResponseWithError(w, err)
		return
	}

	model := responses.InitConsumptionListModel(consumptions)
	SendOkResponse(w, model)
}

func (h *Handler) Consumption(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SetHeaders(w)
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		GetErrorResponse(w, "Geçersiz id")
		return
	}
	consumptionId, _ := strconv.ParseUint(id, 10, 64)
	consumption, err := h.ConsumptionService.Get(uint(consumptionId))
	if err != nil {
		GetErrorResponseWithCode(w, "Yakıt bilgisi bulunamadı", 404)
		return
	}
	SendOkResponse(w, consumption)
}

func (h *Handler) ConsumptionAdd(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SetHeaders(w)
	var consumption consumption.Consumption
	if err := json.NewDecoder(r.Body).Decode(&consumption); err != nil {
		GetErrorResponse(w, "Geçersiz bir talepte bulundunuz: "+err.Error())
		return
	}
	consumption.MemberId = m.Id
	defaultTime := time.Time{}
	if consumption.FuelupDate == defaultTime {
		consumption.FuelupDate = time.Now()
	}
	consumption, err := h.ConsumptionService.Add(consumption)

	if err != nil {
		GetErrorResponse(w, "Ekleme işlemi başarısız")
		return
	}

	SendOkResponse(w, consumption)
}

func (h *Handler) ConsumptionUpdate(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SetHeaders(w)
	var consumption consumption.Consumption
	if err := json.NewDecoder(r.Body).Decode(&consumption); err != nil {
		GetErrorResponse(w, "Geçersiz bir talepte bulundunuz")
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	conId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		GetErrorResponse(w, "Geçersiz id")
		return
	}
	consumption, err = h.ConsumptionService.Update(uint(conId), consumption)
	if err != nil {
		GetErrorResponse(w, "Güncelleme işlemi başarısız")
		return
	}

	SendOkResponse(w, consumption)
}

func (h *Handler) ConsumptionDelete(w http.ResponseWriter, r *http.Request, m responses.MemberResponse) {
	SetHeaders(w)
	vars := mux.Vars(r)
	id := vars["id"]
	conId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		GetErrorResponse(w, "Geçersiz id")
		return
	}

	consumption, _ := h.ConsumptionService.Get(uint(conId))
	if consumption.ID <= 0 || consumption.MemberId != m.Id {
		GetErrorResponseWithCode(w, "Yakıt bilgisi bulunamadı", 404)
		return
	}
	err = h.ConsumptionService.Delete(uint(conId))
	if err != nil {
		GetErrorResponse(w, "Silme işlemi gerçekleştirilemedi")
		return
	}

	GetErrorResponse(w, "Silme işlemi başarılı")
}

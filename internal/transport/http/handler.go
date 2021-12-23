package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/daparadoks/go-fuel-api/internal/consumption"
	"github.com/daparadoks/go-fuel-api/internal/member"
	"github.com/daparadoks/go-fuel-api/internal/responses"
	jwt "github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - stores pointer to our comments service
type Handler struct {
	Router             *mux.Router
	MemberService      *member.Service
	ConsumptionService *consumption.Service
}

type Response struct {
	Success bool
	Message string
	Code    int
}

// NewHandler - returns a pointer to a Handler
func NewHandler(memberService *member.Service, consumptionService *consumption.Service) *Handler {
	return &Handler{
		MemberService:      memberService,
		ConsumptionService: consumptionService,
	}
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
		})
		log.Info("Endpoint hit!")
		next.ServeHTTP(w, r)
	})
}

// BasicAuth - a handy middleware function that logs out incoming requests
func BasicAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if user == "paradox" && pass == "12345" && ok {
			original(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			GetErrorResponseWithCode(w, "not authorized", 401)
		}
	}
}

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request, m responses.MemberResponse)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("jwt auth endpoint hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "not authorized", 401)
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "not authorized", 401)
			return
		}

		if validateToken(authHeaderParts[1], "") {
			deviceTokenHeader := r.Header["DeviceToken"]
			deviceToken := deviceTokenHeader[0]
			currentMember := responses.InitGuest(deviceToken)
			original(w, r, currentMember)
		} else {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "not authorized", 401)
		}
	}
}

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func UserAuth(original func(w http.ResponseWriter, r *http.Request, m responses.MemberResponse), h *Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("user auth")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "authorization required", 401)
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "authorization is not valid", 401)
			return
		}

		if validateToken(authHeaderParts[1], "") {
			tokenHeader := r.Header["Token"]
			token := tokenHeader[0]
			memberToken, err := h.MemberService.GetToken(token)
			if err != nil || memberToken.Token != token {
				SetHeaders(w)
				GetErrorResponseWithCode(w, "user token is not valid", 401)
				return
			}

			member, memberErr := h.MemberService.GetMemberById(memberToken.MemberId)
			if memberErr != nil {
				SetHeaders(w)
				GetErrorResponseWithCode(w, "user token is not valid", 401)
				return
			}

			deviceTokenHeader := r.Header["Authorization"]
			deviceToken := deviceTokenHeader[0]

			var currentMember responses.MemberResponse
			currentMember.Id = memberToken.MemberId
			currentMember.Mail = member.Mail
			currentMember.Username = member.Username
			currentMember.Token = memberToken.Token
			currentMember.DeviceToken = deviceToken
			original(w, r, currentMember)
		} else {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "authorization validation has failed", 401)
			return
		}
	}
}

// validateToken - validates an incoming jwt token
func validateToken(accessToken string, userKey string) bool {
	if len(userKey) == 0 {
		userKey = "paradox"
	}

	var mySigningKey = []byte(userKey)
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func SendOkResponse(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func GetErrorResponseWithError(w http.ResponseWriter, err error) {
	GetErrorResponseWithCode(w, err.Error(), http.StatusBadRequest)
}
func GetErrorResponse(w http.ResponseWriter, message string) {
	GetErrorResponseWithCode(w, message, http.StatusBadRequest)
}
func GetErrorResponseWithCode(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(Response{Success: false, Message: message, Code: code}); err != nil {
		panic(err)
	}
}

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/login", JWTAuth(h.Login)).Methods("POST")

	h.Router.HandleFunc("/api/member", UserAuth(h.GetMember, h)).Methods("GET")
	h.Router.HandleFunc("/api/member", JWTAuth(h.Register)).Methods("POST")

	h.Router.HandleFunc("/api/consumptions", UserAuth(h.Consumptions, h)).Methods("GET")
	h.Router.HandleFunc("/api/consumption/{id}", UserAuth(h.Consumption, h)).Methods("GET")
	h.Router.HandleFunc("/api/consumption", UserAuth(h.ConsumptionAdd, h)).Methods("POST")
	h.Router.HandleFunc("/api/consumption/{id}", UserAuth(h.ConsumptionUpdate, h)).Methods("PUT")
	h.Router.HandleFunc("/api/consumption/{id}", UserAuth(h.ConsumptionDelete, h)).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Success: true, Message: "I'm alive!", Code: http.StatusOK}); err != nil {
			panic(err)
		}
	})
}

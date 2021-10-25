package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/daparadoks/go-fuel-api/internal/member"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - stores pointer to our comments service
type Handler struct {
	Router        *mux.Router
	MemberService *member.Service
}

type Response struct {
	Success bool
	Message string
	Code    int
}

// NewHandler - returns a pointer to a Handler
func NewHandler(memberService *member.Service) *Handler {
	return &Handler{
		MemberService: memberService,
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
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
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

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			SetHeaders(w)
			GetErrorResponseWithCode(w, "not authorized", 401)
		}
	}
}

// validateToken - validates an incoming jwt token
func validateToken(accessToken string) bool {
	// replace this by loading in a private RSA cert for more security
	var mySigningKey = []byte("missionimpossible")
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

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/login", BasicAuth(h.Login)).Methods("POST")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Success: true, Message: "I'm alive!", Code: http.StatusOK}); err != nil {
			panic(err)
		}
	})
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

func GetErrorResponse(w http.ResponseWriter, message string) {
	GetErrorResponseWithCode(w, message, http.StatusBadRequest)
}
func GetErrorResponseWithCode(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(Response{Success: false, Message: message, Code: code}); err != nil {
		panic(err)
	}
}

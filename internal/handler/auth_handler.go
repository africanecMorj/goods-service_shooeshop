package handler

import (
	"encoding/json"
	"net/http"
	"slices"
	"net"
	"log"

	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
)

var reliableIps = []string{
	"93.171.247.178",
	"::1",
}

type AuthHandler struct{ S *service.AuthService }

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	if in.Email == "" || in.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	if err := h.S.Register(r.Context(), in.Email, in.Password); err != nil {
		http.Error(w, "Cannot make record into db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "successfully registered",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if err:=json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	
	role := "user"

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	log.Println(host)
	if err == nil && slices.Contains(reliableIps, host) {
		role = "admin"
	}

	a, rfr, err := h.S.Login(r.Context(), in.Email, in.Password, role)
	if err != nil {
		http.Error(w, "invalid", 401)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access": a, "refresh": rfr})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var in struct{ Token string }
	if err:=json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Bad ruquest", http.StatusBadRequest)
	}
		role := "user"

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && slices.Contains(reliableIps, host) {
		role = "admin"
	}

	a, rfr, err := h.S.Refresh(r.Context(), in.Token, role)
	if err != nil {
		http.Error(w, "invalid", 401)
		return
	}
   
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access": a, "refresh": rfr})
}
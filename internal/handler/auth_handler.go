package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
)

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

	role := "user"
	if err := h.S.Register(r.Context(), in.Email, in.Password, role); err != nil {
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
	

	a, rfr, err := h.S.Login(r.Context(), in.Email, in.Password)
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


	a, rfr, err := h.S.Refresh(r.Context(), in.Token)
	if err != nil {
		http.Error(w, "invalid", 401)
		return
	}
   
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access": a, "refresh": rfr})
}

func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset := (page - 1) * limit

	users, total, err := h.S.GetUsers(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	var respUsers []map[string]interface{}
	for _, u := range users {
		respUsers = append(respUsers, map[string]interface{}{
			"id":    u.ID,
			"email": u.Email,
			"role":  u.Role,
		})
	}

	resp := map[string]interface{}{
		"page":     page,
		"limit":    limit,
		"total":    total,
		"users": respUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service/utils"
)

type UpdateHandler struct {
	AuthService *service.AuthService
	UpdateService *service.UpdateService

}

func sendOk (w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "letter is sent",
	})
}

func (h *UpdateHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
    var req domain.ResetRequest
  	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	code, err := h.UpdateService.StoreCode(r.Context(), req.Email)
	if err != nil {
    	sendOk(w)
		return
	}

    err = utils.SendEmail(req.Email,
        "Password reset",
        "Use this code: "+code )
	if err != nil {
		sendOk(w)
		return
	}

    sendOk(w)
}

func (h *UpdateHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req domain.ResetPassword
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	code := chi.URLParam(r, "code")

	if err := h.UpdateService.ResetPassword(r.Context(), req.Email, code, req.Password); err != nil {
		fmt.Println(err)
		if err != fmt.Errorf("Invalid code") {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
        	return
		}
		http.Error(w, "Invalid code", http.StatusBadRequest)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "succesfully updated",
	})
}

func (h *UpdateHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if err := h.UpdateService.PatchUser(r.Context(), id, payload); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "patched",
	})
}


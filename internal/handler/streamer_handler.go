package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
)

type StreamerHandler struct {
	Service *service.StreamerService
}

func (h *StreamerHandler) PatchImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image required", 400)
		return
	}
	defer file.Close()

	path, err := h.Service.UpdateImage(
		r.Context(),
		id,
		file,
		header,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"image_path": path,
	})
}

func (h *StreamerHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	stream, err := h.Service.GetImageReader(ctx, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	defer stream.Reader.Close()

	if f, ok := stream.Reader.(*os.File); ok {
		if stat, err := f.Stat(); err == nil {
			w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		}
	}

	w.Header().Set("Content-Type", stream.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=86400") 

	if _, err := io.Copy(w, stream.Reader); err != nil {
		return
	}
}
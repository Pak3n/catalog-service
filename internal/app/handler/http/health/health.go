package health

import (
	"log"
	"net/http"

	rhandler "github.com/Pak3n/catalog-service/internal/app/handler/http"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("Failed to write health response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

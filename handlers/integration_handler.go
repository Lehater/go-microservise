// handlers/integration_handler.go
package handlers

import (
	"context"
	"net/http"
	"time"

	"go-microservice/services"

	"github.com/gorilla/mux"
)

type IntegrationHandler struct {
	svc *services.IntegrationService
}

func NewIntegrationHandler(svc *services.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{svc: svc}
}

// RegisterIntegrationRoutes — один тестовый эндпоинт.
func RegisterIntegrationRoutes(r *mux.Router, svc *services.IntegrationService) {
	h := NewIntegrationHandler(svc)
	r.HandleFunc("/integration/upload-test", h.UploadTest).Methods(http.MethodPost)
}

func (h *IntegrationHandler) UploadTest(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.EnsureBucket(ctx); err != nil {
		http.Error(w, "failed to ensure bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.svc.UploadTestObject(ctx, "test.txt", "hello from go-microservice"); err != nil {
		http.Error(w, "failed to upload object: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

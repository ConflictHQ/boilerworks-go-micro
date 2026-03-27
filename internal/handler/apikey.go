package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ConflictHQ/boilerworks-go-micro/internal/database/queries"
)

type ApiKeyHandler struct {
	q *queries.Queries
}

func NewApiKeyHandler(q *queries.Queries) *ApiKeyHandler {
	return &ApiKeyHandler{q: q}
}

type CreateApiKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

func (h *ApiKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"invalid request body"},
		})
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"name is required"},
		})
		return
	}

	if len(req.Scopes) == 0 {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"at least one scope is required"},
		})
		return
	}

	// Generate a random API key
	plaintext, err := generateAPIKey()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to generate API key"},
		})
		return
	}

	hash := sha256.Sum256([]byte(plaintext))
	keyHash := hex.EncodeToString(hash[:])

	apiKey, err := h.q.CreateApiKey(r.Context(), queries.CreateApiKeyParams{
		Name:    req.Name,
		KeyHash: keyHash,
		Scopes:  req.Scopes,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to create API key"},
		})
		return
	}

	writeJSON(w, http.StatusCreated, ApiResponse{
		Ok: true,
		Data: map[string]interface{}{
			"id":        apiKey.ID,
			"name":      apiKey.Name,
			"key":       plaintext,
			"scopes":    apiKey.Scopes,
			"createdAt": apiKey.CreatedAt,
		},
		Message: "Store this key securely -- it will not be shown again",
	})
}

func (h *ApiKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	keys, err := h.q.ListApiKeys(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to list API keys"},
		})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{
		Ok:   true,
		Data: keys,
	})
}

func (h *ApiKeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"invalid API key ID"},
		})
		return
	}

	if err := h.q.RevokeApiKey(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to revoke API key"},
		})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{
		Ok:      true,
		Message: "API key revoked",
	})
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("bw_%s", hex.EncodeToString(bytes)), nil
}

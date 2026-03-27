package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ConflictHQ/boilerworks-go-micro/internal/database/queries"
)

type EventHandler struct {
	q *queries.Queries
}

func NewEventHandler(q *queries.Queries) *EventHandler {
	return &EventHandler{q: q}
}

type CreateEventRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"invalid request body"},
		})
		return
	}

	if req.Type == "" {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"type is required"},
		})
		return
	}

	payload := req.Payload
	if payload == nil {
		payload = json.RawMessage(`{}`)
	}

	event, err := h.q.CreateEvent(r.Context(), queries.CreateEventParams{
		Type:    req.Type,
		Payload: payload,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to create event"},
		})
		return
	}

	writeJSON(w, http.StatusCreated, ApiResponse{
		Ok:   true,
		Data: event,
	})
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	eventType := r.URL.Query().Get("type")

	limit := int32(50)
	offset := int32(0)

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = int32(l)
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = int32(o)
		}
	}

	var events []queries.Event
	var err error

	if eventType != "" {
		events, err = h.q.ListEventsByType(r.Context(), queries.ListEventsByTypeParams{
			Type:   eventType,
			Limit:  limit,
			Offset: offset,
		})
	} else {
		events, err = h.q.ListEvents(r.Context(), queries.ListEventsParams{
			Limit:  limit,
			Offset: offset,
		})
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to list events"},
		})
		return
	}

	count, _ := h.q.CountEvents(r.Context())

	writeJSON(w, http.StatusOK, ApiResponse{
		Ok: true,
		Data: map[string]interface{}{
			"events": events,
			"total":  count,
		},
	})
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"invalid event ID"},
		})
		return
	}

	event, err := h.q.GetEvent(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ApiResponse{
			Ok:      false,
			Message: "event not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{
		Ok:   true,
		Data: event,
	})
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Ok:     false,
			Errors: []string{"invalid event ID"},
		})
		return
	}

	if err := h.q.SoftDeleteEvent(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, ApiResponse{
			Ok:     false,
			Errors: []string{"failed to delete event"},
		})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{
		Ok:      true,
		Message: "event deleted",
	})
}

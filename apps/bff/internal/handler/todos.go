package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	pbtasks "github.com/its-me-mayday/sisyphus/bff/proto/tasks"
)

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	listID := chi.URLParam(r, "listId")
	var body struct {
		Text     string `json:"text"`
		Priority string `json:"priority"`
		DueDate  string `json:"due_date"`
	}
	if err := decode(r, &body); err != nil || body.Text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	resp, err := h.tasks.CreateTodo(r.Context(), &pbtasks.CreateTodoRequest{
		ListId:   listID,
		Text:     body.Text,
		Priority: fromStringPriority(body.Priority),
		DueDate:  body.DueDate,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, resp.Todo)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Text     string `json:"text"`
		Done     bool   `json:"done"`
		Priority string `json:"priority"`
		DueDate  string `json:"due_date"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	resp, err := h.tasks.UpdateTodo(r.Context(), &pbtasks.UpdateTodoRequest{
		Id:       id,
		Text:     body.Text,
		Done:     body.Done,
		Priority: fromStringPriority(body.Priority),
		DueDate:  body.DueDate,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Todo)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.tasks.DeleteTodo(r.Context(), &pbtasks.DeleteTodoRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func fromStringPriority(s string) pbtasks.Priority {
	switch s {
	case "low":
		return pbtasks.Priority_LOW
	case "high":
		return pbtasks.Priority_HIGH
	default:
		return pbtasks.Priority_MEDIUM
	}
}

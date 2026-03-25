package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	pbtasks "github.com/its-me-mayday/sisyphus/bff/proto/tasks"
)

func (h *Handler) CreateList(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := decode(r, &body); err != nil || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	resp, err := h.tasks.CreateList(r.Context(), &pbtasks.CreateListRequest{Name: body.Name})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, resp.List)
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp, err := h.tasks.GetList(r.Context(), &pbtasks.GetListRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.List)
}

func (h *Handler) ListLists(w http.ResponseWriter, r *http.Request) {
	resp, err := h.tasks.ListLists(r.Context(), &pbtasks.ListListsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Lists)
}

func (h *Handler) UpdateList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Name string `json:"name"`
	}
	if err := decode(r, &body); err != nil || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	resp, err := h.tasks.UpdateList(r.Context(), &pbtasks.UpdateListRequest{Id: id, Name: body.Name})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.List)
}

func (h *Handler) DeleteList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.tasks.DeleteList(r.Context(), &pbtasks.DeleteListRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

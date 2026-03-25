package handler

import (
	"encoding/json"
	"net/http"

	pbacl "github.com/its-me-mayday/sisyphus/bff/proto/acl"
	pbtasks "github.com/its-me-mayday/sisyphus/bff/proto/tasks"
)

type Handler struct {
	tasks pbtasks.TasksServiceClient
	acl   pbacl.AclServiceClient
}

func New(tasks pbtasks.TasksServiceClient, acl pbacl.AclServiceClient) *Handler {
	return &Handler{tasks: tasks, acl: acl}
}

// ─── Helpers ──────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

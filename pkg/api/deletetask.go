package api

import (
	"net/http"

	"final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

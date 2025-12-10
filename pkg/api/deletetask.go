package api

import (
	"log"
	"net/http"

	"final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		log.Printf("failed to delete task: %v", err)
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

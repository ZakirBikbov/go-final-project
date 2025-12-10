package api

import (
	"log"
	"net/http"
	"time"

	"final-project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("failed to delete completed task: %v", err)
			writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
			return
		}

		if err := db.UpdateDate(next, id); err != nil {
			log.Printf("failed to update task date: %v", err)
			writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
			return
		}
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"final-project/pkg/db"
)

type updateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	if req.ID == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Неверный формат идентификатора"})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан заголовок задачи"})
		return
	}

	task := &db.Task{
		ID:      id,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := checkDate(task); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	if err := db.UpdateTask(task); err != nil {
		log.Printf("failed to update task: %v", err)
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

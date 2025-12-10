package api

import (
	"encoding/json"
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
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	if req.ID == "" {
		writeJSON(w, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		writeJSON(w, taskResponse{Error: "Неверный формат идентификатора"})
		return
	}

	if req.Title == "" {
		writeJSON(w, taskResponse{Error: "Не указан заголовок задачи"})
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
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	if err := db.UpdateTask(task); err != nil {
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

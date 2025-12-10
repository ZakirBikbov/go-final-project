package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"final-project/pkg/db"
)

type taskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type taskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = nowDate.Format(DateFormat)
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if t.Before(nowDate) {
		if task.Repeat != "" {
			next, err := NextDate(nowDate, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		} else {
			task.Date = nowDate.Format(DateFormat)
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан заголовок задачи"})
		return
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := checkDate(task); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		log.Printf("failed to add task: %v", err)
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, taskResponse{ID: fmt.Sprintf("%d", id)})
}

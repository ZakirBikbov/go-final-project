package api

import (
	"encoding/json"
	"fmt"
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

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
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
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	if req.Title == "" {
		writeJSON(w, taskResponse{Error: "Не указан заголовок задачи"})
		return
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := checkDate(task); err != nil {
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	writeJSON(w, taskResponse{ID: fmt.Sprintf("%d", id)})
}

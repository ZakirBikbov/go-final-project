package api

import (
	"net/http"
	"strconv"

	"final-project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, taskResponse{Error: err.Error()})
		return
	}

	writeJSON(w, TaskResp{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
}

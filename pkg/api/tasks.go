package api

import (
	"fmt"
	"log"
	"net/http"

	"final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*TaskResp `json:"tasks"`
}

type TaskResp struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	limit := 50

	tasks, err := db.Tasks(limit, search)
	if err != nil {
		log.Printf("failed to get tasks: %v", err)
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}

	taskResps := make([]*TaskResp, 0, len(tasks))
	for _, task := range tasks {
		taskResps = append(taskResps, &TaskResp{
			ID:      fmt.Sprintf("%d", task.ID),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: taskResps})
}

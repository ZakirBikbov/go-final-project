package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/signin", signinHandler)
	mux.HandleFunc("/api/task", authMiddleware(taskHandler))
	mux.HandleFunc("/api/task/done", authMiddleware(doneTaskHandler))
	mux.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
}

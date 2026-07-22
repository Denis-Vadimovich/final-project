package api

import "net/http"

func Init(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/task/done", doneHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/Denis-Vadimovich/final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		writeJsonError(w, err)
		return
	}
	writeJsonTasks(w, TasksResp{
		Tasks: tasks,
	})
}

func writeJsonTasks(w http.ResponseWriter, task TasksResp) {

	res, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Denis-Vadimovich/final-project/pkg/db"
)

type Response struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func checkDate(task *db.Task) error {
	nowTime := time.Now().Format(dateFormat)
	now, err := time.Parse(dateFormat, nowTime)
	if err != nil {
		return err
	}

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	//проверяем, что в task.Date указана корректная дата
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	//если определён task.Repeat, то проверяем корректность правила и заодно получаем следующую дату
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)

		if err != nil {
			return err
		}
	}

	// если сегодня (now) больше task.Date (t)

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(dateFormat)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJsonError(w, err)
		return
	}

	if task.Title == "" {
		writeJsonError(w, fmt.Errorf("error - Title is empty"))
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	writeJsonSuccess(w, id)

}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJsonError(w, err)
		return
	}

	if task.Title == "" {
		writeJsonError(w, fmt.Errorf("error - Title is empty"))
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJsonError(w, err)
		return
	}
	id, err := strconv.Atoi(task.ID)

	writeJsonSuccess(w, int64(id))

}

func writeJsonSuccess(w http.ResponseWriter, id int64) {

	response := Response{ID: id}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func writeJsonError(w http.ResponseWriter, mes error) {
	response := Response{Error: mes.Error()}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

func writeJson(w http.ResponseWriter, task db.Task) {

	res, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func writeJsonEmpty(w http.ResponseWriter) {

	type emptyStruct struct{}
	res := emptyStruct{}

	emptyJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(emptyJson)

}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	task, err := db.GetTask(r.URL.Query().Get("id"))
	if err != nil {
		writeJsonError(w, err)
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(task.ID)
		if err != nil {
			writeJsonError(w, err)
			return
		}
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJsonError(w, err)
		}

		err = db.UpdateDate(nextDate, task.ID)
		if err != nil {
			writeJsonError(w, err)
		}
	}
	writeJsonEmpty(w)

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		task, err := db.GetTask(r.URL.Query().Get("id"))
		if err != nil {
			writeJsonError(w, err)
			return
		}
		writeJson(w, task)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		err := db.DeleteTask(r.URL.Query().Get("id"))
		if err != nil {
			writeJsonError(w, err)
			return
		}
		writeJsonEmpty(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

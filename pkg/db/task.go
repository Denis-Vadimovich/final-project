package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, "repeat") VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {

	rows, err := db.Query(`SELECT id, date, title,comment, repeat FROM scheduler ORDER BY date LIMIT ? `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*Task
	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}
		res = append(res, task)
	}

	if res == nil {
		res = make([]*Task, 0)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func GetTask(id string) (Task, error) {
	row := db.QueryRow(`SELECT id, date, title,comment, repeat FROM scheduler WHERE id = ? `, id)

	task := Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return task, fmt.Errorf("task not found")
	}
	if err != nil {
		return task, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {

	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil

}

func UpdateDate(next string, id string) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil

}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

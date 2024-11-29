package tasks

import (
	"database/sql"
	"go_final_project/database"
	"go_final_project/str"
	"net/http"
	"time"
)

func DoneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
			return
		}

		task, err := database.GetTaskByID(db, taskID)
		if err != nil {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		if task.Repeat == "" {
			err = database.DeleteTaskFromDB(db, taskID)
			if err != nil {
				http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
				return
			}
		} else {
			taskDate, err := time.Parse(str.DateFormat, task.Date)
			if err != nil {
				http.Error(w, `{"error":"Неверный формат даты задачи"}`, http.StatusInternalServerError)
				return
			}

			nextDate, err := NextDate(taskDate, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Ошибка расчёта следующей даты"}`, http.StatusInternalServerError)
				return
			}
			task.Date = nextDate

			err = database.UpdateTaskInDB(db, task)
			if err != nil {
				http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
		}

		w.Write([]byte("{}"))
	}
}

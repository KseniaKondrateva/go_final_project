package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/str"
	"io"
	"net/http"
	"time"
)

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		var task str.Task
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error":"Ошибка чтения тела запроса"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка преобразования формата"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()
		if task.Date == "" {
			task.Date = now.Format(str.DateFormat)
		}

		parsedDate, err := time.Parse(str.DateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
			return
		}

		parsedDate = normalizeToDate(parsedDate)
		now = normalizeToDate(now)

		if parsedDate.Before(now) && task.Repeat == "" {
			task.Date = now.Format(str.DateFormat)
		}

		if parsedDate.Before(now) && task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка в правиле повторения: %s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}

		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка добавления задачи: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка получения ID задачи: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

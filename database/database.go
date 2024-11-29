package database

import (
	"database/sql"
	"errors"
	"go_final_project/str"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() *sql.DB {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		appPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Не удалось получить текущий рабочий каталог: %v", err)
		}
		dbFile = filepath.Join(appPath, "scheduler.db")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Ошибка открытия базы данных: %v", err)
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		createTable := `
        CREATE TABLE IF NOT EXISTS scheduler (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          date CHAR(8) NOT NULL,
          title TEXT NOT NULL,
          comment TEXT,
          repeat VARCHAR(128)
        );
        CREATE INDEX IF NOT EXISTS indexDate ON scheduler(date);
        `
		if _, err := db.Exec(createTable); err != nil {
			log.Fatalf("Ошибка создания таблицы: %v", err)
		}
		log.Println("База данных создана")
	} else {
		log.Println("База данных уже существует")
	}

	return db
}

func DeleteTaskFromDB(db *sql.DB, id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func GetTaskByID(db *sql.DB, id string) (str.Task, error) {
	var task str.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return str.Task{}, errors.New("задача не найдена")
	}
	return task, err
}

func UpdateTaskInDB(db *sql.DB, task str.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

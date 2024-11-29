package main

import (
	"go_final_project/database"
	"go_final_project/tasks"
	"log"
	"net/http"
)

func main() {
	db := database.CreateDB()
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", tasks.HandleNextDate)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			tasks.AddTaskHandler(db)(w, r)
		case http.MethodPut:
			tasks.EditTaskHandler(db)(w, r)
		case http.MethodGet:
			tasks.GetTaskHandler(db)(w, r)
		case http.MethodDelete:
			tasks.DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/tasks", TasksHandler(db))
	http.HandleFunc("/api/task/done", tasks.DoneTaskHandler(db))

	port := ":7540"
	log.Println("Сервер запущен на порту", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

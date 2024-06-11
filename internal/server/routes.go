package server

import (
	"api/internal/model"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/", s.HomeHandler)
	r.Post("/", s.CreateTaskHandler)
	r.Put("/update/{id}", s.UpdateTaskHandler)
	r.Delete("/delete/{id}", s.DeleteTaskHandler)
	r.Get("/get-task-row/{id}", s.GetTaskRowHandler)
	r.Get("/get-edit-form/{id}", s.GetEditTaskFormHandler)

	return r
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./internal/template/home.html", "./internal/template/fragments/task-row.html")
	check(err)

	items, _ := s.Db.Tasks()

	data := struct {
		Title string
		Items []model.Task
	}{
		Title: "Task List",
		Items: items,
	}

	err = t.Execute(w, data)
	check(err)
}

func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	data := model.Task{Title: r.FormValue("title"), Description: r.FormValue("description")}
	task, err := s.Db.AddTask(data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("./internal/template/fragments/task-row.html")
	_ = t.Execute(w, task)
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idParam, 10, 32)

	err := s.Db.DeleteTask(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte{})
}

func (s *Server) GetEditTaskFormHandler(w http.ResponseWriter, r *http.Request) {
	task, err := s.getTaskFromRequest(r)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("./internal/template/fragments/task-edit-form.html")
	_ = t.Execute(w, task)
}

func (s *Server) GetTaskRowHandler(w http.ResponseWriter, r *http.Request) {
	task, err := s.getTaskFromRequest(r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("./internal/template/fragments/task-row.html")
	_ = t.Execute(w, task)
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := s.getTaskFromRequest(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	completed := r.FormValue("completed")
	if completed == "on" {
		now := time.Now()
		task.CompletedAt = &now
	} else {
		task.CompletedAt = nil
	}

	title := r.FormValue("title")
	if len(title) > 0 {
		task.Title = title
	}

	description := r.FormValue("description")
	if len(description) > 0 {
		task.Description = description
	}

	_ = s.Db.UpdateTask(task)

	t, _ := template.ParseFiles("./internal/template/fragments/task-row.html")
	_ = t.Execute(w, task)
}

func (s *Server) getTaskFromRequest(r *http.Request) (model.Task, error) {
	idParam := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idParam, 10, 32)

	task, err := s.Db.GetTask(id)
	return task, err
}

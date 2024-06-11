package database

import (
	"api/internal/model"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

type Service interface {
	Close() error
	Tasks() ([]model.Task, error)
	AddTask(task model.Task) (model.Task, error)
	UpdateTask(task model.Task) error
	DeleteTask(id int64) error
	GetTask(id int64) (model.Task, error)
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Run db migrations...
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", database, driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) Tasks() ([]model.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.CompletedAt); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *service) AddTask(task model.Task) (model.Task, error) {
	err := s.db.QueryRow("INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id", task.Title, task.Description).Scan(&task.Id)
	if err != nil {
		return task, fmt.Errorf("AddTask: %v", err)
	}

	return task, nil
}

func (s *service) UpdateTask(task model.Task) error {
	err := s.db.QueryRow("UPDATE tasks SET title = $1, description = $2, completedat = $3 WHERE id = $4", task.Title, task.Description, task.CompletedAt, task.Id)
	if err != nil {
		return fmt.Errorf("UpdateTask: %v", err)
	}

	return nil
}

func (s *service) DeleteTask(id int64) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("DeleteTask: %v", err)
	}
	return nil
}

func (s *service) GetTask(id int64) (model.Task, error) {
	var task model.Task
	err := s.db.QueryRow("SELECT * FROM tasks WHERE id = $1", id).Scan(&task.Id, &task.Title, &task.Description, &task.CompletedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// No rows were returned
			return task, fmt.Errorf("No task found with id: %v", id)
		}
		return task, err
	}

	return task, nil
}

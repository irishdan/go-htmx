package tests

import (
	"api/internal/model"
	"api/internal/server"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDB struct{}

func (mdb *MockDB) Close() error {
	return nil
}

func (mdb *MockDB) Tasks() ([]model.Task, error) {
	return []model.Task{{Id: 1, Title: "Test Task", Description: "Test Description"}}, nil
}

func (mdb *MockDB) AddTask(task model.Task) (model.Task, error) {
	return task, nil
}

func (mdb *MockDB) UpdateTask(task model.Task) error {
	return nil
}

func (mdb *MockDB) DeleteTask(id int64) error {
	return nil
}

func (mdb *MockDB) GetTask(id int64) (model.Task, error) {
	return model.Task{Id: 1, Title: "Test Task", Description: "Test Description"}, nil
}

func TestHomeHandler(t *testing.T) {
	s := &server.Server{Db: &MockDB{}}
	app := httptest.NewServer(http.HandlerFunc(s.HomeHandler))

	resp, err := http.Get(app.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := homepage
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

var homepage = `<title>Task List</title>
<tbody>

<tr id='1'>
<td>Test Task</td>
<td>Test Description</td>

</tbody>`

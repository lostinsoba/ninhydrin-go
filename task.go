package ninhydrin

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
)

const (
	TaskStatusIdle       TaskStatus = "idle"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusTimeout    TaskStatus = "timeout"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusDone       TaskStatus = "done"
)

type TaskStatus string

type Task struct {
	ID          string     `json:"id"`
	NamespaceID string     `json:"namespace_id"`
	Timeout     int64      `json:"timeout,omitempty"`
	RetriesLeft int        `json:"retries_left,omitempty"`
	Status      TaskStatus `json:"status,omitempty"`
}

// taskService provides tasks CRUD
type taskService struct {
	controller *controller
}

// newTaskService returns task service
func newTaskService(client *Client, endpoint string) *taskService {
	return &taskService{
		controller: &controller{
			client:   client,
			endpoint: endpoint,
		},
	}
}

// List returns list of tasks
func (ts *taskService) List(ctx context.Context, namespaceID string) ([]*Task, error) {
	query := []queryArg{
		{
			k: "namespace_id",
			v: namespaceID,
		},
	}
	resp, err := ts.controller.GETWithQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	var data taskListData
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data.List, err
}

// Read returns task by given ID
func (ts *taskService) Read(ctx context.Context, taskID string) (*Task, error) {
	resp, err := ts.controller.GET(ctx, taskID)
	if err != nil {
		return nil, err
	}
	var task *Task
	err = json.Unmarshal(resp, &task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete removes task by given ID
func (ts *taskService) Delete(ctx context.Context, taskID string) error {
	_, err := ts.controller.DELETE(ctx, taskID)
	return err
}

// Register registers task
func (ts *taskService) Register(ctx context.Context, task *Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ts.controller.POST(ctx, reader)
	return err
}

// Capture returns list of tasks
func (ts *taskService) Capture(ctx context.Context, namespaceID string, limit int) ([]*Task, error) {
	query := []queryArg{
		{
			k: "namespace_id",
			v: namespaceID,
		},
		{
			k: "limit",
			v: strconv.Itoa(limit),
		},
	}
	resp, err := ts.controller.GETWithQuery(ctx, query, "capture")
	if err != nil {
		return nil, err
	}
	var data taskListData
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data.List, nil
}

// Release updates tasks with given status
func (ts *taskService) Release(ctx context.Context, status TaskStatus, taskIDs []string) error {
	data := &releaseData{
		Status:  status,
		TaskIDs: taskIDs,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ts.controller.PUT(ctx, reader, "release")
	return err
}

type releaseData struct {
	Status  TaskStatus `json:"status"`
	TaskIDs []string   `json:"task_ids"`
}

type taskListData struct {
	List []*Task `json:"list"`
}

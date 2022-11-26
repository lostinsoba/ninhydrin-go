package ninhydrin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type Task struct {
	ID          string `json:"id"`
	PoolID      string `json:"pool_id"`
	Timeout     int64  `json:"timeout,omitempty"`
	RetriesLeft int    `json:"retries_left,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
	Status      string `json:"status,omitempty"`
}

type TaskListData struct {
	List []*Task `json:"list"`
}

type TaskStatusUpdateData struct {
	Status string `json:"status"`
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
func (ts *taskService) List(ctx context.Context) ([]*Task, error) {
	resp, err := ts.controller.GET(ctx)
	if err != nil {
		return nil, err
	}
	var result *TaskListData
	err = json.Unmarshal(resp, &result)
	return result.List, err
}

// Register registers task
func (ts *taskService) Register(ctx context.Context, task *Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ts.controller.POST(ctx, reader, "register")
	return err
}

// Capture returns list of tasks
func (ts *taskService) Capture(ctx context.Context, limit int) ([]*Task, error) {
	queryParams := url.Values{}
	queryParams.Set("limit", strconv.Itoa(limit))
	resp, err := ts.controller.GET(ctx, queryParams.Encode())
	if err != nil {
		return nil, err
	}
	var result *TaskListData
	err = json.Unmarshal(resp, &result)
	return result.List, err
}

// Update updates task status
func (ts *taskService) UpdateStatus(ctx context.Context, taskID string, status string) error {
	body, err := json.Marshal(&TaskStatusUpdateData{Status: status})
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ts.controller.PUT(ctx, reader, taskID, "status")
	return err
}

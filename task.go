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

// ListIDs returns list of task IDs
func (ts *taskService) ListIDs(ctx context.Context) (ids []string, err error) {
	resp, err := ts.controller.GET(ctx)
	if err != nil {
		return
	}
	var data idListData
	err = json.Unmarshal(resp, &data)
	ids = data.List
	return
}

// Read returns task by given ID
func (ts *taskService) Read(ctx context.Context, taskID string) (task *Task, err error) {
	resp, err := ts.controller.GET(ctx, taskID)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &task)
	return
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

// CaptureIDs returns list of task IDs
func (ts *taskService) CaptureIDs(ctx context.Context, limit int) (ids []string, err error) {
	query := []queryArg{
		{
			k: "limit",
			v: strconv.Itoa(limit),
		},
	}
	resp, err := ts.controller.GETWithQuery(ctx, query, "capture")
	if err != nil {
		return
	}
	var data idListData
	err = json.Unmarshal(resp, &data)
	ids = data.List
	return
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

type idListData struct {
	List []string `json:"list"`
}

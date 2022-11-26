package ninhydrin

import (
	"bytes"
	"context"
	"encoding/json"
)

type Worker struct {
	ID     string   `json:"id"`
	TagIDs []string `json:"tag_ids"`
}

type WorkerListData struct {
	List []*Worker `json:"list"`
}

// workerService provides workers CRUD
type workerService struct {
	controller *controller
}

// newWorkerService returns worker service
func newWorkerService(client *Client, endpoint string) *workerService {
	return &workerService{
		controller: &controller{
			client:   client,
			endpoint: endpoint,
		},
	}
}

// List returns list of workers
func (ws *workerService) List(ctx context.Context) ([]*Worker, error) {
	resp, err := ws.controller.GET(ctx)
	if err != nil {
		return nil, err
	}
	var result WorkerListData
	err = json.Unmarshal(resp, &result)
	return result.List, err
}

// Register registers worker
func (ws *workerService) Register(ctx context.Context, worker *Worker) error {
	body, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ws.controller.POST(ctx, reader, "register")
	return err
}

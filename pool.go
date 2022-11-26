package ninhydrin

import (
	"bytes"
	"context"
	"encoding/json"
)

type Pool struct {
	ID          string   `json:"id"`
	Description string   `json:"description,omitempty"`
	TagIDs      []string `json:"tag_ids"`
}

type PoolListData struct {
	List []*Pool `json:"list"`
}

// poolService provides pools CRUD
type poolService struct {
	controller *controller
}

// newPoolService returns pool service
func newPoolService(client *Client, endpoint string) *poolService {
	return &poolService{
		controller: &controller{
			client:   client,
			endpoint: endpoint,
		},
	}
}

// List returns list of pools
func (ps *poolService) List(ctx context.Context) ([]*Pool, error) {
	resp, err := ps.controller.GET(ctx)
	if err != nil {
		return nil, err
	}
	var result *PoolListData
	err = json.Unmarshal(resp, &result)
	return result.List, err
}

// Register registers pool
func (ps *poolService) Register(ctx context.Context, pool *Pool) error {
	body, err := json.Marshal(pool)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ps.controller.POST(ctx, reader, "register")
	return err
}

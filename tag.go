package ninhydrin

import (
	"bytes"
	"context"
	"encoding/json"
)

type TagData struct {
	ID string `json:"id"`
}

type TagListData struct {
	List []string `json:"list"`
}

// tagService provides tags CRUD
type tagService struct {
	controller *controller
}

// newTagService returns tag service
func newTagService(client *Client, endpoint string) *tagService {
	return &tagService{
		controller: &controller{
			client:   client,
			endpoint: endpoint,
		},
	}
}

// List returns list of tags
func (ts *tagService) List(ctx context.Context) ([]string, error) {
	resp, err := ts.controller.GET(ctx)
	if err != nil {
		return nil, err
	}
	var result *TagListData
	err = json.Unmarshal(resp, &result)
	return result.List, err
}

// Register registers tag
func (ts *tagService) Register(ctx context.Context, id string) error {
	tag := TagData{ID: id}
	body, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ts.controller.POST(ctx, reader, "register")
	return err
}

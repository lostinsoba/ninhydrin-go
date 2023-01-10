package ninhydrin

import (
	"bytes"
	"context"
	"encoding/json"
)

type Namespace struct {
	ID string `json:"id"`
}

// namespaceService provides namespace CRUD
type namespaceService struct {
	controller *controller
}

// newNamespaceService returns namespace service
func newNamespaceService(client *Client, endpoint string) *namespaceService {
	return &namespaceService{
		controller: &controller{
			client:   client,
			endpoint: endpoint,
		},
	}
}

// List returns list of namespaces
func (ns *namespaceService) List(ctx context.Context) ([]*Namespace, error) {
	resp, err := ns.controller.GET(ctx)
	if err != nil {
		return nil, err
	}
	var data namespaceListData
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data.List, nil
}

// Read returns namespace by given ID
func (ns *namespaceService) Read(ctx context.Context, namespaceID string) (*Namespace, error) {
	resp, err := ns.controller.GET(ctx, namespaceID)
	if err != nil {
		return nil, err
	}
	var namespace *Namespace
	err = json.Unmarshal(resp, &namespace)
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

// Delete removes namespace by given ID
func (ns *namespaceService) Delete(ctx context.Context, namespaceID string) error {
	_, err := ns.controller.DELETE(ctx, namespaceID)
	return err
}

// Register registers namespace
func (ns *namespaceService) Register(ctx context.Context, namespace *Namespace) error {
	body, err := json.Marshal(namespace)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	_, err = ns.controller.POST(ctx, reader)
	return err
}

type namespaceListData struct {
	List []*Namespace `json:"list"`
}

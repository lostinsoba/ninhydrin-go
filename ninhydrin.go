package ninhydrin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// defaultBaseURL is a default URL of Ninhydrin API
const defaultBaseURL = "http://ninhydrin:8080/v1"

// workerIDHeader indicates that client is being used for worker
const workerIDHeader = "X-Ninhydrin-Worker-ID"

// getDefaultHeaders returns collection of default headers
func getDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":   "Ninhydrin Go API Client",
		"Content-Type": "application/json",
	}
}

// getDefaultHTTPClient returns default HTTP Client
func getDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}
}

// New returns new Ninhydrin API client
func New(options ...Option) *Client {
	client := &Client{
		baseURL:    defaultBaseURL,
		headers:    getDefaultHeaders(),
		httpClient: getDefaultHTTPClient(),
	}

	for _, opt := range options {
		opt(client)
	}

	client.Tag = newTagService(client, "/tag")
	client.Pool = newPoolService(client, "/pool")
	client.Task = newTaskService(client, "/task")
	client.Worker = newWorkerService(client, "/worker")

	return client
}

// Client is Ninhydrin API client
type Client struct {
	baseURL    string
	headers    map[string]string
	httpClient *http.Client

	Tag    *tagService
	Pool   *poolService
	Task   *taskService
	Worker *workerService
}

// Option defines an option for a Client
type Option func(*Client)

// OptionBaseURL rewrites the API URL value
func OptionBaseURL(url string) func(*Client) {
	return func(c *Client) { c.baseURL = url }
}

// OptionUserAgent rewrites "User-Agent" header value
func OptionUserAgent(userAgent string) func(*Client) {
	return func(c *Client) { c.headers["User-Agent"] = userAgent }
}

// OptionHTTPClient sets another http.Client
func OptionHTTPClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) { c.httpClient = httpClient }
}

// OptionWorkerID sets worker ID header
func OptionWorkerID(workerID string) func(*Client) {
	return func(c *Client) { c.headers[workerIDHeader] = workerID }
}

// controller is a controller that provides basic http methods for every service
// use it as an extension point to declare endpoint-specific behaviour
type controller struct {
	client   *Client
	endpoint string
}

// GET performs GET request
// GET(ctx, "12345") will perform GET "/serviceEndpoint/12345" request
func (ctrl *controller) GET(ctx context.Context, entries ...string) ([]byte, error) {
	path := ctrl.pathJoin(entries...)
	return ctrl.processRequest(ctx, http.MethodGet, path, nil)
}

// DELETE performs DELETE request
// DELETE(ctx, "12345") will perform DELETE "/serviceEndpoint/12345" request
func (ctrl *controller) DELETE(ctx context.Context, entries ...string) ([]byte, error) {
	path := ctrl.pathJoin(entries...)
	return ctrl.processRequest(ctx, http.MethodDelete, path, nil)
}

// POST performs POST request
// POST(ctx, payload) will perform POST "/serviceEndpoint" request with payload
func (ctrl *controller) POST(ctx context.Context, body io.Reader, entries ...string) ([]byte, error) {
	path := ctrl.pathJoin(entries...)
	return ctrl.processRequest(ctx, http.MethodPost, path, body)
}

// PATCH performs PATCH request
// PATCH(ctx, payload, "12345", "task") will perform PATCH "/serviceEndpoint/12345/task" request with payload
func (ctrl *controller) PATCH(ctx context.Context, body io.Reader, entries ...string) ([]byte, error) {
	path := ctrl.pathJoin(entries...)
	return ctrl.processRequest(ctx, http.MethodPatch, path, body)
}

// PUT performs PUT request
// PUT(ctx, payload, "12345") will perform PUT "/serviceEndpoint/12345" request with payload
func (ctrl *controller) PUT(ctx context.Context, body io.Reader, entries ...string) ([]byte, error) {
	path := ctrl.pathJoin(entries...)
	return ctrl.processRequest(ctx, http.MethodPut, path, body)
}

func (ctrl *controller) pathJoin(entries ...string) string {
	var res = ctrl.client.baseURL + ctrl.endpoint
	for _, entry := range entries {
		if !strings.HasSuffix(res, "/") {
			res += "/"
		}
		res += url.PathEscape(entry)
	}
	return res
}

func (ctrl *controller) processRequest(ctx context.Context, method string, url string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("invalid request parameters: %w", err)
	}
	for k, v := range ctrl.client.headers {
		request.Header.Set(k, v)
	}
	response, err := ctrl.client.httpClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if !isValidStatusCode(response.StatusCode) {
		if message := extractErrorMessage(data); message != "" {
			return nil, fmt.Errorf("invalid status: %s, message: %s", response.Status, message)
		}
		return nil, fmt.Errorf("invalid status: %s", response.Status)
	}
	return data, nil
}

func isValidStatusCode(statusCode int) bool {
	return (statusCode >= 200) && (statusCode <= 299)
}

func extractErrorMessage(data []byte) string {
	var res struct {
		Value string `json:"message"`
	}
	_ = json.Unmarshal(data, &res)
	return res.Value
}

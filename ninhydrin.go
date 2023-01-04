package ninhydrin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// defaultBaseURL is a default URL of Ninhydrin API
const defaultBaseURL = "http://localhost:8080/v1"

// getDefaultHeaders returns collection of default headers
func getDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":   "Ninhydrin Go API Client",
		"Content-Type": "application/json; charset=utf-8",
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

	client.Task = newTaskService(client, "/task")

	return client
}

// Client is Ninhydrin API client
type Client struct {
	baseURL    string
	headers    map[string]string
	httpClient *http.Client

	Task *taskService
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

// OptionHeader adds or sets existing header value
func OptionHeader(key, value string) func(*Client) {
	return func(c *Client) { c.headers[key] = value }
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

type queryArg struct {
	k, v string
}

// GETWithQuery performs GET request with query
// GETWithQuery(ctx, query, "12345") will perform GET "/serviceEndpoint/12345?=query" request
func (ctrl *controller) GETWithQuery(ctx context.Context, query []queryArg, entries ...string) ([]byte, error) {
	path := ctrl.pathJoin(entries...)
	return ctrl.processRequest(ctx, http.MethodGet, path, nil, query...)
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
	var res strings.Builder
	res.WriteString(ctrl.client.baseURL)
	res.WriteString(ctrl.endpoint)
	for _, entry := range entries {
		res.WriteString("/")
		res.WriteString(entry)
	}
	return res.String()
}

func (ctrl *controller) processRequest(ctx context.Context, method string, url string, body io.Reader, query ...queryArg) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("invalid request parameters: %w", err)
	}
	for k, v := range ctrl.client.headers {
		request.Header.Set(k, v)
	}
	if len(query) > 0 {
		q := request.URL.Query()
		for _, arg := range query {
			q.Set(arg.k, arg.v)
		}
		request.URL.RawQuery = q.Encode()
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

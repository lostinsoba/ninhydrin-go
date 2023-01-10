package ninhydrin

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func flushAll(api *Client, ctx context.Context) (err error) {
	namespaces, err := api.Namespace.List(ctx)
	if err != nil {
		return
	}
	for _, namespace := range namespaces {
		var (
			tasks []*Task
		)
		tasks, err = api.Task.List(ctx, namespace.ID)
		if err != nil {
			return
		}
		for _, task := range tasks {
			err = api.Task.Delete(ctx, task.ID)
			if err != nil {
				return
			}
		}
		err = api.Namespace.Delete(ctx, namespace.ID)
		if err != nil {
			return
		}
	}
	return
}

type logger interface {
	Printf(format string, v ...interface{})
}

func newDebugHTTPClient() *http.Client {
	return &http.Client{
		Transport: dbgRoundTripper{
			log.New(os.Stdout, "", log.LstdFlags),
			http.DefaultTransport,
		},
	}
}

type dbgRoundTripper struct {
	log logger
	rt  http.RoundTripper
}

func (drt dbgRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	var (
		requestID = generateRandomID("req")
	)
	switch req.Method {
	case http.MethodPost, http.MethodPut:
		payload := copyFrom(req.Body)
		drt.log.Printf("[%s] request: %s %s %s\n", requestID, req.Method, req.URL, payload)
	default:
		drt.log.Printf("[%s] request: %s %s\n", requestID, req.Method, req.URL)
	}
	res, err = drt.rt.RoundTrip(req)
	if err != nil {
		drt.log.Printf("[%s] request error: %s", requestID, err.Error())
	} else {
		drt.log.Printf("[%s] response: %s", requestID, res.Status)
	}
	return
}

func generateRandomID(prefix string) string {
	id := rand.Intn(100)
	return fmt.Sprintf("%s-%d", prefix, id)
}

func copyFrom(src io.Reader) string {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	_, _ = io.Copy(w, src)
	return b.String()
}

func assertEqual(t *testing.T, name string, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("[%s] expected: %v, actual: %v", name, expected, actual)
	}
}

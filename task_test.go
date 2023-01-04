package ninhydrin

import (
	"context"
	"testing"
)

func TestTaskService(t *testing.T) {
	var (
		ctx = context.Background()
		dbg = newDebugHTTPClient()
		api = New(OptionHTTPClient(dbg))
	)

	task := &Task{
		ID:          generateRandomID("Task"),
		Timeout:     10,
		RetriesLeft: 5,
		Status:      TaskStatusIdle,
	}

	err := flushAll(api, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = api.Task.Register(ctx, task)
	assertEqual(t, "register task", nil, err)

	list, err := api.Task.ListIDs(ctx)
	assertEqual(t, "list task ids returned no error", nil, err)
	assertEqual(t, "list task ids with registered task", []string{task.ID}, list)

	actualTask, err := api.Task.Read(ctx, task.ID)
	assertEqual(t, "read task by id returned no error", nil, err)
	assertEqual(t, "read task by id", task, actualTask)

	tasks, err := api.Task.CaptureIDs(ctx, 10)
	assertEqual(t, "capture tasks ids returned no error", nil, err)
	assertEqual(t, "capture tasks ids", []string{task.ID}, tasks)

	err = api.Task.Release(ctx, TaskStatusDone, []string{task.ID})
	assertEqual(t, "release task", nil, err)

	err = api.Task.Delete(ctx, task.ID)
	assertEqual(t, "delete task", nil, err)

	list, err = api.Task.ListIDs(ctx)
	assertEqual(t, "list task ids without deleted task returned no error", nil, err)
	assertEqual(t, "list task ids without deleted task", []string{}, list)
}

func flushAll(api *Client, ctx context.Context) error {
	var (
		ids []string
		err error
	)
	ids, err = api.Task.ListIDs(ctx)
	if err != nil {
		return err
	}
	for _, id := range ids {
		err = api.Task.Delete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

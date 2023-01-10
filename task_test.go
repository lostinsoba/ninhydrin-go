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

	namespace := &Namespace{
		ID: generateRandomID("Namespace"),
	}

	task := &Task{
		ID:          generateRandomID("Task"),
		NamespaceID: namespace.ID,
		Timeout:     10,
		RetriesLeft: 5,
		Status:      TaskStatusIdle,
	}

	taskCaptured := &Task{
		ID:          task.ID,
		NamespaceID: task.NamespaceID,
		Timeout:     task.Timeout,
		RetriesLeft: task.RetriesLeft - 1,
		Status:      TaskStatusInProgress,
	}

	err := flushAll(api, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = api.Namespace.Register(ctx, namespace)
	assertEqual(t, "register namespace", nil, err)

	err = api.Task.Register(ctx, task)
	assertEqual(t, "register task", nil, err)

	list, err := api.Task.List(ctx, namespace.ID)
	assertEqual(t, "list tasks returned no error", nil, err)
	assertEqual(t, "list tasks with registered task", []*Task{task}, list)

	actualTask, err := api.Task.Read(ctx, task.ID)
	assertEqual(t, "read task by id returned no error", nil, err)
	assertEqual(t, "read task by id", task, actualTask)

	tasks, err := api.Task.Capture(ctx, namespace.ID, 10)
	assertEqual(t, "capture tasks returned no error", nil, err)
	assertEqual(t, "capture tasks", []*Task{taskCaptured}, tasks)

	err = api.Task.Release(ctx, TaskStatusDone, []string{task.ID})
	assertEqual(t, "release task", nil, err)

	err = api.Task.Delete(ctx, task.ID)
	assertEqual(t, "delete task", nil, err)

	list, err = api.Task.List(ctx, namespace.ID)
	assertEqual(t, "list tasks without deleted task returned no error", nil, err)
	assertEqual(t, "list tasks without deleted task", []*Task{}, list)
}

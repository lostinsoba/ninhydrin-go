package ninhydrin

import (
	"context"
	"testing"
)

func TestNamespaceService(t *testing.T) {
	var (
		ctx = context.Background()
		dbg = newDebugHTTPClient()
		api = New(OptionHTTPClient(dbg))
	)

	namespace := &Namespace{
		ID: generateRandomID("Namespace"),
	}

	err := flushAll(api, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = api.Namespace.Register(ctx, namespace)
	assertEqual(t, "register namespace", nil, err)

	list, err := api.Namespace.List(ctx)
	assertEqual(t, "list namespaces returned no error", nil, err)
	assertEqual(t, "list namespaces with registered namespace", []*Namespace{namespace}, list)

	actualNamespace, err := api.Namespace.Read(ctx, namespace.ID)
	assertEqual(t, "read namespace by id returned no error", nil, err)
	assertEqual(t, "read namespace by id", namespace, actualNamespace)

	err = api.Namespace.Delete(ctx, namespace.ID)
	assertEqual(t, "delete namespace", nil, err)

	list, err = api.Namespace.List(ctx)
	assertEqual(t, "list namespaces without deleted namespace returned no error", nil, err)
	assertEqual(t, "list namespaces without deleted namespace", []*Namespace{}, list)
}

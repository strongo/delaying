package delaying

import (
	"context"
	"testing"
)

func TestNewFunction(t *testing.T) {
	testNewDelayer(t, NewFunction)
}

func TestNewDelayer(t *testing.T) {
	testNewDelayer(t, NewDelayer)
}
func testNewDelayer(t *testing.T, newDelayer func(
	id string,
	implementation any,
	enqueueWork func(c context.Context, params Params, args ...interface{}) error,
	enqueueWorkMulti func(c context.Context, params Params, args ...[]interface{}) error,
) Delayer) {
	t.Run("EnqueueWork", func(t *testing.T) {
		var singleArgs []any
		enqueueWork := func(c context.Context, params Params, args ...any) error {
			singleArgs = args
			return nil
		}
		enqueueWorkMulti := func(c context.Context, params Params, args ...[]any) error {
			panic("unexpected call")
		}
		f := NewDelayer("EnqueueWorkTest", func() {}, enqueueWork, enqueueWorkMulti)
		if f == nil {
			t.Fatal("f is nil")
		}
		if singleArgs != nil {
			t.Fatal("singleArgs is not nil")
		}
		err := f.EnqueueWork(context.Background(), With("queue1", "path1", 0), 1, 2, 3)
		if err != nil {
			t.Fatal(err)
		}
		if singleArgs == nil {
			t.Fatal("singleArgs is nil")
		}
		if len(singleArgs) != 3 {
			t.Fatal("singleArgs is not 3")
		}
		for i, v := range singleArgs {
			if v != i+1 {
				t.Fatalf("singleArgs[%d] is not %d", i, i+1)
			}
		}

	})
	t.Run("EnqueueWorkMulti", func(t *testing.T) {
		var multiArgs [][]any
		enqueueWork := func(c context.Context, params Params, args ...any) error {
			panic("unexpected call")
		}
		enqueueWorkMulti := func(c context.Context, params Params, args ...[]any) error {
			multiArgs = args
			return nil
		}
		f := newDelayer("EnqueueWorkMultiTest", func() {}, enqueueWork, enqueueWorkMulti)
		err := f.EnqueueWorkMulti(context.Background(), With("queue1", "path1", 0), []any{1, 2}, []any{3, 4})
		if err != nil {
			t.Fatal(err)
		}
		if multiArgs == nil {
			t.Fatal("multiArgs is nil")
		}
		if len(multiArgs) != 2 {
			t.Fatal("len(multiArgs) is not 2")
		}
		if len(multiArgs[0]) != 2 {
			t.Fatal("len(multiArgs[0]) is not 2")
		}
		if len(multiArgs[1]) != 2 {
			t.Fatal("len(multiArgs[1]) is not 2")
		}
		if multiArgs[0][0] != 1 {
			t.Errorf("multiArgs[0][0] is not 1")
		}
		if multiArgs[0][1] != 2 {
			t.Errorf("multiArgs[0][1] is not 2")
		}
		if multiArgs[1][0] != 3 {
			t.Errorf("multiArgs[1][0] is not 3")
		}
		if multiArgs[1][1] != 4 {
			t.Errorf("multiArgs[1][1] is not 4")
		}
	})
}

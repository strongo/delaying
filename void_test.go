package delaying

import (
	"context"
	"testing"
)

func TestVoid(t *testing.T) {
	f := VoidWithLog("test", func() {})
	if f == nil {
		t.Fatal("VoidWithLog() returned nil")
	}
	t.Run("EnqueueWork", func(t *testing.T) {
		if err := f.EnqueueWork(context.Background(), params{queue: "queue1"}, 1, 2, 3); err != nil {
			t.Errorf("EnqueueWork() returned unexpcted error: = %v", err)
		}
	})
	t.Run("EnqueueWorkMulti", func(t *testing.T) {
		if err := f.EnqueueWorkMulti(context.Background(), params{queue: "queue1"}, []any{1, 2}, []any{3, 4}); err != nil {
			t.Errorf("EnqueueWorkMulti() returned unexpcted error: = %v", err)
		}
	})
}

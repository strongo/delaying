package delaying

import (
	"testing"
	"time"
)

func TestWith(t *testing.T) {
	type args struct {
		queue string
		path  string
		delay time.Duration
	}
	tests := []struct {
		name         string
		args         args
		want         params
		expectsPanic bool
	}{
		{"full", args{"queue1", "path1", 3}, params{"queue1", "path1", 3}, false},
		{"empty_queue", args{"", "path1", 3}, params{}, true},
		{"empty_path", args{"queue1", "", 3}, params{}, true},
		{"negative_delay", args{"queue1", "path1", -1}, params{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectsPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("With() did not panic")
					}
				}()
			}
			p := With(tt.args.queue, tt.args.path, tt.args.delay)
			if !tt.expectsPanic {
				if p != tt.want {
					t.Errorf("With() = %v, want %v", p, tt.want)
				}
			}
		})
	}
}

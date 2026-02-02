package fn_test

import (
	"test_api/fn"
	"testing"
)

func TestRandString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		flag   string
		length int
		want   string
	}{
		{
			name:   "test",
			flag:   "1",
			length: 10,
			want:   "1234567890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fn.RandString(tt.flag, tt.length)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("RandString() = %v, want %v", got, tt.want)
			}
		})
	}
}

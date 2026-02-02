package tool_test

import (
	"test_api/tool"
	"testing"
)

func TestEvaluateWithGoval(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		expr    string
		want    any
		wantErr bool
	}{
		{
			name:    "test1",
			expr:    "(16*32/(4*5)>=25)==true",
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tool.EvaluateWithGoval(tt.expr)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("EvaluateWithGoval() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("EvaluateWithGoval() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("EvaluateWithGoval() = %v, want %v", got, tt.want)
			}
		})
	}
}

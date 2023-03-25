package controller

import "testing"

func Test_isCurrentTimeBetween(t *testing.T) {
	type args struct {
		t1 string
		t2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test", args: args{t1: "2023-03-25 12:00:00", t2: "2023-03-25 18:00:00"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := InTimeSpan(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("isCurrentTimeBetween() = %v, want %v err: %v", got, tt.want, err)
			}
		})
	}
}

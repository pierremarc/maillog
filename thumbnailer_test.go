package main

import "testing"

func Test_thumSize(t *testing.T) {
	type args struct {
		w uint
		h uint
		s uint
	}
	tests := []struct {
		name  string
		args  args
		want  uint
		want1 uint
	}{
		{
			"Thumbnail Size",
			args{1600, 800, 800},
			800,
			400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := thumSize(tt.args.w, tt.args.h, tt.args.s)
			if got != tt.want {
				t.Errorf("thumSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("thumSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

package main

import "testing"

func Test_decodeSubject(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"encoded subject",
			args{s: "=?iso-8859-1?Q?Passage_=E9troit_?="},
			"Passage étroit ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeSubject(tt.args.s); got != tt.want {
				t.Errorf("decodeSubject() = `%v`, want `%v`", got, tt.want)
			}
		})
	}
}

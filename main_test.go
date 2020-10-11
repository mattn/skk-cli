package main

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "aiuEo",
			want:  []string{"aiu", "Eo"},
		},
		{
			input: "aiueo",
			want:  []string{"aiueo"},
		},
		{
			input: "AIUEO",
			want:  []string{"AIUEO"},
		},
		{
			input: "AIUeo",
			want:  []string{"AIUeo"},
		},
	}

	for _, test := range tests {
		got := split(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
}

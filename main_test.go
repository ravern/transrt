package main_test

import (
	"reflect"
	"testing"

	main "github.com/ravernkoh/translate-srt"
)

func TestGroupLinesSentences(t *testing.T) {
	tests := []struct {
		lines  []string
		groups [][]string
	}{}

	for i, test := range tests {
		groups := main.GroupLinesSentences(test.lines)
		if !reflect.DeepEqual(groups, test.groups) {
			t.Errorf("Test %d: Expected %v but got %v.", i+1, test.groups, groups)
		}
	}
}

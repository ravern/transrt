package main_test

import (
	"reflect"
	"testing"

	main "github.com/ravernkoh/translate-srt"
)

func TestSplitAfter(t *testing.T) {
	tests := []struct {
		s      string
		delims []byte
		res    []string
	}{
		{
			"hello. world, today",
			[]byte{'.', ','},
			[]string{"hello.", " world,", " today"},
		},
		{
			"hello. world! now.",
			[]byte{'.', '!', ','},
			[]string{"hello.", " world!", " now.", ""},
		},
	}

	for i, test := range tests {
		res := main.SplitAfter(test.s, test.delims)
		if !reflect.DeepEqual(res, test.res) {
			t.Errorf("Test %d: Expected %v but got %v.", i+1, test.res, res)
		}
	}
}

func TestGroupLinesSentences(t *testing.T) {
	tests := []struct {
		lines  []string
		groups [][]string
	}{
		{
			[]string{
				"hello this",
				"is a human. i am",
				"interested.",
			},
			[][]string{
				[]string{
					"hello this",
					"is a human.",
				},
				[]string{
					" i am",
					"interested.",
				},
				[]string{
					"",
				},
			},
		},
	}

	for i, test := range tests {
		groups := main.GroupLinesSentences(test.lines)
		if !reflect.DeepEqual(groups, test.groups) {
			t.Errorf("Test %d: Expected %v but got %v.", i+1, test.groups, groups)
		}
	}
}

func TestNumWords(t *testing.T) {
	tests := []struct {
		words []string
		num   int
	}{
		{
			[]string{
				"hello this",
				"is a human. i am",
				"interested.",
			},
			8,
		},
		{
			[]string{
				"Another string",
				"to test-my-beloved",
				"function.",
			},
			5,
		},
	}

	for i, test := range tests {
		num := main.NumWords(test.words)
		if num != test.num {
			t.Errorf("Test %d: Expected %d but got %d.", i+1, test.num, num)
		}
	}
}

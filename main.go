package main

import (
	"fmt"
	"os"
	"regexp"

	astisub "github.com/asticode/go-astisub"
)

func main() {
	subs, err := OpenFile("sample.srt")
	check(err)

	lines := ExtractLines(subs)
	fmt.Println(lines)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// OpenFile reads the contents of the file
func OpenFile(path string) (*astisub.Subtitles, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	subs, err := astisub.OpenFile(wd + "/" + path)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

// ExtractLines convert the contents into just pure lines
func ExtractLines(subs *astisub.Subtitles) []string {
	lines := []string{}

	for _, i := range subs.Items {
		for _, l := range i.Lines {
			for _, i2 := range l.Items {
				lines = append(lines, i2.Text)
			}
		}
	}

	return lines
}

// GroupLinesSentences groups the lines up into sentences to be sent to DeepL
func GroupLinesSentences(lines []string) [][]string {
	cur := 0
	groups := [][]string{[]string{}}

	for _, l := range lines {
		sens := splitSentences(l)

		fst := sens[0]
		sens = sens[1:]

		// Take the first element and skip to next one if there isn't any left
		groups[cur] = append(groups[cur], fst)
		if len(sens) == 0 {
			continue
		}

		lst := sens[len(sens)-1]
		sens = sens[:len(sens)-1]

		// Loop through and add sentences
		for _, s := range sens {
			groups = append(groups, []string{s})
		}

		// Add the last one and change the current index
		groups = append(groups, []string{lst})
		cur = len(groups) - 1
	}

	return groups
}

// splitSentences splits a string into sentences but it keeps the seperators
func splitSentences(line string) []string {
	regexp.MustCompile("")
	return nil
}

// TranslateGroup translates a single group
func TranslateGroup(group []string) ([]string, error) {
	return nil, nil
}

// UngroupLines converts the groups back into their ungrouped state
func UngroupLines(groups [][]string) []string {
	return nil
}

// InsertLines inserts the lines back into the subtitles object
func InsertLines(subs *astisub.Subtitles, lines []string) {
}

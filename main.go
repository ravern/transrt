package main

import (
	"fmt"
	"os"
	"strings"

	astisub "github.com/asticode/go-astisub"
	"github.com/ravernkoh/translate-srt/deepl"
)

func main() {
	subs, err := OpenFile("sample.srt")
	check(err)
	lines := ExtractLines(subs)
	groups := GroupLinesSentences(lines)
	groups, err = TranslateGroups(groups)
	check(err)
	lines = UngroupLines(groups)

	for _, l := range lines {
		fmt.Printf("%s\n\n", l)
	}
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

// WriteFile writes the subtitles to the file
func WriteFile(subs *astisub.Subtitles, path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = subs.Write(wd + "/" + path)
	if err != nil {
		return err
	}

	return nil
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
	groups := [][]string{[]string{}}

	for _, l := range lines {
		sens := SplitAfter(l, []byte{'.', '!', '?'})

		fst := sens[0]
		sens = sens[1:]

		// Take the first element and skip to next one if there isn't any left
		cur := len(groups) - 1
		groups[cur] = append(groups[cur], fst)

		// Loop through and add sentences
		for _, s := range sens {
			groups = append(groups, []string{s})
		}
	}

	return groups
}

// SplitAfter is similar the strings.SplitAfter function but with multiple delimeters.
func SplitAfter(s string, delims []byte) []string {
	if len(delims) == 0 {
		return nil
	}

	// Accumulator
	strs := [][]byte{[]byte{}}

	for _, b := range []byte(s) {
		// Add it to the current string
		cur := len(strs) - 1
		strs[cur] = append(strs[cur], b)

		// If it is a delim, move on to next string
		for _, d := range delims {
			if b == d {
				strs = append(strs, []byte{})
				break
			}
		}
	}

	ret := []string{}
	for _, s := range strs {
		ret = append(ret, string(s))
	}

	return ret
}

// TranslateGroups translates a slice of groups
func TranslateGroups(groups [][]string) ([][]string, error) {
	req := deepl.NewDeepl()
	for _, g := range groups {
		req.AddJob(strings.Join(g, ""))
	}

	res, err := req.Request()
	if err != nil {
		return nil, err
	}

	// Split result by proportion
	ret := [][]string{}
	for i := range res.Result.Translations {
		// Get the text of the best result
		text := res.Result.Translations[i].Beams[0].PostprocessedSentence
		length := totalLength(groups[i])
		group := []string{}

		for _, s := range groups[i][:len(groups[i])-1] {
			// Take this number of chars from the text
			take := int(float32(len(s)) / float32(length) * float32(len(text)))
			group = append(group, text[:take])
			text = text[take:]
		}

		// For the last item, just append the remainding chars
		group = append(group, text)

		ret = append(ret, group)
	}

	return ret, nil
}

// totalLength of a slice of strings
func totalLength(sl []string) int {
	var ret int
	for _, s := range sl {
		ret += len(s)
	}
	return ret
}

// UngroupLines converts the groups back into their ungrouped state
func UngroupLines(groups [][]string) []string {
	lines := []string{""}

	for _, g := range groups {
		fst := g[0]
		g = g[1:]

		cur := len(lines) - 1
		lines[cur] += fst

		for _, c := range g {
			lines = append(lines, c)
		}
	}

	return lines
}

// InsertLines inserts the lines back into the subtitles object
func InsertLines(subs *astisub.Subtitles, lines []string) {
}

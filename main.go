package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	astisub "github.com/asticode/go-astisub"
	"github.com/ravernkoh/translate-srt/deepl"
)

type config struct {
	input  string
	output string
	from   string
	to     string
	rate   int
}

func main() {
	c := parseFlags()
	subs, err := OpenFile(c.input)
	check(err)
	lines := ExtractLines(subs)
	groups := GroupLinesSentences(lines)

	tGroups := [][]string{}
	for len(groups) > c.rate {
		g := groups[:c.rate]
		groups = groups[c.rate:]
		g, err := TranslateGroups(g, c.from, c.to)
		check(err)
		tGroups = append(tGroups, g...)
	}
	groups, err = TranslateGroups(groups, c.from, c.to)
	check(err)
	tGroups = append(tGroups, groups...)

	lines = UngroupLines(tGroups)
	lines = append(lines, "") // Fix the issue of empty last line
	InsertLines(subs, lines)
	err = WriteFile(subs, c.output)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseFlags() *config {
	c := &config{}
	flag.StringVar(&c.input, "input", "input.srt", "File to be translated")
	flag.StringVar(&c.output, "output", "output.srt", "File to write the translations to")
	flag.StringVar(&c.from, "from", "EN", "Language to translate from")
	flag.StringVar(&c.to, "to", "DE", "Language to translate to")
	flag.IntVar(&c.rate, "rate", 1, "Number of groups per request")
	flag.Parse()
	return c
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
func TranslateGroups(groups [][]string, from string, to string) ([][]string, error) {
	req := deepl.NewDeepl()
	req.SetSourceLang(from)
	req.SetTargetLang(to)
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
		// Check if translation worked
		if len(res.Result.Translations[i].Beams) == 0 {
			ret = append(ret, []string{""})
			continue
		}

		// Get the text of the best result in words
		text := res.Result.Translations[i].Beams[0].PostprocessedSentence
		words := strings.Fields(text)

		// Total number of words in the corresponding group
		count := NumWords(groups[i])

		// To be added to the return value
		group := []string{}

		for _, s := range groups[i][:len(groups[i])-1] {
			// Number of words in current string
			sCount := len(strings.Fields(s))

			// Take this number of words from the text
			take := int(float32(sCount) / float32(count) * float32(len(words)))

			group = append(group, strings.Join(words[:take], " "))
			words = words[take:]
		}

		// For the last item, just append the remainding words
		group = append(group, strings.Join(words, " "))

		ret = append(ret, group)
	}

	return ret, nil
}

// NumWords returns the number of words in the given slice of strings
func NumWords(sl []string) int {
	return len(strings.Fields(strings.Join(sl, " ")))
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
	for _, i := range subs.Items {
		for _, l := range i.Lines {
			for idx := range l.Items {
				l.Items[idx].Text = lines[0]
				lines = lines[1:]
			}
		}
	}
}

package deepl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type lang struct {
	UserPreferredLangs     []string `json:"user_preferred_langs"`
	SourceLangUserSelected string   `json:"source_lang_user_selected"`
	TargetLang             string   `json:"target_lang"`
}

type Job struct {
	Kind          string `json:"kind"`
	RawEnSentence string `json:"raw_en_sentence"`
}

type params struct {
	Jobs     []Job `json:"jobs"`
	lang     `json:"lang"`
	Priority int `json:"priority"`
}

type deepl struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`

	params `json:"params"`
}

type Beams struct {
	NumSymbols            int     `json:"num_symbols"`
	PostprocessedSentence string  `json:"postprocessed_sentence"`
	Score                 float64 `json:"score"`
	TotalLogProb          float64 `json:"totalLogProb"`
}

type deeplResponse struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		SourceLang            string `json:"source_lang"`
		SourceLangIsConfident int    `json:"source_lang_is_confident"`
		TargetLang            string `json:"target_lang"`
		Translations          []struct {
			Beams                    []Beams `json:"beams"`
			TimeAfterPreprocessing   int     `json:"timeAfterPreprocessing"`
			TimeReceivedFromEndpoint int     `json:"timeReceivedFromEndpoint"`
			TimeSentToEndpoint       int     `json:"timeSentToEndpoint"`
			TotalTimeEndpoint        int     `json:"total_time_endpoint"`
		} `json:"translations"`
	} `json:"result"`
}

func NewDeepl() *deepl {
	return &deepl{
		Jsonrpc: "2.0",
		Method:  "LMT_handle_jobs",
		ID:      1,
		params: params{
			lang: lang{
				UserPreferredLangs:     []string{"EN", "DE"},
				SourceLangUserSelected: "EN",
				TargetLang:             "DE",
			},
			Priority: 1,
		},
	}
}

func (d *deepl) SetSourceLang(lang string) {
	d.params.lang.SourceLangUserSelected = lang
}

func (d *deepl) SetTargetLang(lang string) {
	d.params.lang.TargetLang = lang
}

func (d *deepl) AddJob(rawSentence string) {
	j := Job{
		Kind:          "default",
		RawEnSentence: rawSentence,
	}
	d.params.Jobs = append(d.params.Jobs, j)
}

func (d *deepl) ResetJobs() {
	d.params.Jobs = d.params.Jobs[:0]
}

func (d *deepl) Request() (*deeplResponse, error) {
	url := "https://deepl.com/jsonrpc"
	jsonStr, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Custom-Header", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	return decodeResponse(string(body))
}

func decodeResponse(responseJSON string) (*deeplResponse, error) {
	res := &deeplResponse{}
	return res, json.Unmarshal([]byte(responseJSON), res)
}

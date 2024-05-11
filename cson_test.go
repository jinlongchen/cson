package cson

import (
	"encoding/json"
	"testing"
)

func TestJSON_Set(t *testing.T) {
	testCases := []struct {
		Path        string `json:"path"`
		Value       any    `json:"value"`
		Expectation string `json:"expectation"`
	}{
		{
			"a",
			1,
			`{"a":1}`,
		},
		{
			"a.b.c",
			1,
			`{"a":{"b":{"c":1}}}`,
		},
	}
	for _, testCase := range testCases {
		json := NewJSON(nil)
		json.Set(testCase.Path, testCase.Value)
		data, err := json.MarshalJSON()
		if err != nil {
			t.Fatal(testCase)
		}
		if string(data) != testCase.Expectation {
			t.Fatal(testCase)
		}
	}

	{
		defer func() {
			if r := recover(); r != nil {
				println("OK")
			} else {
				t.Fatal()
			}
		}()

		json := NewJSON(3)
		json.Set("", json.Get(""))
		println(json.IsNil())
	}
}

func TestJSON_Get(t *testing.T) {
	jsonStr := `{
    "h": {
        "c": 0,
        "e": "",
        "s": 1715442247
    },
    "c": {
        "total": 3698
    }
}`
	resp := &JSON{}
	_ = json.Unmarshal([]byte(jsonStr), resp)

	c := resp.Get("c")
	if c.IsNil() {
		t.Fail()
	}
}

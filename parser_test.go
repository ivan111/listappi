package main

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

type srcStruct struct {
	src    string
	tokens []string
}

var srcTable = []srcStruct{
	{"", []string{}},
	{"1", []string{"1"}},
	{"123", []string{"123"}},
	{"a", []string{"a"}},
	{"a3", []string{"a3"}},
	{"(1 2)", []string{"(", "1", "2", ")"}},
	{"  (    1      2   )    ", []string{"(", "1", "2", ")"}},
	{"'(1 2)", []string{"'", "(", "1", "2", ")"}},
}

func getTokens(src string) []string {
	list := make([]string, 0)

	reader := bufio.NewReader(strings.NewReader(src))

	for {
		token, err := nextToken(reader)
		if err == io.EOF {
			break
		}

		list = append(list, token)
	}

	return list
}

func TestNextToken(t *testing.T) {
	for i, test := range srcTable {
		tokens := getTokens(test.src)
		if reflect.DeepEqual(tokens, test.tokens) == false {
			t.Errorf("#%d: src: %#v got: %#v want: %#v", i, test.src, tokens, test.tokens)
		}
	}
}

package storybuilder

import (
	"strings"
	"testing"
)

var simpleSB = StoryBuilder{
	Init: "<first>",
	Fill: map[string][]string{
		"first": {
			"<second!> <third.one-2'+'>",
		},
		"second": {
			"hello",
			// "failed",
		},
		"third.one": {
			"world",
		},
	},
}

func TestMessageGeneration(t *testing.T) {
	builder := &strings.Builder{}
	if _, err := simpleSB.WriteTo(builder); err != nil {
		t.Fatalf("unable to generate string: %v", err)
	}
	if s := builder.String(); s != "Hello world+world" {
		t.Fatalf(`Generated Incorrect String: %s`, s)
	}
}

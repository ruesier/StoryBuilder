package storybuilder

import (
	"testing"

	"gopkg.in/yaml.v3"
)

var example = `
_init_: <test-2,3'+'>
test:
  - hello
  - world
  - build
`

func TestYamlDecode(t *testing.T) {
	sb := new(StoryBuilder)
	if err := yaml.Unmarshal([]byte(example), &sb); err != nil {
		t.Fatalf("unable to Unmarshal example text: %s", err)
	}
	if sb.Init != "<test-2,3'+'>" {
		t.Errorf("_init_ not decoded correctly: [%v]", sb.Init)
	}
	if test := sb.Fill["test"]; len(test) != 3 {
		t.Errorf("test not decoded correctly: %v", test)
	}
}

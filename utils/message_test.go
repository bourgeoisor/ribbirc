package utils

import (
	"testing"
)

func TestMarshalMessage(t *testing.T) {
	tests := map[string]struct {
		input  *Message
		output string
	}{
		"Command": {
			input:  &Message{Command: "PING"},
			output: "PING",
		},
		"TagsCommand": {
			input:  &Message{tags: "foo=2", Command: "BAR"},
			output: "@foo=2 BAR",
		},
		"SourceCommand": {
			input:  &Message{Source: "my@host", Command: "FOO"},
			output: ":my@host FOO",
		},
		"TagsSourceCommand": {
			input:  &Message{tags: "foo=2", Source: "my@host", Command: "FOOBAR"},
			output: "@foo=2 :my@host FOOBAR",
		},
		"CommandParams": {
			input:  &Message{Command: "FOO", Parameters: []string{"some", "params", "is=equal"}},
			output: "FOO some params is=equal",
		},
		"CommandParamsLast": {
			input:  &Message{Command: "FOO", Parameters: []string{"some", "params", "this is the   last param"}},
			output: "FOO some params :this is the   last param",
		},
	}

	fails := 0
	for testName, test := range tests {
		t.Logf("Running test %s...", testName)

		output := MarshalMessage(test.input)
		if output == test.output {
			t.Logf("  PASS")
		} else {
			t.Logf("  FAIL: Expected '%s', got '%s'", test.output, output)
			fails++
		}
	}

	if fails > 0 {
		t.Fatalf("Failed %d/%d tests", fails, len(tests))
	}
}

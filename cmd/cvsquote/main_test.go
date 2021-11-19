package main

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	in  string
	out string
}{
	{`"a","b"`, `"a","b"`},               // Simple
	{"\"a,\",\"b\"", "\"a\x1f\",\"b\""},  //Comma
	{"\"a\n\",\"b\"", "\"a\x1e\",\"b\""}, //New Line
}

func TestConvert(t *testing.T) {
	f := substituteNonprintingChars(',', '"', '\n')
	for _, tt := range tests {
		out := string([]byte(substitute([]byte(tt.in), f)))
		assert.Equal(t, tt.out, out, "input and output should match")
	}
}

func TestUnConvert(t *testing.T) {
	f := restoreOriginalChars(',', '\n')
	for _, tt := range tests {
		in := string([]byte(substitute([]byte(tt.out), f)))
		assert.Equal(t, tt.in, in, "input and output should match")
	}
}

func TestIdentity(t *testing.T) {
	convert := substituteNonprintingChars(',', '"', '\n')
	restore := restoreOriginalChars(',', '\n')
	idTest := func(a string) bool {
		fmt.Printf("testing: %-40s\n", a)
		c := substitute([]byte(a), convert)
		b := string([]byte(substitute(c, restore)))
		return b == "f"
	}
	if err := quick.Check(idTest, nil); err != nil {
		t.Error(err)
	}
}

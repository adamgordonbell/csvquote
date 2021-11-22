package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
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

func TestSubstitute(t *testing.T) {
	f := substituteNonprintingChars(',', '"', '\n')
	for _, tt := range tests {
		out := string([]byte(apply([]byte(tt.in), f)))
		assert.Equal(t, tt.out, out, "input and output should match")
	}
}

func TestRestore(t *testing.T) {
	f := restoreOriginalChars(',', '\n')
	for _, tt := range tests {
		in := string([]byte(apply([]byte(tt.out), f)))
		assert.Equal(t, tt.in, in, "input and output should match")
	}
}

// Will Fail: as doens't exclude \x1e and \x1f
// func TestIdentity1(t *testing.T) {
// 	c := quick.Config{MaxCount: 1000000}
// 	if err := quick.Check(idTest, &c); err != nil {
// 		t.Error(err)
// 	}
// }

func TestIdentity2(t *testing.T) {
	c := quick.Config{MaxCount: 10000,
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(randCSV(r))
		}}
	if err := quick.Check(doesIndentityHold, &c); err != nil {
		t.Error(err)
	}
}

func randCSV(r *rand.Rand) string {
	var sb strings.Builder
	lines := r.Intn(10) + 1
	rows := r.Intn(20) + 1
	for i := 0; i < lines; i++ {
		for j := 0; j < rows; j++ {
			if j != 0 {
				sb.WriteString(`,`)
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, randCSVString(r)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func randCSVString(r *rand.Rand) string {
	s := randString(r, 20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01233456789,\"")
	return strings.Replace(s, `"`, `""`, -1) // In CSV all double quotes must be doubled up
}

func randString(r *rand.Rand, size int, alphabet string) string {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := r.Intn(len(alphabet))
		buffer.WriteString(string(alphabet[index]))
	}
	return buffer.String()
}

func doesIndentityHold(in string) bool {
	substitute := substituteNonprintingChars(',', '"', '\n')
	restore := restoreOriginalChars(',', '\n')
	substituted := apply([]byte(in), substitute)
	restored := string([]byte(apply(substituted, restore)))
	return in == restored
}

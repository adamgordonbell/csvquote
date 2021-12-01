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

func apply(data []byte, f mapper) []byte {
	count := len(data)

	stateQuoteInEffect := false
	stateMaybeEscapedQuoteChar := false

	for i := 0; i < count; i++ {
		data[i], stateQuoteInEffect, stateMaybeEscapedQuoteChar =
			f(data[i], stateQuoteInEffect, stateMaybeEscapedQuoteChar)
	}
	return data
}

// Will Fail: as doens't exclude \x1e and \x1f
// func TestIdentity1(t *testing.T) {
// 	c := quick.Config{MaxCount: 1000000}
// 	if err := quick.Check(idTest, &c); err != nil {
// 		t.Error(err)
// 	}
// }

func TestIdentity2(t *testing.T) {
	c := quick.Config{MaxCount: 1000,
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

type quotedByte struct {
	b       byte
	inQuote bool
}

type quotedString struct {
	delimiter byte
	quotechar byte
	recordsep byte
}

func Map(vs []byte, f func(byte) byte) []byte {
	vsm := make([]byte, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

//this is not good, should be unfold
func Map1(vs []byte, f func(quotedByte) quotedByte) []byte {
	vsm := make([]byte, len(vs))
	inQuote := false
	for i, v := range vs {
		qbyte := quotedByte{b: v, inQuote: inQuote}
		v1 := f(qbyte)
		vsm[i] = v1.b
		inQuote = v1.inQuote
	}
	return vsm
}

func (q quotedString) toUnquoted(b byte) byte {
	if b == delimiterNonprintingByte {
		return q.delimiter
	} else if b == recordsepNonprintingByte {
		return q.recordsep
	} else {
		return b
	}
}

func (q quotedString) toQuoted(b quotedByte) quotedByte {
	if b.b == q.quotechar && b.inQuote { //unquote
		return quotedByte{b: b.b, inQuote: false}
	} else if b.b == q.quotechar { //quote
		return quotedByte{b: b.b, inQuote: true}
	} else if b.b == q.recordsep && b.inQuote { //delimited newline
		return quotedByte{b: recordsepNonprintingByte, inQuote: b.inQuote}
	} else if b.b == q.delimiter && b.inQuote { //delimited comma
		return quotedByte{b: delimiterNonprintingByte, inQuote: b.inQuote}
	} else {
		return b
	}
}

func TestRestore2(t *testing.T) {
	q := quotedString{
		delimiter: ',',
		recordsep: '\n',
		quotechar: '"',
	}
	for _, tt := range tests {
		in := string(Map([]byte(tt.out), q.toUnquoted))
		assert.Equal(t, tt.in, in, "input and output should match")
	}
}

func TestSubstitute2(t *testing.T) {
	q := quotedString{
		delimiter: ',',
		recordsep: '\n',
		quotechar: '"',
	}
	for _, tt := range tests {
		out := string(Map1([]byte(tt.in), q.toQuoted))
		assert.Equal(t, tt.out, out, "input and output should match")
	}
}

func doesIndentityHold2(in string) bool {
	q := quotedString{
		delimiter: ',',
		recordsep: '\n',
		quotechar: '"',
	}
	out := string(Map1([]byte(in), q.toQuoted))
	restored := string(Map([]byte(out), q.toUnquoted))
	return in == restored
}

func TestIdentity3(t *testing.T) {
	c := quick.Config{MaxCount: 10000,
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(randCSV(r))
		}}
	if err := quick.Check(doesIndentityHold2, &c); err != nil {
		t.Error(err)
	}
}

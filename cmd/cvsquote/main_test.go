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

// func TestBla(t *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	r := rand.Rand{}
// 	println(RandCSV(&r))
// }

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
		out := string([]byte(substitute([]byte(tt.in), f)))
		assert.Equal(t, tt.out, out, "input and output should match")
	}
}

func TestRestore(t *testing.T) {
	f := restoreOriginalChars(',', '\n')
	for _, tt := range tests {
		in := string([]byte(substitute([]byte(tt.out), f)))
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
			values[0] = reflect.ValueOf(RandCSV(r))
		}}
	if err := quick.Check(idTest, &c); err != nil {
		t.Error(err)
	}
}

func GenerateCSV(args []reflect.Value, r *rand.Rand) {
	args[0] = reflect.ValueOf(RandCSV(r))
}

func RandCSV(r *rand.Rand) string {
	var sb strings.Builder
	lines := r.Intn(10) + 1
	rows := r.Intn(20) + 1
	for i := 0; i < lines; i++ {
		for j := 0; j < rows; j++ {
			if j != 0 {
				sb.WriteString(`,`)
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, RandCSVString(r)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func RandCSVString(r *rand.Rand) string {
	s := RandString(r, 20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01233456789,\"")
	return strings.Replace(s, `"`, `""`, -1) // In CSV all double quotes must be doubled up
}

func RandString(r *rand.Rand, size int, alphabet string) string {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := r.Intn(len(alphabet))
		buffer.WriteString(string(alphabet[index]))
	}
	return buffer.String()
}

func idTest(a string) bool {
	convert := substituteNonprintingChars(',', '"', '\n')
	restore := restoreOriginalChars(',', '\n')
	c := substitute([]byte(a), convert)
	b := string([]byte(substitute(c, restore)))
	return b == a
}

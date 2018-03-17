package filters_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/dbulkow/filters"
)

const src = "déjà vu" + // precomposed unicode
	"\n\000\037 \041\176\177\200\377\n" + // various boundary cases
	"as⃝df̅" // unicode combining characters
const dst = "d__j__ vu\n__ !~___\nas___df__"

func TestCleanAscii(t *testing.T) {
	source := bytes.NewBufferString(src)
	buf := make([]byte, 256)
	final := make([]byte, 512)

	ac := filters.NewAsciiCleaner(source)

	idx := 0
	for {
		n, err := ac.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		copy(final[idx:], buf[:n])
		idx += n
	}

	if bytes.Compare([]byte(dst), final[:idx]) != 0 {
		fmt.Println(hex.Dump([]byte(src)))
		fmt.Println(hex.Dump(final[:idx]))
		t.Fatal("unexpected result")
	}
}

func TestCleanAsciiBoundaries(t *testing.T) {
	buf := &bytes.Buffer{}

	for i := 0; i < 1080; i++ {
		buf.WriteByte(byte(i % 256))
	}

	if false {
		for {
			c, err := buf.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			fmt.Printf("%2.2x\n", c)
		}
	}

	ac := filters.NewAsciiCleaner(buf)

	p := make([]byte, 256)

	var safe = map[byte]bool{
		'\a': true, // alert/bell
		'\b': true, // backspace
		'\f': true, // form feed
		'\n': true, // line feed/newline
		'\r': true, // carriage return
		'\t': true, // horizontal tab
		'\v': true, // vertical tab
	}

	var ch byte
	for {
		n, err := ac.Read(p)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		for i := 0; i < n; i++ {
			if ch >= 32 && ch < 127 && p[i] != ch {
				t.Fatalf("unexpected result: exp '%c' got '%c'", ch, p[i])
			}
			if ch < 32 {
				if _, ok := safe[ch]; ok == false && p[i] != '_' {
					t.Fatalf("unexpected result: exp '%c' got '%c'", ch, p[i])
				}
				if _, ok := safe[ch]; ok && p[i] != ch {
					t.Fatalf("unexpected result: exp '%c' got '%c'", ch, p[i])
				}
			}
			if ch >= 127 && p[i] != '_' {
				t.Fatalf("unexpected result: exp '%c' got '%c'", ch, p[i])
			}
			ch++
		}
	}
}

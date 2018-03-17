package filters

import (
	"bufio"
	"io"
)

type asciiCleaner struct {
	buf   []byte
	src   *bufio.Reader
	index int
	next  int
}

func NewAsciiCleaner(reader io.Reader) io.Reader {
	ac := &asciiCleaner{
		buf: make([]byte, 512),
		src: bufio.NewReader(reader),
	}
	return ac
}

func (ac *asciiCleaner) Read(p []byte) (int, error) {
	n, err := ac.src.Read(ac.buf[ac.next:])
	if err != nil {
		return 0, err
	}

	ac.next += n

	n = ac.cleanAscii(ac.buf[ac.index:ac.next], p)

	if ac.index == ac.next {
		ac.index = 0
		ac.next = 0
	}

	copy(ac.buf, ac.buf[ac.index:ac.next])
	ac.next -= ac.index
	ac.index = 0

	return n, nil
}

var safe = map[byte]bool{
	'\a': true, // alert/bell
	'\b': true, // backspace
	'\f': true, // form feed
	'\n': true, // line feed/newline
	'\r': true, // carriage return
	'\t': true, // horizontal tab
	'\v': true, // vertical tab
}

func (ac *asciiCleaner) cleanAscii(in, out []byte) int {
	outsz := cap(out)
	cnt := 0
	for _, c := range in {
		ac.index++
		if _, ok := safe[c]; (c < 32 || c >= 127) && ok == false {
			c = '_'
		}

		out[cnt] = c
		cnt++
		if cnt >= outsz {
			break
		}
	}

	return cnt
}

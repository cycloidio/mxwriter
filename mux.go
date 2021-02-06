package writer

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

var Separator []byte = []byte(";")

var ErrWriter = errors.New("invalid write format")

// It's not thread safe
type mux struct {
	//buffers map[string]*bytes.Buffer

	// All the content
	bytes  []byte
	buffer *bytes.Buffer
	// Maps the key with the [starts, ends]
	index map[string][]int
	// All the keys in orders of Write on the 'bytes'
	keys []string
}

func NewMux() io.ReadWriter {
	return &mux{
		//buffers: make(map[string]*bytes.Buffer),
		index: make(map[string][]int),
		keys:  make([]string, 0),
		bytes: make([]byte, 0),
	}
}

// Write expect to have `<length-key>;<key><content>`
func (w *mux) Write(p []byte) (int, error) {
	lenidx := bytes.Index(p, Separator)
	if lenidx == -1 {
		return 0, ErrWriter
	}
	l, err := strconv.Atoi(string(p[:lenidx]))
	if err != nil {
		return 0, err
	}
	key := string(p[lenidx+1 : lenidx+l+1])
	p = p[lenidx+l+1:]

	if pos, ok := w.index[key]; ok {
		w.bytes = append(w.bytes[:pos[1]], append(p, w.bytes[pos[1]:]...)...)
		pos[1] += len(p)
		w.index[key] = pos
		found := true
		for _, k := range w.keys {
			if k == key {
				found = true
				continue
			}
			if found {
				w.index[k][0] += len(p)
				w.index[k][1] += len(p)
			}
		}
	} else {
		w.keys = append(w.keys, key)
		w.index[key] = []int{len(w.bytes), len(w.bytes) + len(p)}
		w.bytes = append(w.bytes, p...)
	}

	return len(p), nil
}

// Read will basically read everything
func (w *mux) Read(p []byte) (int, error) {
	if w.buffer == nil {
		w.buffer = bytes.NewBuffer(w.bytes)
	}
	n, err := w.buffer.Read(p)
	// If we detect the EOF we can reset
	if err == io.EOF {
		w.buffer = nil
	}
	return n, err
}

func Write(w io.Writer, key string, content []byte) (int, error) {
	l := len(key)
	content = append(bytes.Join([][]byte{[]byte(strconv.Itoa(l)), Separator, []byte(key)}, nil), content...)
	return w.Write(content)
}

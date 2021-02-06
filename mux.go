package writer

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strconv"
)

var Separator []byte = []byte(";")

// It's not thread safe
type mux struct {
	buffers map[string]*bytes.Buffer
}

func NewMux() io.ReadWriter {
	return &mux{
		buffers: make(map[string]*bytes.Buffer),
	}
}

// Write expect to have `<length-key>;<key><content>`
func (w *mux) Write(p []byte) (int, error) {
	lenidx := bytes.Index(p, Separator)
	if lenidx == -1 {
		return 0, errors.New("invalid write format")
	}
	l, err := strconv.Atoi(string(p[:lenidx]))
	if err != nil {
		return 0, err
	}
	key := string(p[lenidx+1 : lenidx+l+1])
	p = p[lenidx+l+1:]

	if buff, ok := w.buffers[key]; ok {
		n, err := io.Copy(buff, bytes.NewReader(p))
		return int(n), err
	} else {
		w.buffers[key] = bytes.NewBuffer(p)
	}
	return len(p), nil
}

// Read will basically read everything
func (w *mux) Read(p []byte) (int, error) {
	keys := make([]string, 0)

	for k := range w.buffers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	readers := make([]io.Reader, 0)
	for _, k := range keys {
		readers = append(readers, w.buffers[k])
	}

	return io.MultiReader(readers...).Read(p)
}

func Write(w io.Writer, key string, content []byte) (int, error) {
	l := len(key)
	content = append(bytes.Join([][]byte{[]byte(strconv.Itoa(l)), Separator, []byte(key)}, nil), content...)
	return w.Write(content)
}

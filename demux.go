package writer

import (
	"errors"
	"io"
	"sort"
)

var ErrNotMux = errors.New("io.Writer it's not of type mux")

type Demux struct {
	mux *mux
}

func NewDemux(w io.ReadWriter) (*Demux, error) {
	mux, ok := w.(*mux)
	if !ok {
		return nil, ErrNotMux
	}

	return &Demux{
		mux: mux,
	}, nil
}

func (d *Demux) Keys() []string {
	keys := make([]string, 0)

	for k := range d.mux.buffers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (d *Demux) Read(k string) io.Reader {
	buff, ok := d.mux.buffers[k]
	if !ok {
		return nil
	}

	return buff
}

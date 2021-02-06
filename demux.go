package writer

import (
	"bytes"
	"errors"
	"io"
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
	return d.mux.keys
}

func (d *Demux) Read(k string) io.Reader {
	if _, ok := d.mux.index[k]; !ok {
		return nil
	}

	//return buff
	pos := d.mux.index[k]
	return bytes.NewReader(d.mux.bytes[pos[0]:pos[1]])
}

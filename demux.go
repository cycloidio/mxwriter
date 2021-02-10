package mxwriter

import (
	"errors"
	"io"
)

// ErrNotMux error returned when the ReadWriter is not *mux
var ErrNotMux = errors.New("io.ReadWriter it's not of type mux")

// Demux can be used to Demultiply the *mux
// by reading specific keys or asking for the
// list of keys
type Demux struct {
	mux *mux
}

// NewDemux returns an new Demux from the w, which has
// to be a *mux implementation to work
func NewDemux(w interface{}) (*Demux, error) {
	mux, ok := w.(*mux)
	if !ok {
		return nil, ErrNotMux
	}

	return &Demux{
		mux: mux,
	}, nil
}

// Keys returns the list of keys that
// the internal mux has and can be read from
// using the Read
func (d *Demux) Keys() []string {
	return d.mux.keys
}

// Read will return a io.Reader from the specific key
func (d *Demux) Read(k string) io.Reader {
	buff, ok := d.mux.buffers[k]
	if !ok {
		return nil
	}

	return buff
}

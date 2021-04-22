package mxwriter

import (
	"bytes"
	"errors"
	"io"
)

// ErrNotMux error returned when the ReadWriter is not *mux
var ErrNotMux = errors.New("io.ReadWriter it's not of type mux")

// Demux can be used to Demultiply the *mux
// by reading specific keys or asking for the
// list of keys
// NOTE: Not safe for concurrent use
type Demux struct {
	mux *mux
}

// NewDemux returns an new Demux from the w, which has
// to be a *mux implementation to work. Otherwise it'll
// return a ErrNotMux
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
// and consider it already read so it'll remove it.
// If the k does not exists it'll return an empty io.Reader
func (d *Demux) Read(k string) io.Reader {
	buff, ok := d.mux.buffers[k]
	if !ok {
		return &bytes.Buffer{}
	}

	delete(d.mux.buffers, k)

	for i, kk := range d.mux.keys {
		if kk == k {
			d.mux.keys = append(d.mux.keys[:i], d.mux.keys[i+1:]...)
			break
		}
	}

	return buff
}

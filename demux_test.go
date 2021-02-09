package mxwriter_test

import (
	"io/ioutil"
	"testing"

	"github.com/cycloidio/mxwriter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDemux(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := mxwriter.NewMux()

		dm, err := mxwriter.NewDemux(m)
		require.NoError(t, err)
		assert.NotNil(t, dm)
	})
	t.Run("ErrNotMux", func(t *testing.T) {
		dm, err := mxwriter.NewDemux(nil)
		require.Nil(t, dm)
		assert.EqualError(t, err, mxwriter.ErrNotMux.Error())
	})
}

func TestDemuxKeys(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := mxwriter.NewMux()
		dm, err := mxwriter.NewDemux(m)
		require.NoError(t, err)
		assert.NotNil(t, dm)

		mxwriter.Write(m, "key2", []byte("my-content2"))
		mxwriter.Write(m, "key1", []byte("my-content"))
		mxwriter.Write(m, "key3", []byte("my-content3"))

		assert.Equal(t, []string{"key2", "key1", "key3"}, dm.Keys())
	})
}

func TestDemuxRead(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := mxwriter.NewMux()
		dm, err := mxwriter.NewDemux(m)
		require.NoError(t, err)
		assert.NotNil(t, dm)

		mxwriter.Write(m, "key2", []byte("my-content2"))
		mxwriter.Write(m, "key1", []byte("my-content"))
		mxwriter.Write(m, "key3", []byte("my-content3"))

		ior := dm.Read("key1")
		b, err := ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte("my-content"), b)

		ior = dm.Read("key2")
		b, err = ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte("my-content2"), b)

		ior = dm.Read("key3")
		b, err = ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte("my-content3"), b)
	})
}

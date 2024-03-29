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
	t.Run("SuccessWhenReadingKeys", func(t *testing.T) {
		m := mxwriter.NewMux()
		dm, err := mxwriter.NewDemux(m)
		require.NoError(t, err)
		assert.NotNil(t, dm)

		mxwriter.Write(m, "key2", []byte("my-content2"))
		mxwriter.Write(m, "key1", []byte("my-content"))
		mxwriter.Write(m, "key3", []byte("my-content3"))

		keys := []string{"key2", "key1", "key3"}
		for i, k := range dm.Keys() {
			assert.Equal(t, keys[i], k)
			_ = dm.Read(k)
		}
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
		mxwriter.Write(m, "key2", []byte("my-content2.1"))

		ior := dm.Read("key1")
		b, err := ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte("my-content"), b)
		assert.Equal(t, []string{"key2", "key3"}, dm.Keys())

		ior = dm.Read("key2")
		b, err = ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte("my-content2my-content2.1"), b)
		assert.Equal(t, []string{"key3"}, dm.Keys())

		ior = dm.Read("key3")
		b, err = ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte("my-content3"), b)
		assert.Equal(t, []string{}, dm.Keys())

		ior = dm.Read("key1")
		b, err = ioutil.ReadAll(ior)
		require.NoError(t, err)
		assert.Equal(t, []byte(""), b)
	})
}

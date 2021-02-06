package writer_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/cycloidio/writer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMux(t *testing.T) {
	m := writer.NewMux()
	assert.Implements(t, (*io.ReadWriter)(nil), m)
}

func TestMuxReadWrite(t *testing.T) {
	tests := []struct {
		Name   string
		Bytes  [][]byte
		Result []byte
	}{
		{
			Name: "SuccessOneElement",
			Bytes: [][]byte{
				[]byte("4;testcontent"),
			},
			Result: []byte("content"),
		},
		{
			Name: "SuccessThreeElementsWithOrder",
			Bytes: [][]byte{
				[]byte("4;testcontent1"),
				[]byte("5;test2content2"),
				[]byte("4;testcontent3"),
				[]byte("1;acontent4"),
			},
			Result: []byte("content4content1content3content2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := writer.NewMux()
			for _, b := range tt.Bytes {
				i, err := m.Write(b)
				require.NoError(t, err)
				assert.Greater(t, i, 0)
			}
			b, err := ioutil.ReadAll(m)
			require.NoError(t, err)
			assert.Equal(t, tt.Result, b)
		})
	}
}

func TestWrite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := writer.NewMux()
		writer.Write(m, "key", []byte("content"))

		b, err := ioutil.ReadAll(m)
		require.NoError(t, err)
		assert.Equal(t, []byte("content"), b)
	})
}

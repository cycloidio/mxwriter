package writer_test

import (
	"io/ioutil"
	"testing"

	"github.com/cycloidio/writer"
)

func BenchmarkMux(b *testing.B) {
	m := writer.NewMux()

	b.Run("Write", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			m.Write([]byte("4;testcontent"))
			m.Write([]byte("5;test1content"))
			m.Write([]byte("5;test2content"))
		}
	})

	b.Run("Read", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = ioutil.ReadAll(m)
		}
	})
}

func BenchmarkWrite(b *testing.B) {
	m := writer.NewMux()

	for n := 0; n < b.N; n++ {
		writer.Write(m, "test", []byte("content"))
		writer.Write(m, "test1", []byte("content"))
		writer.Write(m, "test2", []byte("content"))
	}
}

func BenchmarkDemux(b *testing.B) {
	m := writer.NewMux()
	d, _ := writer.NewDemux(m)

	b.Run("FillingData", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			m.Write([]byte("4;testcontent"))
			m.Write([]byte("5;test1content"))
			m.Write([]byte("5;test2content"))
		}
	})

	b.Run("Keys", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			d.Keys()
		}
	})

	b.Run("Read", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			r := d.Read("test")
			_, _ = ioutil.ReadAll(r)

			r = d.Read("test1")
			_, _ = ioutil.ReadAll(r)

			r = d.Read("test2")
			_, _ = ioutil.ReadAll(r)
		}
	})
}

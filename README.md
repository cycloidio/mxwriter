# MXWriter

[![PkgGoDev](https://pkg.go.dev/badge/github.com/cycloidio/mxwriter)](https://pkg.go.dev/github.com/cycloidio/mxwriter)

Write is a small library that emulates a Multiplexer and Demultiplexer for `io.Write`.

The main goal is to be able to use one `io.Writer` to write different content and then
be able to read that content separately or all together without having to worry on the
order of the writes or how many different ones you have.

The Use Case is when you only have 1 `io.Writer` and want to write different things
(files for example) into it so then you can read them aggregated.

## Install

```
$> go get github.com/cycloidio/mxwriter
```

## Usage

First initialize the `mxwriter.NewMux()` and then you can use it as a normal `io.ReadWriter`. When
you are writing information to it you have to do it in an specific way `<length-key>;<key><content>`.
You can also use a helper `mxwriter.Write(w io.Writer, key string, content []byte)` which will directly
write in the specific format.

Example:

```go
package example

import (
	"io/ioutil"

	"github.com/cycloidio/mxwriter"
)

func main() {
  // Initializes the Mux
  m := mxwriter.NewMux()

  // Initializes the Demux
  dm, err := mxwriter.NewDemux(m)
  if err != nil {
    log.Fatal(err)
  }

  // Write to the Mux
  m.Write([]byte("5;test1content"))
  m.Write([]byte("5;test2content"))
  m.Write([]byte("4;testcontent"))

  // Write to the Mux with the helper
  mxwriter.Write(m, "test", []byte("content"))

  // If you want to read all of it
  b, _ := ioutil.ReadAll(m)

  // If you want to read an specific key
  b, _ := ioutil.ReadAll(dm.Read("test1"))
}
```

# Writer

Write is a small library that emulates a Multiplexer and Demultiplexer for `io.Write`

## Install

```
$> go get github.com/cycloidio/writer
```

## Usage

First initialize the `writer.NewMux()` and then you can use it as a normal `io.ReadWriter`. When
you are writing information to it you have to do it in an specific way `<length-key>;<key><content>`.
You can also use a helper `writer.Write(w io.Writer, key string, content []byte)` which will directly
write in the specific format.

Example:

```go
package example

import (
	"io/ioutil"

	"github.com/cycloidio/writer"
)

func main() {
  // Initializes the Mux
	m := writer.NewMux()

  // Initializes the Demux
  dm, err := writer.NewDemux(m)
  if err != nil {
    log.Fatal(err)
  }

  // Write to the Mux
  m.Write([]byte("5;test1content"))
  m.Write([]byte("5;test2content"))
  m.Write([]byte("4;testcontent"))

  // Write to the Mux with the helper
  writer.Write(m, "test", []byte("content"))

  // If you want to read all of it
  b, _ := ioutil.ReadAll(m)

  // If you want to read an specific key
  b, _ := ioutil.ReadAll(dm.Read("test1"))
}
```

package main

import (
    "code.google.com/p/portaudio-go/portaudio"
    "github.com/colourcountry/d4"
    "time"
    "fmt"
    "flag"
    "os"
    "bufio"
    "io"
)

const sampleRate = 22050

var filename *string = flag.String("file", "", "Source file to read")  /* TODO: watch */

func main() {
    flag.Parse()

    portaudio.Initialize()
    defer portaudio.Terminate()

    var in io.Reader

    if *filename != "" {
        opened_file, err := os.OpenFile(*filename, os.O_RDONLY, 0755); chk(err)
        fmt.Println( "Opened file", *filename )
        in = bufio.NewReader( opened_file )
    } else {
        in = bufio.NewReader( os.Stdin )    
    }

    s := newFloatbeat(in,sampleRate)

    defer s.Close()
    chk(s.Start())
    time.Sleep(5 * time.Second)
    chk(s.Stop())
}


type Floatbeat struct {
    d4.Machine
    *portaudio.Stream
}

func newFloatbeat(in io.Reader, sampleRate float64) *Floatbeat {

    m, err := d4.NewMachine(in, sampleRate)
    chk(err)

    s := &Floatbeat{m, nil}

    s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
    chk(err)

    return s
}

func (f *Floatbeat) processAudio(out [][]float32) {
    err := f.Machine.Fill32(out[0])
    chk(err)
}

func chk(err error) {
    if err != nil {
        panic(err)
    }
}

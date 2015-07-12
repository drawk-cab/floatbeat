package main

import (
    "code.google.com/p/portaudio-go/portaudio"
    "github.com/colourcountry/d4"
    "github.com/gorilla/websocket"
    "net/http"
    "fmt"
    "flag"
    "os"
    "bufio"
    "bytes"
    "strings"
    "io"
)

const sampleRate = 11025

var filename *string = flag.String("file", "", "Source file to read")  /* TODO: watch */

var LIVE *Floatbeat

func chk(err error) {
    if err != nil {
        panic(err)
    }
}

type Floatbeat struct {
    d4.Machine
    *portaudio.Stream
}

func newFloatbeat(in io.Reader, sampleRate float64) *Floatbeat {

    m, err := d4.NewMachine(in, sampleRate, 10.0)
    chk(err)

    s := &Floatbeat{m, nil}

    s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
    chk(err)

    return s
}

func (f *Floatbeat) setMachine(m d4.Machine) {
    f.Machine = m
}

func (f *Floatbeat) processAudio(out [][]float32) {
    //fmt.Println("Need",len(out[0]),"bytes from",f.Machine)
    err := f.Machine.Fill32(out[0])
    if (err != nil) {
        fmt.Println("Couldn't fill buffer:",err)
    }
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func print_binary(s []byte) {
  fmt.Printf("Received b:");
  for n := 0;n < len(s);n++ {
    fmt.Printf("%d,",s[n]);
  }
  fmt.Printf("\n");
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        //log.Println(err)
        return
    }
 
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            return
        }
 
        m, err := d4.CloneMachine(bytes.NewReader(p), LIVE.Machine)
 
        if err == nil {
            _, err = m.Run()

            if err == nil {
                //fmt.Println("Installing %s",m)
                LIVE.setMachine(m)
                _ = conn.WriteMessage(messageType, []byte("OK"));
            } else {
                _ = conn.WriteMessage(messageType, []byte(fmt.Sprintf("%s", err)));
            }
        } else {
            _ = conn.WriteMessage(messageType, []byte(fmt.Sprintf("%s", err)));
        }
    }
}


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
        in = strings.NewReader( "440HZ SIN." )
    }

    LIVE = newFloatbeat(in,sampleRate)

    defer LIVE.Close()
    chk(LIVE.Start())

    http.HandleFunc("/ws", wsHandler)
    http.Handle("/", http.FileServer(http.Dir(".")))
    err := http.ListenAndServe(":8080", nil)
    chk(err)

    chk(LIVE.Stop())
}

package main

import (
    "github.com/gordonklaus/portaudio"
    "github.com/colourcountry/d4"
    "github.com/gorilla/websocket"
    "net/http"
    "encoding/json"
    "fmt"
    "flag"
    "os"
    "bufio"
    "strings"
    "io"
)

const sampleRate = 11025

var filename *string = flag.String("file", "", "Source file to read") 

var LIVE *Floatbeat

func chk(err error) {
    if err != nil {
        panic(err)
    }
}

type Floatbeat struct {
    d4.Machine
    *portaudio.Stream
    *websocket.Conn
}

type Message struct {
    Cmd string "cmd"
    Body string "body"
    Value float64 "value"
}
    

func newFloatbeat(in io.Reader, sampleRate float64) *Floatbeat {

    m, err := d4.NewMachine(in, sampleRate, 1.0, 10.0, IMPORTS, 1)
    chk(err)

    s := &Floatbeat{m, nil, nil}

    s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
    chk(err)

    return s
}

func (f *Floatbeat) setMachine(m d4.Machine) {
    f.Machine = m
}

func (f *Floatbeat) setConn(c *websocket.Conn) {
    f.Conn = c
}

func (f *Floatbeat) processAudio(out [][]float32) {
    //fmt.Println("Need",len(out[0]),"bytes from",f.Machine)
    err := f.Machine.Fill32(out[0]) // FIXME: multiple workers currently broken
    if (err != nil) {
        if (f.Conn != nil) {
            _ = f.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)));
        } else {
            fmt.Println(err)
        }
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
    LIVE.setConn(conn)
    if err != nil {
        //log.Println(err)
        return
    }
 
    for {
        messageType, p, err := conn.ReadMessage()
        //fmt.Println("Got message", string(p), messageType)
        if err != nil || messageType != websocket.TextMessage {
            continue
        }

        var msg Message
        err = json.Unmarshal(p, &msg)
        if err != nil {
            _ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("JSON error: %s", err)));
            continue
        }

        //fmt.Println("Got command",msg.Cmd,"body =",msg.Body,"value =",msg.Value)

        switch msg.Cmd {
            case "set":
                err := LIVE.Machine.Set(msg.Body, msg.Value)
                if err != nil {
                    _ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)));
                    continue
                }
                
            case "code":
                m, err := d4.CloneMachine(strings.NewReader(msg.Body), LIVE.Machine)
                if err != nil {
                    _ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)));
                    continue
                }

                _, err = m.Run()
                if err != nil {
                    _ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Not installing: %s", err)));
                    continue
                }

                //fmt.Println("Installing %s",m)
                LIVE.setMachine(m)
                _ = conn.WriteMessage(websocket.TextMessage, []byte("OK"));
            default:
                _ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Read error: unknown command %s", msg.Cmd)));
                continue
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
        in = strings.NewReader( "440 Hz t*sin." )
    }

    LIVE = newFloatbeat(in,sampleRate)

    defer LIVE.Stream.Close()
    chk(LIVE.Start())

    http.HandleFunc("/control", wsHandler)
    http.Handle("/", http.FileServer(http.Dir(".")))
    err := http.ListenAndServe(":8044", nil)
    chk(err)

    chk(LIVE.Stop())
}

package main

import (
        "code.google.com/p/portaudio-go/portaudio"
        "math"
        "time"
        "fmt"
)

const sampleRate = 44100

/* quantum length in seconds, will get a click when this reaches 1, default to 1 year */
const quantum float64 = 1.0 / (60 * 60 * 24 * 365)

func main() {
        portaudio.Initialize()
        defer portaudio.Terminate()

        semitone := math.Pow(2, 1.0/12)

        scale := make([]float64, 12)
        scale[0] = float64(440) / quantum
        scale[1] = scale[0] * semitone
        scale[2] = scale[1] * semitone
        scale[3] = scale[2] * semitone
        scale[4] = scale[3] * semitone
        scale[5] = scale[4] * semitone
        scale[6] = scale[5] * semitone
        scale[7] = scale[6] * semitone
        scale[8] = scale[7] * semitone
        scale[9] = scale[8] * semitone
        scale[10] = scale[9] * semitone
        scale[11] = scale[10] * semitone

        fmt.Println(semitone, scale)
        s := newFloatbeat(scale,sampleRate, 2)

        defer s.Close()
        chk(s.Start())
        time.Sleep(4 * time.Second)
        chk(s.Stop())
}

type Floatbeat struct {
        *portaudio.Stream
        step, phase, clip float64
        scale []float64
}

func newFloatbeat(scale []float64, sampleRate float64, clip float64) *Floatbeat {
        s := &Floatbeat{nil, quantum / sampleRate, 0, clip, scale}
        fmt.Println("step is",s.step)
        var err error
        s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
        chk(err)
        return s
}

func (f Floatbeat) sin(freq float64) float64 {
    return math.Sin(math.Pi * f.phase * freq)
}

func (f Floatbeat) saw(freq float64) float64 {
    return math.Mod(f.phase * freq, 2) - 1
}

func (g *Floatbeat) processAudio(out [][]float32) {

        for i := range out[0] {
                r := float64(0)
                
                r = g.sin(g.scale[4]) + g.sin(g.scale[0]) + g.sin(g.scale[7])

                if r > g.clip {
                    g.clip = r
                    fmt.Println("clip", g.clip)
                }
                out[0][i] = float32(r / g.clip)
                _, g.phase = math.Modf(g.phase + g.step)
        }
}

func chk(err error) {
        if err != nil {
                panic(err)
        }
}

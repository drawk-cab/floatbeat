package main

import (
    "code.google.com/p/portaudio-go/portaudio"
    "math"
    "time"
    "fmt"
    "flag"
    "os"
    "bufio"
    "io"
    "strconv"
    "unicode"
    "unicode/utf8"
)

const sampleRate = 44100

const LETTER = 0
const DIGIT = 1
const OTHER = 2

/* quantum length in seconds, will get a click when this reaches 1, default to 1 year(!) */
const quantum float64 = 1.0 / (60 * 60 * 24 * 365)

var filename *string = flag.String("file", "", "Source file to read")  /* TODO: watch */
var semitone = math.Pow(2, 1.0/12)

var A = float64(440) / quantum
var sec = float64(1) / quantum

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

    s := newFloatbeat(in,sampleRate, 2)

    defer s.Close()
    chk(s.Start())
    time.Sleep(5 * time.Second)
    chk(s.Stop())
}

// ScanForthWords is a split function for a Scanner that 
// looks for a contiguous sequence of either
// unicode.IsLetter
// unicode.IsDigit or '.'
// unicode.IsSpace
// other
// and either ignores it, in the case of IsSpace,
// or returns it.
// It will never return an empty string.
func ScanForthWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
    // Skip leading spaces.
    start := 0
    var r rune

    for width := 0; start < len(data); start += width {
        r, width = utf8.DecodeRune(data[start:])
        if !unicode.IsSpace(r) {
            break
        }
    }

    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }

    var seq int
    switch {
        case unicode.IsLetter(r):
            seq = LETTER
        case unicode.IsDigit(r), r == '.':
            seq = DIGIT
        default:
            seq = OTHER
    }

    // Scan until rune not matching current set.
    for width, i := 0, start; i < len(data); i += width {
        r, width = utf8.DecodeRune(data[i:])
        if (unicode.IsLetter(r) != (seq == LETTER)) ||
           ((unicode.IsDigit(r) || r == '.') != (seq == DIGIT)) ||
           (unicode.IsSpace(r)) {
                return i, data[start:i], nil
        }
    }
    // If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
    if atEOF && len(data) > start {
        return len(data), data[start:], nil
    }
    // Request more data.
    return 0, nil, nil
}

type Floatbeat struct {
    *portaudio.Stream
    step, phase, clip float64
    code []string
}

func newFloatbeat(in io.Reader, sampleRate float64, clip float64) *Floatbeat {

    // Set phase = 0.5 s to avoid zeros everywhere during dummy run
    s := &Floatbeat{nil, quantum / sampleRate, 0.1/sec, clip, nil}
    fmt.Println("step is",s.step)

    s.code = s.read_code(in)

    fmt.Println("Dummy run produced",s.run(true))
    s.phase = 1.1/sec
    fmt.Println("Dummy run produced",s.run(true))
    s.phase = 0

    var err error
    s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
    chk(err)
    return s
}

func (f Floatbeat) read_code(in io.Reader) []string {
    scanner := bufio.NewScanner(in)

    scanner.Split(ScanForthWords)

    var result []string

    for scanner.Scan() {
        word := scanner.Text()
        result = append(result, word)
    }
    fmt.Println("Code:",result)
    chk(scanner.Err())

    return result
}

func (f Floatbeat) sin(freq float64) float64 {
    return math.Sin(2 * math.Pi * f.phase * freq)
}

func (f Floatbeat) saw(freq float64) float64 {
    return math.Mod(f.phase * freq, 2) - 1
}

func (f Floatbeat) sq(freq float64) float64 {
    return math.Copysign(1, math.Mod(f.phase * freq, 2) - 1)
}

func (f *Floatbeat) run(debug bool) []float64 {
    var stack []float64
    var pop float64

    for _, w := range f.code {
        l := len(stack)-1
        switch w {
            case "t":
                stack = append(stack, f.phase * sec)

            /* Forth words */

            case "frac":
                stack[l] = math.Mod(stack[l],1)

            case "*":
                pop, stack = stack[l], stack[:l]
                stack[l-1] *= pop
            case "/":
                pop, stack = stack[l], stack[:l]
                stack[l-1] /= pop
            case "+":
                pop, stack = stack[l], stack[:l]
                stack[l-1] += pop
            case "-":
                pop, stack = stack[l], stack[:l]
                stack[l-1] -= pop
            case "max":
                pop, stack = stack[l], stack[:l]
                stack[l-1] = math.Max(pop,stack[l-1])
            case "min":
                pop, stack = stack[l], stack[:l]
                stack[l-1] = math.Min(pop,stack[l-1])

            /* musical words */

            case "A":
                stack = append(stack, A)
            case "#":
                stack[l] *= semitone
            case "b":
                stack[l] /= semitone
            case "'":
                stack[l] *= 2
            case ",":
                stack[l] /= 2
            case "sin":
                stack[l] = f.sin(stack[l])
            case "saw":
                stack[l] = f.saw(stack[l])
            case "sq":
                stack[l] = f.sq(stack[l])
            default:
                num, err := strconv.ParseFloat(w, 64); chk(err)
                stack = append(stack, num)
        }
        if debug == true {
            fmt.Println(w,"::",stack)
        }
    }
    return stack
}

func (f *Floatbeat) processAudio(out [][]float32) {
    for i := range out[0] {

        stack := f.run(false)

        /* Add up whatever is left on the stack */
        r := float64(0)
        for _, s := range stack {
            r += s
        }
        /* Tweak the scale if it's clipping */
        if r > f.clip {
            f.clip = r
            fmt.Println("clip", f.clip)
        }
        out[0][i] = float32(r / f.clip)
        _, f.phase = math.Modf(f.phase + f.step)
    }
}

func chk(err error) {
    if err != nil {
        panic(err)
    }
}

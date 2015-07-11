floatbeat
=========

A programmatic sound/music generator inspired by bytebeat.
Also an excuse to do something in Go.

Uses
* `code.google.com/p/portaudio-go/portaudio`
* `github.com/colourcountry/d4`
* `github.com/gorilla/websocket`

The generator takes as input a file which defines what it will initially play,
using a [Forth-like language](http://github.com/colourcountry/d4).

It also runs a web server which accepts replacement programs over a Web Socket.

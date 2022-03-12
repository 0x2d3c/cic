package main

import (
	"sync"
)

var (
	cfg        Config
	poolPacket = sync.Pool{New: func() interface{} { return &Packet{} }}
)

type Packet struct {
	Typ uint8
	Who string
	Msg string
}

type Config struct {
	Uri, Name, Addr string
}

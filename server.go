package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	ACK = iota + 1
	Join
	Notice
	Message
)

var (
	punishErr = "what you want, connection closed? OK"

	server     = StreamServer{kv: make(map[string]*Stream, 1024)}
	poolStream = sync.Pool{New: func() interface{} { return &Stream{} }}
)

type StreamServer struct {
	kv map[string]*Stream
	sync.RWMutex
}

func (ss *StreamServer) Join(s *Stream) {
	ss.RLock()
	ss.kv[s.id] = s
	ss.RUnlock()
}

func (ss *StreamServer) Abort(id string) {
	ss.RLock()
	delete(ss.kv, id)
	ss.RUnlock()
}

func (ss *StreamServer) Allot(m *Packet, id string) {
	for user, stream := range ss.kv {
		if user == id {
			continue
		}
		stream.conn.WriteJSON(m)
	}
}

type Stream struct {
	id   string
	wc   chan<- *Packet
	conn *websocket.Conn
}

func (stream *Stream) Run() {
	for {
		packet := PoolPacket.Get().(*Packet)
		if err := stream.conn.ReadJSON(&packet); err != nil {
			fmt.Println("read err", err.Error())
			break
		}

		switch packet.Typ {
		case Join:
			stream.id = packet.Who
			server.Join(stream)
		case Message:
			packet.Who = stream.id
			server.Allot(packet, stream.id)

			packet.Typ = ACK
			packet.Msg = ""

			stream.conn.WriteJSON(packet)
		default:
			packet.Typ = Notice
			packet.Msg = punishErr
			stream.conn.WriteJSON(packet)

			goto END
		}
	}

END:
	if stream.id != "" {
		server.Abort(stream.id)
	}
}

var upgrade = websocket.Upgrader{} // use default options
func Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	in := poolStream.Get().(*Stream)
	in.conn = c
	in.wc = make(chan<- *Packet)

	in.Run()
}

func ServerRun() {
	http.HandleFunc("/chat", Handler)
	if err := http.ListenAndServe(cfg.Addr, nil); err != nil {
		fmt.Println(err.Error())
	}
}

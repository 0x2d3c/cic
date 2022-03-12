package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

var cli StreamClient

type StreamClient struct {
	id   string
	conn *websocket.Conn
}

func (sc *StreamClient) In() {
	for {
		msg := PoolPacket.Get().(*Packet)
		if err := sc.conn.ReadJSON(&msg); err != nil {
			continue
		}
		switch msg.Typ {
		case ACK:
		case Message:
			fmt.Print("\n", msg.Who, " #", msg.Msg, "\n", cfg.Name, " >")
		default:
			break
		}
	}
}

func (sc *StreamClient) Out() {
	sc.conn.WriteJSON(&Packet{Typ: Join, Who: cfg.Name})
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print(cfg.Name, " >")
		txt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		msg := PoolPacket.Get().(*Packet)
		msg.Msg = txt
		msg.Typ = Message

		sc.conn.WriteJSON(msg)
	}
}

func ClientRun() {
	u := url.URL{Scheme: "ws", Host: cfg.Addr, Path: "/chat"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	cli.id = cfg.Name
	cli.conn = conn

	go cli.In()

	cli.Out()
}

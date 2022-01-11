package main

import (
	"log"
	"net"
	"syscall/js"

	"github.com/elewis787/hacket"
)

const (
	chat hacket.PacketType = iota
)

type hacketService struct {
	server hacket.PacketServer
	client hacket.PacketClient
	msg    chan string
}

func (hs *hacketService) sendMessage() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		addr, err := net.ResolveUDPAddr("udp", args[0].String())
		if err != nil {
			return err.Error()
		}
		hs.client.WriteTo([]byte(args[1].String()), addr)
		return "sent"
	})
}

func (hs *hacketService) getMessage() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(<-hs.msg)
	})
}

func newHacketService(this js.Value, args []js.Value) interface{} {
	options := []hacket.Options{}

	s, c, err := hacket.New("udp", "127.0.0.1:1234", options...)
	if err != nil {
		return err.Error()
	}

	hs := &hacketService{
		server: s,
		client: c,
		msg:    make(chan string),
	}

	mux := hacket.NewPacketMux()
	mux.PacketHandlerFunc(chat, func(packet hacket.Packet, pw hacket.PacketWriter) {
		if len(packet.Msg()) > 0 {
			hs.msg <- string(packet.Msg())
		}
	})

	go func() {
		if err := hs.server.Serve(mux); err != nil {
			log.Println(err)
		}
	}()

	return js.ValueOf(map[string]interface{}{
		"SendMessage": hs.sendMessage(),
		"GetMessage":  hs.getMessage(),
	})
}

func main() {
	c := make(chan struct{})
	js.Global().Set("NewHacketService", js.FuncOf(newHacketService))
	<-c
}

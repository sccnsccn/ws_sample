// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const fileConfig = "config.json"

func rcvMsg(connection *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		if mt == websocket.CloseMessage {
			return
		}
	}
}

// go run client.go config_reader.go
func main() {
	log.SetFlags(0)

	config := ReadConfig(fileConfig)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: config.GetHost(), Path: config.Path}
	log.Printf("connecting to %s", u.String())

	var header = http.Header{
		config.CustomHeader: {config.Token},
	}

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = config.GetSecondsTimeout()
	connection, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err.Error())
	}
	log.Println("connect to success")

	defer connection.Close()

	done := make(chan struct{})

	go rcvMsg(connection, done)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

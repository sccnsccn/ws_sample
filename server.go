// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	fileTokens     = "valid_tokens.json"
	fileValidLines = "valid_tokens.json"
	customHeader   = "X-Authorization"
	host           = "localhost"
	port           = 8080
	path           = "/connect"
)

type ServerWS struct {
	clients         map[*websocket.Conn]bool
	token_validator TokenValidator
	msg_generator   MsgGenerator
	upgrader        websocket.Upgrader
}

func newServerWS() *ServerWS {

	token_validator := TokenValidator{Reader: StringJsonReader{}}
	msg_generator := MsgGenerator{Reader: StringJsonReader{}}

	token_validator.LoadTokens(fileTokens)
	msg_generator.LoadValidLine(fileValidLines)

	server := ServerWS{
		clients:         make(map[*websocket.Conn]bool),
		token_validator: token_validator,
		msg_generator:   msg_generator,
	}

	server.upgrader = websocket.Upgrader{
		CheckOrigin: server.CheckOrigin,
	}

	return &server
}

func (s *ServerWS) CheckOrigin(r *http.Request) bool {
	var token = r.Header.Get(customHeader)
	log.Print("token: ", token)

	if !s.token_validator.CheckToken(token) {
		log.Print("invalid token ", token)
		return false
	}
	return true
}

func (s *ServerWS) Connect(w http.ResponseWriter, r *http.Request) {
	connection, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	s.clients[connection] = true

	go s.RcvMsg(connection)
}

func (s *ServerWS) RcvMsg(connection *websocket.Conn) {
	defer connection.Close()
	defer delete(s.clients, connection)
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %d %s", mt, message)
		if mt == websocket.CloseMessage {
			return
		}
	}
}

func (s *ServerWS) SendMsg() {
	for {
		time.Sleep(time.Second)
		for conn := range s.clients {
			msg := s.msg_generator.GenerateMsg()
			text := strings.Join(msg[:], "\r\n")
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}
}

// go run token_validator.go string_reader.go server.go msg_generator.go
func main() {
	log.SetFlags(0)

	serverWS := newServerWS()

	go serverWS.SendMsg()

	http.HandleFunc(path, serverWS.Connect)
	log.Println("server start")
	log.Fatal(http.ListenAndServe(host+":"+strconv.Itoa(port), nil))

}

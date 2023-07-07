package main

import (
	"math/rand"
)

type MsgGenerator struct {
	Reader StringJsonReader
}

func (mg *MsgGenerator) LoadValidLine(fileNmae string) {
	mg.Reader.LoadData(fileNmae)
}

func (mg *MsgGenerator) GenerateMsg() []string {

	len := len(mg.Reader.Data)
	count := rand.Intn(4) + 1

	var msg []string

	for i := 1; i <= count; i++ {
		rnd := rand.Intn(len)
		value := mg.Reader.Data[rnd]
		msg = append(msg, value)
	}

	return msg
}

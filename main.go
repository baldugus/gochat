package main

import (
	"github.com/baldugus/gochat/internal/broker"
)

func main() {
	client := broker.NewClient()
	defer client.Close()

	StartUi(client)
}

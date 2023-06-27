package main

import (
	"flag"
	"fmt"
	"websockets/internal/chat"
)

func main() {
	flag.Parse()
	err := chat.StartServer()
	if err != nil {
		fmt.Printf("An error occurred during execution: %v", err)
	}
}

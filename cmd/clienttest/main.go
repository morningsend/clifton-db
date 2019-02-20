package main

import (
	"github.com/zl14917/MastersProject/kvclient"
	"log"
)

const PORT = 10030

func main() {
	client := kvclient.NewClient()
	log.Println("connecting to server on ")
	err := client.Connect("localhost", PORT)

	if err != nil {
		log.Fatalln("failed to connect to server:", err)
	}

	log.Println("PUT value: <hello:world>")
	_, err = client.Put("hello", []byte("world"))

	if err != nil {
		log.Println("failed to put", err)
	}
}

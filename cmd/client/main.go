package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marten-seemann/webtransport-go"
)

func main() {

	var d webtransport.Dialer
	rsp, conn, err := d.Dial(context.Background(), "https://localhost:4433/counter", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("resp status: %d\n", rsp.StatusCode)
	fmt.Printf("session opened to %s", conn.RemoteAddr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(nil)
	}
	n, err := stream.Write([]byte("HELLO WORLD"))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sent %d bytes\n", n)
	stream.Close()
	conn.Close()
}

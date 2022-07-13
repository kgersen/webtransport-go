package main

import (
	"context"
	"fmt"
	"io"
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
	fmt.Printf("session opened to %s\n", conn.RemoteAddr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		log.Fatal(err)
	}
	n, err := stream.Write([]byte("HELLO WORLD\n"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sent %d bytes\n", n)
	// fmt.Printf("sleeping 3 seconds\n")
	// time.Sleep(time.Second * 5)
	stream.Close()
	// stream is closed but we can read...
	reply, err := io.ReadAll(stream)
	fmt.Printf("got reply:%s\n", reply)
	conn.Close()
}

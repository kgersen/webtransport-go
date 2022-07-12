package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

func main() {
	s := webtransport.Server{
		H3: http3.Server{Addr: ":4433"},
	}

	// Create a new HTTP endpoint /webtransport.
	http.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}
		// Handle the connection. Here goes the application logic.
		fmt.Printf("got something from %s, waiting for a stream\n", conn.RemoteAddr())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		stream, err := conn.AcceptStream(ctx)
		fmt.Printf("got a stream\n")

		if errors.Is(err, &webtransport.ConnectionError{Message: "EOF"}) {
			fmt.Println("got EOF")
		}
		if err != nil && !errors.Is(err, io.EOF) {
			s := fmt.Sprintf("error on AcceptStream: %+v", err)
			http.Error(w, s, http.StatusInternalServerError)
			fmt.Println(s)
			return
		}

		// dump stream to stdout
		fmt.Println("copy started")
		n, err := io.Copy(os.Stdout, stream)
		fmt.Println("copy ended")
		if err != nil {
			s := fmt.Sprintf("error on Copy: %v", err)
			http.Error(w, s, http.StatusInternalServerError)
			fmt.Println(s)
			return
		}
		fmt.Printf("%d bytes received\n", n)
		stream.Write([]byte("SERVER IS OK"))
		stream.Close()
		fmt.Printf("replied and closed\n")
	})

	err := s.ListenAndServeTLS("localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

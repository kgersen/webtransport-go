package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
		stream, err := conn.AcceptStream(r.Context())
		fmt.Printf("got a stream\n")
		if err != nil {
			s := fmt.Sprintf("error on AcceptStream: %v", err)
			http.Error(w, s, http.StatusInternalServerError)
			fmt.Println(s)
		}
		// dump stream to stdout
		n, err := io.Copy(os.Stdout, stream)
		if err != nil {
			s := fmt.Sprintf("error on Copy: %v", err)
			http.Error(w, s, http.StatusInternalServerError)
			fmt.Println(s)
		}
		fmt.Printf("%d bytes received\n", n)

	})

	err := s.ListenAndServeTLS("localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

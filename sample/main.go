package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mpao/ssc"
)

func main() {
	t, err := ssc.NewTrust("localhost.crt", "localhost.key", "localCA.crt", tls.RequireAndVerifyClientCert)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello server"))
	})
	go func() {
		log.Println("server started at port 8080")
		if e := t.StartServer(mux, 8080); e != nil {
			log.Fatal(e)
		}
	}()
	time.Sleep(1 * time.Second)
	client := t.Client(10 * time.Second)
	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		log.Println(err)
	}
	body := resp.Body
	b, _ := io.ReadAll(body)
	resp.Body.Close()
	log.Println("***** Trusted client:")
	log.Printf("> %s\n", string(b))
	log.Println("***** Bad client error:")
	_, err = http.Get("https://localhost:8080")
	if err != nil {
		log.Println(err)
	}
	log.Println("Server closed, BYE :)")
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lucas-clemente/quic-go/http3"
	log "github.com/sirupsen/logrus"
)

func main() {

	address := flag.String("addr", "localhost", "Specify server address")
	port := flag.Int("port", 4242, "Specify server port")
	cert := flag.String("cert", "../keys/cert.pem", "Specify Certificate path")
	key := flag.String("key", "../keys/priv.key", "Specify Servers Private key")
	flag.Parse()

	sig := make(chan os.Signal, 1)

	addr := fmt.Sprintf("%s:%d", *address, *port)
	http.HandleFunc("/echo", echoHandler)
	log.Printf("Server is running...\n")
	if err := http3.ListenAndServeQUIC(addr, *cert, *key, nil); err != nil {
		log.Errorf("Got an error while listening: %v\n", err)
	}

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	close(sig)

}

func echoHandler(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Unable to read request body: %v\n", err)
	}
	log.Printf("Received: %s\n", body)
	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Errorf("Unable to close request body: %v\n", err)
		}
	}()
	resp := fmt.Sprintf("Server Echo's: %s\n", body)
	log.Printf("Sending: %s\n", resp)
	if _, err := res.Write([]byte(resp)); err != nil {
		log.Errorf("Unable to write the response: %v\n", err)
	}
}

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
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
	addr := flag.String("addr", "localhost", "Specify server address")
	port := flag.Int("port", 4242, "Specify server port")
	keyPath := flag.String("keys", "../keys/ca.pem", "Specify Certificate Authority key path")
	msg := flag.String("msg", "hello http/3", "Specify message")
	flag.Parse()

	sig := make(chan os.Signal, 1)

	caCertRaw, err := ioutil.ReadFile(*keyPath)
	if err != nil {
		panic(err)
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	if ok := pool.AppendCertsFromPEM(caCertRaw); !ok {
		panic("Could not add root ceritificate to pool.")
	}

	rt := &http3.RoundTripper{TLSClientConfig: &tls.Config{
		RootCAs: pool,
	},
	}
	defer func() {
		if err := rt.Close(); err != nil {
			log.Error("Unable to close round-tripper")
		}
	}()

	client := http.Client{
		Transport: rt,
	}

	url := fmt.Sprintf("https://%s:%d/echo", *addr, *port)

	log.Printf("Sending following message: %s\n", *msg)
	req, err := http.NewRequest(http.MethodGet, url, ioutil.NopCloser(bytes.NewReader([]byte(*msg))))
	if err != nil {
		log.Errorf("Unable to create HTTP request\n")
	}
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Got error while sending GET request: %v\n", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error("Unable to close response body")
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Unable to read response body: %v\n", err)
	}
	log.Printf("Received message: %s\n", body)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	close(sig)
}
